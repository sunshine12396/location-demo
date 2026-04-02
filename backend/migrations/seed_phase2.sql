-- Seed Data Phase 2: Posts, Stats, Trending, and Slugs
-- Run after 002_add_posts_stats_trending.up.sql

-- ============================================================
-- 1. Add slugs to existing locations
-- ============================================================
UPDATE locations SET slug = 'vietnam'           WHERE id = 1;
UPDATE locations SET slug = 'ho-chi-minh-city'  WHERE id = 2;
UPDATE locations SET slug = 'hanoi'             WHERE id = 3;
UPDATE locations SET slug = 'district-1'        WHERE id = 4;
UPDATE locations SET slug = 'bitexco-tower'     WHERE id = 5;
UPDATE locations SET slug = 'hoan-kiem-lake'    WHERE id = 6;
UPDATE locations SET slug = 'da-nang'           WHERE id = 7;
UPDATE locations SET slug = 'japan'             WHERE id = 8;
UPDATE locations SET slug = 'tokyo'             WHERE id = 9;
UPDATE locations SET slug = 'shibuya'           WHERE id = 10;

-- Mark countries and major cities as verified
UPDATE locations SET is_verified = TRUE WHERE id IN (1, 2, 3, 7, 8, 9);

-- ============================================================
-- 2. Sample posts
-- ============================================================
INSERT INTO posts (user_id, content, media_type, location_id) VALUES
    (1, 'Beautiful morning in Ho Chi Minh City! The energy here is incredible 🏙️', 'photo', 2),
    (2, 'Street food in Saigon is the best in the world 🍜', 'photo', 2),
    (1, 'Walking through District 1, amazing architecture everywhere', 'photo', 4),
    (3, 'Bitexco Tower at sunset — breathtaking view! 🌅', 'photo', 5),
    (2, 'Bitexco skydeck vlog is up!', 'video', 5),
    (1, 'Peaceful morning at Hoan Kiem Lake 🌿', 'photo', 6),
    (4, 'Old Quarter vibes near Hoan Kiem 🏮', 'photo', 6),
    (3, 'Exploring Hanoi — the capital city has so much history', 'text', 3),
    (1, 'Ha Long Bay day trip from Hanoi was worth it!', 'photo', 3),
    (2, 'Da Nang beach at night is magical ✨', 'photo', 7),
    (4, 'Dragon Bridge fire show in Da Nang!', 'video', 7),
    (1, 'First time in Tokyo — Shibuya crossing is insane! 🚶', 'photo', 9),
    (3, 'Tokyo street food tour — takoyaki is life 🐙', 'video', 9),
    (2, 'Shibuya at night, neon everywhere 🌃', 'photo', 10),
    (1, 'Lost in the backstreets of Shibuya', 'text', 10),
    (4, 'Cherry blossom season in Japan 🌸', 'photo', 8),
    (1, 'Vietnam is absolutely beautiful!', 'text', 1),
    (2, 'Vietnam travel diary — 2 weeks across the country', 'video', 1),
    (3, 'Coffee culture in Saigon is next level ☕', 'photo', 2),
    (1, 'Rainy season in Ho Chi Minh City has its own charm 🌧️', 'photo', 2);

-- ============================================================
-- 3. Initialize location_stats from actual post data
-- ============================================================
INSERT INTO location_stats (location_id, total_posts, total_photos, total_videos, last_post_at, trending_score)
SELECT
    p.location_id,
    COUNT(*) AS total_posts,
    COUNT(*) FILTER (WHERE p.media_type = 'photo') AS total_photos,
    COUNT(*) FILTER (WHERE p.media_type = 'video') AS total_videos,
    MAX(p.created_at) AS last_post_at,
    -- Simple trending score formula
    (COUNT(*) * 1.0) +
    (COUNT(*) FILTER (WHERE p.media_type = 'photo') * 1.5) +
    (COUNT(*) FILTER (WHERE p.media_type = 'video') * 2.0) AS trending_score
FROM posts p
GROUP BY p.location_id
ON CONFLICT (location_id) DO UPDATE SET
    total_posts = EXCLUDED.total_posts,
    total_photos = EXCLUDED.total_photos,
    total_videos = EXCLUDED.total_videos,
    last_post_at = EXCLUDED.last_post_at,
    trending_score = EXCLUDED.trending_score;

-- ============================================================
-- 4. Insert today's trending snapshot
-- ============================================================
INSERT INTO trending_locations (location_id, score, date)
SELECT location_id, trending_score, CURRENT_DATE
FROM location_stats
WHERE trending_score > 0
ON CONFLICT (location_id, date) DO UPDATE SET score = EXCLUDED.score;
