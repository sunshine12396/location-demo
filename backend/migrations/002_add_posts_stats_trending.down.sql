-- Migration: 002_add_posts_stats_trending (DOWN)
DROP TABLE IF EXISTS trending_locations;
DROP TABLE IF EXISTS location_stats;
DROP TABLE IF EXISTS posts;
ALTER TABLE locations DROP COLUMN IF EXISTS slug;
ALTER TABLE locations DROP COLUMN IF EXISTS is_verified;
ALTER TABLE locations DROP COLUMN IF EXISTS updated_at;
