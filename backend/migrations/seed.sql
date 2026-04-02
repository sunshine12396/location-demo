-- Seed Data: Vietnam location hierarchy with translations and aliases
-- Run after migrations to have test data ready.

-- ============================================================
-- 1. Insert Locations (Hierarchy: Country → City → District → Landmark)
-- ============================================================
INSERT INTO locations (id, external_id, type, lat, lng, parent_id, path) VALUES
    (1, 'ext_vietnam',       'country',  16.0471,   108.2068, NULL, '1'),
    (2, 'ext_hcm',           'city',     10.7769,   106.7009, 1,    '1/2'),
    (3, 'ext_hanoi',         'city',     21.0285,   105.8542, 1,    '1/3'),
    (4, 'ext_district1',     'district', 10.7756,   106.7019, 2,    '1/2/4'),
    (5, 'ext_bitexco',       'landmark', 10.7714,   106.7043, 4,    '1/2/4/5'),
    (6, 'ext_hoan_kiem',     'landmark', 21.0288,   105.8525, 3,    '1/3/6'),
    (7, 'ext_da_nang',       'city',     16.0544,   108.2022, 1,    '1/7'),
    (8, 'ext_japan',         'country',  36.2048,   138.2529, NULL, '8'),
    (9, 'ext_tokyo',         'city',     35.6762,   139.6503, 8,    '8/9'),
    (10, 'ext_shibuya',      'district', 35.6619,   139.7041, 9,    '8/9/10')
ON CONFLICT (external_id) DO NOTHING;

-- Reset sequence to avoid conflicts
SELECT setval('locations_id_seq', (SELECT MAX(id) FROM locations));

-- ============================================================
-- 2. Insert Translations (English + Vietnamese + Japanese)
-- ============================================================
INSERT INTO location_translations (location_id, lang_code, name) VALUES
    -- Vietnam
    (1, 'en', 'Vietnam'),
    (1, 'vi', 'Việt Nam'),
    (1, 'ja', 'ベトナム'),
    -- Ho Chi Minh City
    (2, 'en', 'Ho Chi Minh City'),
    (2, 'vi', 'Thành phố Hồ Chí Minh'),
    (2, 'ja', 'ホーチミン市'),
    -- Hanoi
    (3, 'en', 'Hanoi'),
    (3, 'vi', 'Hà Nội'),
    (3, 'ja', 'ハノイ'),
    -- District 1
    (4, 'en', 'District 1'),
    (4, 'vi', 'Quận 1'),
    -- Bitexco Tower
    (5, 'en', 'Bitexco Financial Tower'),
    (5, 'vi', 'Tháp Tài chính Bitexco'),
    -- Hoan Kiem Lake
    (6, 'en', 'Hoan Kiem Lake'),
    (6, 'vi', 'Hồ Hoàn Kiếm'),
    (6, 'ja', 'ホアンキエム湖'),
    -- Da Nang
    (7, 'en', 'Da Nang'),
    (7, 'vi', 'Đà Nẵng'),
    -- Japan
    (8, 'en', 'Japan'),
    (8, 'vi', 'Nhật Bản'),
    (8, 'ja', '日本'),
    -- Tokyo
    (9, 'en', 'Tokyo'),
    (9, 'vi', 'Tokyo'),
    (9, 'ja', '東京'),
    -- Shibuya
    (10, 'en', 'Shibuya'),
    (10, 'ja', '渋谷区')
ON CONFLICT (location_id, lang_code) DO NOTHING;

-- ============================================================
-- 3. Insert Aliases (Smart Search)
-- ============================================================
INSERT INTO location_alias (location_id, alias) VALUES
    (2, 'sai gon'),
    (2, 'saigon'),
    (2, 'hcm'),
    (2, 'tp hcm'),
    (2, 'thanh pho ho chi minh'),
    (3, 'ha noi'),
    (7, 'da nang'),
    (7, 'danang'),
    (9, 'toukyou');
