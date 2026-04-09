-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS ltree;
-- =========================
-- 1. LOCATIONS
-- =========================
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,

    -- External reference (e.g. Google Place ID)
    external_id VARCHAR(255) UNIQUE,

    -- Internal classification
    type VARCHAR(50) NOT NULL DEFAULT '',

    -- Raw provider type (Google, OSM, etc.)
    external_type VARCHAR(100),

    -- Coordinates
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,

    -- Geography (auto-generated from lat/lng)
    geog GEOGRAPHY(Point, 4326)
    GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography
    ) STORED,

    -- Hierarchy
    parent_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,

    -- Path using ltree (e.g. '1.5.10')
    path LTREE,

    -- Data provider
    provider VARCHAR(50) DEFAULT 'google',

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Fast lookup by external_id
CREATE INDEX IF NOT EXISTS idx_locations_external_id
ON locations (external_id);

-- Spatial index (VERY important)
CREATE INDEX IF NOT EXISTS idx_locations_geog
ON locations USING GIST (geog);

-- Hierarchy index
CREATE INDEX IF NOT EXISTS idx_locations_path
ON locations USING GIST (path);

-- Parent lookup
CREATE INDEX IF NOT EXISTS idx_locations_parent_id
ON locations (parent_id);

-- =========================
-- 2. LOCATION TRANSLATIONS
-- =========================
CREATE TABLE IF NOT EXISTS location_translations (
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    lang_code VARCHAR(5) NOT NULL,
    name TEXT NOT NULL,
    formatted_address TEXT,
    short_formatted_address TEXT,
    PRIMARY KEY (location_id, lang_code)
);

CREATE INDEX IF NOT EXISTS idx_location_translations_name 
ON location_translations 
USING gin(to_tsvector('simple', name));

-- =========================
-- 3. LOCATION ALIAS
-- =========================
CREATE TABLE IF NOT EXISTS location_aliases (
    id SERIAL PRIMARY KEY,
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    alias TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_location_aliases_alias 
ON location_aliases(LOWER(alias));

-- =========================
-- 4. POSTS
-- =========================
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL DEFAULT 1,
    content TEXT NOT NULL,
    media_type VARCHAR(20) DEFAULT 'text',  -- text / photo / video
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_location_id 
ON posts(location_id);

CREATE INDEX IF NOT EXISTS idx_posts_created_at 
ON posts(created_at DESC);

-- =========================
-- 5. LOCATION STATS
-- =========================
CREATE TABLE IF NOT EXISTS location_stats (
    location_id INTEGER PRIMARY KEY REFERENCES locations(id) ON DELETE CASCADE,
    total_posts BIGINT DEFAULT 0,
    total_photos BIGINT DEFAULT 0,
    total_videos BIGINT DEFAULT 0,
    last_post_at TIMESTAMPTZ,
    trending_score DOUBLE PRECISION DEFAULT 0
);

-- =========================
-- 6. TRENDING LOCATIONS
-- =========================
CREATE TABLE IF NOT EXISTS trending_locations (
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (location_id, date)
);

CREATE INDEX IF NOT EXISTS idx_trending_locations_date 
ON trending_locations(date DESC, score DESC);