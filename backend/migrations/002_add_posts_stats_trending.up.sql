-- Migration: 002_add_posts_stats_trending (UP)
-- Adds posts, location_stats, trending_locations tables and slug/verified columns.

-- 1. Add slug and is_verified columns to locations
ALTER TABLE locations ADD COLUMN IF NOT EXISTS slug VARCHAR(255) UNIQUE;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE locations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Index for slug lookups
CREATE INDEX IF NOT EXISTS idx_locations_slug ON locations(slug);

-- 2. Posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL DEFAULT 1,
    content TEXT NOT NULL,
    media_type VARCHAR(20) DEFAULT 'text',  -- text / photo / video
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_location_id ON posts(location_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- 3. Location stats table (pre-aggregated counters)
CREATE TABLE IF NOT EXISTS location_stats (
    location_id INTEGER PRIMARY KEY REFERENCES locations(id) ON DELETE CASCADE,
    total_posts BIGINT DEFAULT 0,
    total_photos BIGINT DEFAULT 0,
    total_videos BIGINT DEFAULT 0,
    last_post_at TIMESTAMP WITH TIME ZONE,
    trending_score DOUBLE PRECISION DEFAULT 0
);

-- 4. Trending locations table (daily snapshots)
CREATE TABLE IF NOT EXISTS trending_locations (
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    score DOUBLE PRECISION NOT NULL DEFAULT 0,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (location_id, date)
);

CREATE INDEX IF NOT EXISTS idx_trending_locations_date ON trending_locations(date DESC, score DESC);
