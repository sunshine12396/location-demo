import json
import os
import unicodedata

def remove_accents(input_str: str) -> str:
    if not input_str:
        return input_str

    # Handle Vietnamese-specific characters that are NOT decomposed
    input_str = input_str.replace('đ', 'd').replace('Đ', 'D')

    # Normalize Unicode (decompose characters into base + diacritics)
    nfkd_form = unicodedata.normalize('NFKD', input_str)

    # Remove combining characters (diacritics)
    result = "".join(
        c for c in nfkd_form if not unicodedata.combining(c)
    )

    return result


# SOURCE: 
# https://github.com/dr5hn/countries-states-cities-database
def migrate():
    # Base directory is one level up from the script's directory (backend/)
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    countries_json_path = os.path.join(base_dir, "resources", "countries.json")
    states_json_path = os.path.join(base_dir, "resources", "states.json")
    output_sql_path = os.path.join(base_dir, "migrations", "02_seed_locations.sql")

    if not os.path.exists(countries_json_path) or not os.path.exists(states_json_path):
        print(f"Source JSON files not found at: {countries_json_path}")
        return

    with open(countries_json_path, 'r', encoding='utf-8') as f:
        countries_data = json.load(f)
    
    with open(states_json_path, 'r', encoding='utf-8') as f:
        states_data = json.load(f)

    print(f"Loaded {len(countries_data)} countries and {len(states_data)} states.")

    # Build maps for hierarchy resolution
    country_pg_id_map = {} # mysql_country_id -> pg_id
    pg_id_to_data = {}     # pg_id -> {'path': str, 'parent_id': int}
    
    current_pg_id = 1
    
    # Pre-process countries to assign PG IDs
    for c in countries_data:
        pg_id = current_pg_id
        country_pg_id_map[c['id']] = pg_id
        pg_id_to_data[pg_id] = {'path': str(pg_id), 'parent_id': None}
        c['pg_id'] = pg_id
        current_pg_id += 1

    # Pre-process states to build lookup and assign PG IDs
    states_by_id = {s['id']: s for s in states_data}
    resolved_states = {} # mysql_id -> pg_id
    pg_id_to_state = {}

    def resolve_state(s_id):
        if s_id in resolved_states:
            return resolved_states[s_id]
        
        s = states_by_id[s_id]
        nonlocal current_pg_id
        pg_id = current_pg_id
        current_pg_id += 1
        resolved_states[s_id] = pg_id
        s['pg_id'] = pg_id
        pg_id_to_state[pg_id] = s
        
        # Initialize default
        pg_id_to_data[pg_id] = {'path': str(pg_id), 'parent_id': None}
        
        # Determine parent and path
        parent_mysql_id = s.get('parent_id')
        if parent_mysql_id:
            try:
                parent_mysql_id = int(parent_mysql_id)
            except:
                parent_mysql_id = None

        if parent_mysql_id and parent_mysql_id in states_by_id:
            parent_pg_id = resolve_state(parent_mysql_id)
            parent_path = pg_id_to_data[parent_pg_id]['path']
            pg_id_to_data[pg_id] = {'path': f"{parent_path}.{pg_id}", 'parent_id': parent_pg_id}
        else:
            # Fallback to country
            country_mysql_id = s.get('country_id')
            if country_mysql_id:
                try:
                    country_mysql_id = int(country_mysql_id)
                except:
                    country_mysql_id = None
                    
            if country_mysql_id and country_mysql_id in country_pg_id_map:
                parent_pg_id = country_pg_id_map[country_mysql_id]
                parent_path = pg_id_to_data[parent_pg_id]['path']
                pg_id_to_data[pg_id] = {'path': f"{parent_path}.{pg_id}", 'parent_id': parent_pg_id}
        
        return pg_id

    for s_id in states_by_id:
        resolve_state(s_id)

    with open(output_sql_path, 'w', encoding='utf-8') as out:
        out.write("-- Seed data for locations and translations (Generated from JSON sources)\n")
        out.write("BEGIN;\n\n")

        out.write("-- 1. Countries\n")
        for c in countries_data:
            pg_id = c['pg_id']
            lat = c.get('latitude', 0.0)
            lng = c.get('longitude', 0.0)
            code = c.get('iso2', '')
            path = pg_id_to_data[pg_id]['path']
            
            out.write(f"INSERT INTO locations (id, code, external_id, type, external_type, lat, lng, path, provider) VALUES ({pg_id}, '{code}', 'country_{c['id']}', 'country', 'country', {lat}, {lng}, '{path}', 'local') ON CONFLICT (external_id) DO NOTHING;\n")
            
            # Translations
            trans = c.get('translations', {})
            # Add native field as a translation
            if c.get('native'):
                trans['native'] = c['native']
            
            # Ensure English name is present
            if 'en' not in trans:
                trans['en'] = c['name']
            
            trans_values = []
            for lang, name in trans.items():
                if not name: continue
                name_esc = name.replace("'", "''")
                trans_values.append(f"({pg_id}, '{lang}', '{name_esc}')")
            
            if trans_values:
                out.write(f"INSERT INTO location_translations (location_id, lang_code, name) VALUES {', '.join(trans_values)} ON CONFLICT DO NOTHING;\n")

        out.write("\n-- 2. States (Cities)\n")
        # Sort by path length to ensure parents are inserted first
        sorted_state_pg_ids = sorted(resolved_states.values(), key=lambda x: pg_id_to_data[x]['path'].count('.'))

        for pg_id in sorted_state_pg_ids:
            s = pg_id_to_state[pg_id]
            lat = s.get('latitude', 0.0)
            lng = s.get('longitude', 0.0)
            code = s.get('iso2', '')
            path = pg_id_to_data[pg_id]['path']
            parent_id = pg_id_to_data[pg_id]['parent_id']
            parent_id_str = str(parent_id) if parent_id else 'NULL'
            
            external_type = s.get('type', 'state')
            out.write(f"INSERT INTO locations (id, code, external_id, type, external_type, lat, lng, parent_id, path, provider) VALUES ({pg_id}, '{code}', 'state_{s['id']}', 'city', '{external_type}', {lat}, {lng}, {parent_id_str}, '{path}', 'local') ON CONFLICT (external_id) DO NOTHING;\n")
            
            # Translations
            trans = s.get('translations', {})
            if s.get('native'):
                trans['native'] = s['native']
            
            if 'en' not in trans:
                trans['en'] = remove_accents(s['name'])
            
            trans_values = []
            for lang, name in trans.items():
                if not name: continue
                name_esc = name.replace("'", "''")
                trans_values.append(f"({pg_id}, '{lang}', '{name_esc}')")
            
            if trans_values:
                out.write(f"INSERT INTO location_translations (location_id, lang_code, name) VALUES {', '.join(trans_values)} ON CONFLICT DO NOTHING;\n")

        out.write(f"\nSELECT setval('locations_id_seq', (SELECT MAX(id) FROM locations));\n")
        out.write("\nCOMMIT;\n")

if __name__ == "__main__":
    migrate()
