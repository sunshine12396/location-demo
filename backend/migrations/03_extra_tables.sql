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
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE, =1243
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS post_location (
    post_id SERIAL PRIMARY KEY,
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