// API client for the Location Demo backend.
// All backend communication is centralized here.

const IS_SERVER = typeof window === 'undefined';
const API_BASE = IS_SERVER
  ? (process.env.INTERNAL_API_URL || "http://backend:8088/api/v1")
  : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8088/api/v1");


export interface SearchResult {
  id: number;
  name: string;
  type: string;
  country?: string;
}

export interface LocationDetail {
  id: number;
  name: string;
  type: string;
  lat: number;
  lng: number;
  is_verified?: boolean;
  slug?: string;
  parent?: {
    id: number;
    name: string;
    type: string;
  };
  stats?: {
    total_posts: number;
    total_photos: number;
    total_videos: number;
    trending_score: number;
  };
  translations?: {
    location_id: number;
    lang_code: string;
    name: string;
  }[];
}

export interface TrendingLocation {
  location_id: number;
  name: string;
  type: string;
  score: number;
  date: string;
}

export interface Post {
  id: number;
  user_id: number;
  content: string;
  media_type: string;
  location_id: number;
  location_name: string;
  location_type: string;
  created_at: string;
}

export interface ApiResponse<T> {
  data: T;
}

/**
 * Search locations by query string.
 * Uses the Waterfall Strategy on the backend: alias → translation → external API.
 */
export async function searchLocations(
  query: string,
  lang: string = "en"
): Promise<SearchResult[]> {
  const params = new URLSearchParams({ q: query, lang });
  const res = await fetch(`${API_BASE}/locations/search?${params}`);

  if (!res.ok) {
    throw new Error(`Search failed: ${res.statusText}`);
  }

  const json: ApiResponse<SearchResult[]> = await res.json();
  return json.data || [];
}

/**
 * Get a single location's full detail by ID.
 */
export async function getLocation(
  id: number,
  lang: string = "en"
): Promise<LocationDetail> {
  const params = new URLSearchParams({ lang });
  const res = await fetch(`${API_BASE}/locations/${id}?${params}`);

  if (!res.ok) {
    throw new Error(`Location not found: ${res.statusText}`);
  }

  const json: ApiResponse<LocationDetail> = await res.json();
  return json.data;
}

/**
 * Get child locations (e.g., cities in a country).
 */
export async function getChildren(
  parentId: number,
  lang: string = "en"
): Promise<SearchResult[]> {
  const params = new URLSearchParams({ lang });
  const res = await fetch(
    `${API_BASE}/locations/${parentId}/children?${params}`
  );

  if (!res.ok) {
    throw new Error(`Children fetch failed: ${res.statusText}`);
  }

  const json: ApiResponse<SearchResult[]> = await res.json();
  return json.data || [];
}

/**
 * Get trending locations
 */
export async function getTrending(lang: string = "en", limit: number = 5): Promise<TrendingLocation[]> {
  const params = new URLSearchParams({ lang, limit: limit.toString() });
  const res = await fetch(`${API_BASE}/locations/trending?${params}`);
  if (!res.ok) throw new Error("Failed to fetch trending locations");
  const json: ApiResponse<TrendingLocation[]> = await res.json();
  return json.data || [];
}

/**
 * Get posts for a specific location
 */
export async function getPostsByLocation(locationId: number, lang: string = "en"): Promise<Post[]> {
  const params = new URLSearchParams({ lang });
  const res = await fetch(`${API_BASE}/locations/${locationId}/posts?${params}`);
  if (!res.ok) throw new Error("Failed to fetch posts");
  const json: ApiResponse<Post[]> = await res.json();
  return json.data || [];
}

/**
 * Create a new post
 */
export async function createPost(content: string, mediaType: string, locationId: number): Promise<Post> {
  const res = await fetch(`${API_BASE}/posts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      content,
      media_type: mediaType,
      location_id: locationId
    })
  });
  
  if (!res.ok) throw new Error("Failed to create post");
  const json: ApiResponse<Post> = await res.json();
  return json.data;
}
