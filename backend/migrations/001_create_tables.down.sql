-- Migration: 001_create_tables (DOWN)
-- Drops all tables in reverse dependency order.

DROP TABLE IF EXISTS location_alias;
DROP TABLE IF EXISTS location_translations;
DROP TABLE IF EXISTS locations;
