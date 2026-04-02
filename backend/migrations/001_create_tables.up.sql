CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(255) UNIQUE,
    type VARCHAR(50) NOT NULL DEFAULT 'city',
    lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    lng DOUBLE PRECISION NOT NULL DEFAULT 0,
    parent_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
    path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_locations_parent_id ON locations(parent_id);
CREATE INDEX IF NOT EXISTS idx_locations_path ON locations(path);

CREATE TABLE IF NOT EXISTS location_translations (
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    lang_code VARCHAR(5) NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (location_id, lang_code)
);

CREATE INDEX IF NOT EXISTS idx_location_translations_name ON location_translations USING gin(to_tsvector('simple', name));

CREATE TABLE IF NOT EXISTS location_alias (
    id SERIAL PRIMARY KEY,
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    alias TEXT NOT NULL
);

-- Index for alias search
CREATE INDEX IF NOT EXISTS idx_location_alias_alias ON location_alias(LOWER(alias));
