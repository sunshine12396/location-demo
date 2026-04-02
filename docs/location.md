## Location System Design Document

This document outlines the design and architecture for a robust location system.

### 1. System Objectives

The Location system supports the following key functionalities:

* **Assign Locations:** Attach specific locations to user posts.

* **Efficient Display:** Rapidly display location names without repeated external API calls.

* **Trending Locations:** Identify and rank popular or trending locations.

* **Post Filtering:** Enable users to filter posts by city, country, or landmark.

* **Multi-language Support:** Provide location names in multiple languages.

* **Nearby Search:** Facilitate discovery of nearby locations and posts (leveraging PostGIS if implemented).

* **Scalability:** Ensure the system scales effectively with large datasets.

### 2. Database Architecture

#### Core Tables Overview

The system design revolves around a set of interconnected tables:

* `locations`: Stores primary location data, excluding localized names. This includes external IDs, types, geographic coordinates, hierarchical relationships, SEO slugs, and verification status.

* `location_translations`: Manages multi-language support for location names, linking to `locations` via `location_id` and `lang_code`.

* `location_alias`: Facilitates flexible search by storing alternative names, spellings, and abbreviations for locations.

* `posts`: Contains user-generated content, linked to specific locations.

* `location_stats`: Optimizes performance by pre-aggregating statistics (e.g., total posts, photos, videos) for locations, avoiding costly `COUNT(*)` operations on large datasets.

* `trending_locations`: Stores calculated trending scores for locations on a daily basis.

#### Table Schemas

**2.1 Table:&#x20;**`locations`\
Stores core location information.

```sql
CREATE TABLE locations (
    id BIGINT PRIMARY KEY,
    external_id VARCHAR UNIQUE,      -- Google Place ID / Mapbox ID
    type VARCHAR,                    -- country / city / district / landmark / venue / address
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    parent_id BIGINT,                -- Hierarchy (city -> country, landmark -> city)
    path VARCHAR,                    -- Hierarchy path: 1/2/3/4
    slug VARCHAR UNIQUE,             -- Used for SEO-friendly URLs
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**2.2 Table:&#x20;**`location_translations`\
Supports multi-language location names.

```sql
CREATE TABLE location_translations (
    id BIGINT PRIMARY KEY,
    location_id BIGINT,
    lang_code VARCHAR(5),           -- vi / en / ja / ko / fr
    name VARCHAR,
    UNIQUE(location_id, lang_code)
);
```

**2.3 Table:&#x20;**`location_alias`\
Enhances search flexibility.

```sql
CREATE TABLE location_alias (
    id BIGINT PRIMARY KEY,
    location_id BIGINT,
    alias VARCHAR
);
```

**2.4 Table:&#x20;**`posts`\
Stores user posts with associated location.

```sql
CREATE TABLE posts (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,
    content TEXT,
    location_id BIGINT,
    created_at TIMESTAMP
);
```

**2.5 Table:&#x20;**`location_stats`\
Optimizes performance for location-based statistics.

```sql
CREATE TABLE location_stats (
    location_id BIGINT PRIMARY KEY,
    total_posts BIGINT DEFAULT 0,
    total_photos BIGINT DEFAULT 0,
    total_videos BIGINT DEFAULT 0,
    last_post_at TIMESTAMP,
    trending_score DOUBLE PRECISION
);
```

**2.6 Table:&#x20;**`trending_locations`\
Records daily trending scores for locations.

```sql
CREATE TABLE trending_locations (
    location_id BIGINT,
    score DOUBLE PRECISION,
    date DATE
);
```

### 3. Location Hierarchy Management

The system utilizes `parent_id` and `path` columns within the `locations` table to manage hierarchical relationships effectively. This allows for clear representation of geographical and administrative divisions.

**Example Hierarchy:**

```javascript
Vietnam
└── Ho Chi Minh City
    └── District 1
        └── Bitexco Tower
```

**Example Data Representation:**

| id | name             | parent_id | path    |
| -- | ---------------- | --------- | ------- |
| 1  | Vietnam          | NULL      | 1       |
| 2  | Ho Chi Minh City | 1         | 1/2     |
| 3  | District 1       | 2         | 1/2/3   |
| 4  | Bitexco Tower    | 3         | 1/2/3/4 |

**Query Example: Retrieve all posts in Ho Chi Minh City**

```sql
SELECT p.*
FROM posts p
JOIN locations l ON l.id = p.location_id
WHERE l.path LIKE '1/2/%';
```

### 4. Multi-language Support

Location names are stored in the `location_translations` table, enabling display in various languages based on user preferences.

**Example Data:**

| location_id | lang | name           |
| ----------- | ---- | -------------- |
| 10          | vi   | Hồ Hoàn Kiếm   |
| 10          | en   | Hoan Kiem Lake |
| 10          | ja   | ホアンキエム湖        |

**Query Example: Retrieve posts with location names in English**

```sql
SELECT p.*, lt.name
FROM posts p
JOIN location_translations lt
ON lt.location_id = p.location_id
AND lt.lang_code = 'en';
```

### 5. Core Operational Flow

**Step 1: User Searches for a Location**\
The application queries external APIs (e.g., Google Places API, Mapbox API) to find potential locations.

**Step 2: User Selects a Location**\
The backend processes the user's selection:

1. Checks if the `external_id` already exists in the `locations` table.

2. If the location is new, it inserts the data into the `locations` table.

3. Inserts the location's translation(s) based on the user's language into `location_translations`.

4. Inserts any necessary parent hierarchy data (e.g., city, country) into `locations`.

**Step 3: User Creates a Post**\
The system records the post with its associated `location_id`.

```sql
INSERT INTO posts (user_id, content, location_id) VALUES (?, ?, ?);
```

**Step 4: Update Location Statistics**\
The system increments relevant statistics in the `location_stats` table.

```sql
UPDATE location_stats SET total_posts = total_posts + 1 WHERE location_id = ?;
```

### 6. Trending Location Logic

Relying solely on `total_posts` is insufficient for accurately identifying trending locations. A more comprehensive trending score is proposed:

**Proposed Trending Score Formula:**

```javascript
score = (post_count * 1.0) +
        (photo_count * 1.5) +
        (video_count * 2.0) +
        (comment_count * 0.3) +
        (share_count * 2.5) +
        recent_post_weight
```

This formula assigns different weights to various engagement metrics and incorporates a factor for recent activity.

### 7. Nearby Search (Optional with PostGIS)

PostGIS is highly recommended for systems requiring advanced geospatial queries, such as "Nearby posts," "Nearby cafes," or "Explore location map" features.

**Suggested Column:**

```sql
location GEOGRAPHY(POINT)
```

**Query Example: Find locations within a 2km radius**

```sql
SELECT *
FROM locations
WHERE ST_DWithin(
    location,
    ST_Point(:lng, :lat),
    2000
);
```

### 8. Key Enhancements

* **Standardized Location Types:** Implement a comprehensive set of location types (e.g., country, state/region, city, district, ward, street, landmark, venue, airport, building) for better categorization and filtering.

* **Verified Locations:** Utilize an `is_verified` boolean flag to prevent user-generated fake locations and ensure data integrity.

* **SEO-Friendly Slugs:** Generate unique, human-readable slugs for locations (e.g., `/location/ho-chi-minh-city`, `/location/hoan-kiem-lake`) to improve SEO and URL structure.

* **Location Aliases for Enhanced Search:** Leverage the `location_alias` table to support flexible search queries, including common misspellings, abbreviations, and alternative names (e.g., "sai gon," "hcm," "tp hcm").

### 9. Final Architecture Summary

The core architecture for the Location System comprises the following key tables:

* `locations`

* `location_translations`

* `location_alias`

* `location_stats`

* `trending_locations`

* `posts`

### 10. Current Implementation Status

The Location Demo System has been fully implemented out of the Phase 1 and Phase 2 specifications. 

**What is completed:**
1. **Core Database Schema**: All tables (`locations`, `location_translations`, `location_alias`, `posts`, `location_stats`, `trending_locations`) are fully implemented and seeded via Docker Postgres init scripts.
2. **Waterfall Search**: The backend uses local Aliases, local Translations, and falls back to **OpenStreetMap (OSM) Nominatim** for dynamic external discovery with intelligent deduplication and caching.
3. **Hierarchical Graph**: Navigation seamlessly uses `parent_id` arrays to render breadcrumbs, while Post fetching leverages the string-based `path` column (e.g. `1/2/%`) to retrieve community posts across an entire country or city in a single query.
4. **Interactive UI**: A Next.js frontend is built utilizing Tailwind CSS and glassmorphism. It provides multi-language support (English, Vietnamese, Japanese), verified location badges, dynamic SEO slug display, and live post creation right in the browser.

**What is remaining:**
- PostGIS Nearby Search: The `GEOGRAPHY(POINT)` data model is proposed but currently excluded to maintain a simple standard PostgreSQL dependency.

### Appendix: Running the Demo Local Environment

The entire stack is containerized for instant local development.
Ensure Docker is installed and running on your system, then simply execute:
```bash
make up
```

This will automatically spin up:
- **db**: PostgreSQL instance mapped to port `5433`. It auto-runs the phase 1 and phase 2 schema and seed SQL upon first startup.
- **backend**: The Golang API server running on port `8088`.
- **frontend**: The Next.js application running on port `3001`.

*Note: Use `make clean` to completely destroy the containers and volumes if you need to perform a fresh database seed.*
