-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS ltree;
-- =========================
-- LOCATIONS
-- =========================
CREATE TABLE IF NOT EXISTS locations (
    id BIGSERIAL PRIMARY KEY,

    code TEXT,

    -- Internal classification
    type TEXT NOT NULL DEFAULT '',

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
    provider TEXT DEFAULT 'google',
    -- External reference (e.g. Google Place ID)
    external_id TEXT UNIQUE,
    -- Raw provider type (Google, OSM, etc.)
    external_type TEXT,
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
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
-- LOCATION TRANSLATIONS
-- =========================
CREATE TABLE IF NOT EXISTS location_translations (
    location_id BIGINT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    lang_code TEXT NOT NULL,
    name TEXT NOT NULL,
    formatted_address TEXT,
    short_formatted_address TEXT,
    PRIMARY KEY (location_id, lang_code)
);

CREATE INDEX IF NOT EXISTS idx_location_translations_name 
ON location_translations 
USING gin(to_tsvector('simple', name));
