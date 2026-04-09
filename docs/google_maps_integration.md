# Google Maps Integration Guide

This application integrates with the **Google Maps Places API (New)** to provide robust location discovery, normalization, and hierarchical translation.

## 1. Setup and Authentication

1. **Create a Project:** In the [Google Cloud Console](https://console.cloud.google.com/), create a new project.
2. **Enable Billing:** A linked billing account is required, even for the free tier.
3. **Enable APIs:** Go to **APIs & Services > Library** and enable the **Places API (New)**.
4. **Generate API Key:** Under **Credentials**, create an API Key. **Crucial:** Restrict the key to only the Places API and limit it to your application's IPs/domains to prevent abuse.

## 2. Location Type Mapping

Google returns an array of `types` for a location. We map these to our internal, simplified categorization system. Always evaluate types from **most specific to least specific** (e.g., check for a business venue before falling back to a general address).

| Internal Type | Primary Google Types | Notes |
| :--- | :--- | :--- |
| **country** | `country` | Top-level node (e.g., Việt Nam). |
| **city** | `administrative_area_level_1` | Tỉnh/Thành phố trực thuộc TW (e.g., TP. Hồ Chí Minh, Bình Dương). |
| **district** | `administrative_area_level_2`, `locality` | Quận/Huyện hoặc Thành phố thuộc tỉnh (e.g., Quận Bình Thạnh, TP. Thủ Đức). *Note: Google uses locality for provincial cities/towns.* |
| **ward** | `sublocality_level_1`, `sublocality` | Phường/Xã (e.g., Phường 22). |
| **landmark** | `tourist_attraction`, `point_of_interest`, `natural_feature`, `park`, `museum`, `stadium` | Specific famous places or natural features. |
| **venue** <br/>*(business/events)* | `restaurant`, `cafe`, `hotel`, `store`, `gym`, `establishment` | Specific businesses. `establishment` is the generic fallback. |
| **address** | `street_address`, `route`, `premise` | Specific street or building addresses. |

> **Note on "political":** Almost all geographies (Country, City, District) will include the type `political`. Ignore it, as it is too broad.

### Important Types (Hierarchy Filter)

When processing `addressComponents`, only these Google types are considered for building the parent hierarchy. All other types (e.g., `postal_code`, `route`) are ignored.

```go
importantTypes := map[string]bool{
    "country":                     true,
    "administrative_area_level_1": true, 
    "administrative_area_level_2": true, 
    "locality":                    true, 
    "sublocality":                 true,
}
```

## 3. The Data Mapping Strategy

When converting Google's JSON response to our SQL schema:

| Database Column | Google Maps Field | Notes |
| :--- | :--- | :--- |
| **`locations.external_id`** | `id` (Place ID) | Unique Google identifier. Resolved via `searchText` for parent components. |
| **`locations.type`** | `types` | Mapped using the table in Section 2. |
| **`locations.lat` / `lng`** | `location.latitude` / `.longitude` | Geographic coordinates. |
| **`locations.parent_id`** | Computed | FK to the parent location in the hierarchy. |
| **`locations.path`** | Computed | Materialized path string (e.g., `/1/4/22/`). |
| **`locations.provider`** | Hardcoded | Always `'google'`. |
| **`location_translations`** | `displayName.text` | Stored per requested `lang_code`. |
| **`location_alias`** | `formattedAddress` | Stored to allow searching via full physical address. |

## 4. End-to-End Implementation Flow

The integration involves three major steps across the frontend and backend.

### Step 1: Client-Side Search (Autocomplete Proxy)
The frontend calls the internal backend endpoint (`/api/v1/locations/autocomplete`) which proxies to Google. This keeps our API key secure and allows us to pre-filter or cache suggestions if needed.

**Backend Implementation:**
```go
reqURL := "https://places.googleapis.com/v1/places:autocomplete"
payload := map[string]interface{}{
    "input": query,
    "languageCode": lang,
}
// Returns place_id (external_id) for Step 2
```

### Step 2: Post Creation & JIT Hydration
When a user selects a location and clicks "Create Post", the frontend sends the `external_id` (Google Place ID) to our backend.
The backend's `CreatePost` service then triggers the **Just-in-Time Hydration** flow:
1. **Check Existence:** Does the `external_id` exist in our `locations` table?
2. **Details Fetch:** If missing, the backend calls the Google Place Details API (REST) using a field mask.
3. **Hierarchy Sync:** The `IngestGooglePlace` flow (Step 3) is triggered to build the full parent chain.
4. **Instant Save:** The new location is saved to our DB, and the post is immediately linked to it.

### Step 3: Hierarchy Generation & Saving (`IngestGooglePlace`)

This is the core hydration logic. It receives the full Google Place Details response and ensures a complete parent hierarchy exists in the local DB before saving the target location.

#### 3a. Helper: `resolveNameToID` (Name → Google Place ID)

When a parent component (e.g., "Quận Bình Thạnh") is missing from the DB, we need its Google Place ID. This helper resolves it via the `searchText` API:

```go
func resolveNameToID(ctx context.Context, name string, includedType string, lat, lng float64) (string, error) {
    // 1. Check cache first (avoid redundant API calls)
    if id, found := cache.Get(name); found {
        return id.(string), nil
    }

    // 2. Call Google Places searchText with strict filtering
    payload := map[string]interface{}{
        "textQuery":    name,
        "languageCode": "vi",           // Vietnamese for admin accuracy
        "includedType": includedType,   // e.g., "administrative_area_level_1"
        "locationBias": map[string]interface{}{
            "circle": map[string]interface{}{
                "center": map[string]interface{}{
                    "latitude": lat, "longitude": lng,
                },
                "radius": 5000.0, // 5km radius for disambiguation
            },
        },
    }
    // POST https://places.googleapis.com/v1/places:searchText
    // X-Goog-FieldMask: places.id
    // Return: result.places[0].id
}
```

**Key parameters:**
- **`languageCode: "vi"`**: Ensures Vietnamese administrative names are matched correctly.
- **`includedType`**: Constrains results to the exact administrative layer (prevents "Hồ Chí Minh" from resolving to a restaurant instead of the city).
- **`locationBias`**: Uses coordinates from the original place to disambiguate common names (e.g., "Quận 1" in HCMC vs. "Quận 1" in Hà Nội).

#### 3b. The Ingestion Loop

```
 Google addressComponents (specific → general):
   [Landmark 81, Phường 22, Q. Bình Thạnh, TP.HCM, Việt Nam]
                          ↓ reverse
 Processing order (general → specific):
   [Việt Nam, TP.HCM, Q. Bình Thạnh, Phường 22]
                          ↓ loop
 DB result:
   Việt Nam (id=1) → TP.HCM (id=4, parent=1) → Q. Bình Thạnh (id=22, parent=4) → ...
```

**Steps:**

1. **Reverse** the `addressComponents` array so we process from root (Country) to leaf.
2. **Filter** each component — only process those whose `types[0]` is in the `importantTypes` map.
3. **For each important component:**
   - **Search Local DB** by `name`, `type`, and `parent_id` (scoped by `lastParentID`).
   - **If not found (JIT Resolution):**
     - Call `resolveNameToID(name, includedType, lat, lng)` to get the Google Place ID.
     - `INSERT` the new location with `external_id`, `type`, `parent_id = lastParentID`, and computed `path`.
   - **Update loop state:** Set `lastParentID = loc.ID` and extend `currentPath`.
4. **Finally**, upsert the target location itself (e.g., Landmark 81) with `parent_id = lastParentID` and the fully built `path`.

#### 3c. Database Transactions

All inserts within the ingestion happen inside a single transaction:
   - `INSERT` into `locations` (core data + computed path).
   - `INSERT` into `location_translations` (using `displayName.text`).
   - `INSERT` into `location_alias` (using `formattedAddress`).

## 5. Lifecycle & Data Refresh (Auto-Sync)

To handle "outdated info" while staying compliant with Google's caching policies:
*   **`updated_at` Tracking:** Every location in our DB has an `updated_at` column.
*   **Threshold-based Sync:** In the `GetByID` (view) and `CreatePost` (hydration) flows, the backend checks if `time.Since(updated_at) > LOCATION_SYNC_DAYS` (default: 30 days).
*   **Background Refresh:** If a location is outdated, a background goroutine re-fetches the Google Details API and overwrites the local record (coordinates, address, and primary translation) to keep it fresh.
*   **Permanent ID:** We store the `external_id` (Place ID) permanently as the stable anchor for these refreshes.