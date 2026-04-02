'use client';

import { useState, useRef, useEffect } from 'react';
import { Search, MapPin, Loader2 } from 'lucide-react';
import { searchLocations, SearchResult } from '@/lib/api';
import { useRouter } from 'next/navigation';

interface LocationSearchProps {
  language: string;
  customOnSelect?: (id: number) => void;
}

export default function LocationSearch({ language, customOnSelect }: LocationSearchProps) {

  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<SearchResult[]>([]);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  useEffect(() => {
    if (query.length < 2) {
      setResults([]);
      setIsOpen(false);
      return;
    }

    const timeout = setTimeout(async () => {
      setLoading(true);
      try {
        const data = await searchLocations(query, language);
        setResults(data);
        setIsOpen(true);
      } catch (err) {
        console.error('Search failed:', err);
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => clearTimeout(timeout);
  }, [query, language]);

  const handleSelect = (id: number) => {
    if (customOnSelect) {
      customOnSelect(id);
    } else {
      router.push(`/location/${id}?lang=${language}`);
    }
    setQuery('');
    setIsOpen(false);
  };


  return (
    <div className="relative">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" />
        <input
          ref={inputRef}
          type="text"
          placeholder="Search locations... (try 'sai gon' or 'tokyo')"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => query.length >= 2 && setIsOpen(true)}
          className="w-full pl-10 pr-4 py-3 bg-slate-800/50 border border-slate-700/50 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all"
        />
        {loading && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-blue-500 animate-spin" />
        )}
      </div>

      {/* Results Dropdown */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-slate-800 border border-slate-700 rounded-xl shadow-2xl z-[100] max-h-96 overflow-hidden flex flex-col ring-1 ring-white/10">
          {results.length > 0 ? (
            <ul className="divide-y divide-slate-700/50 overflow-y-auto">
              {results.map((location) => (
                <li key={location.id}>
                  <button
                    onClick={() => handleSelect(location.id)}
                    className="w-full text-left px-4 py-3.5 hover:bg-slate-700/80 transition-colors flex items-center justify-between group"
                  >
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <div className="w-8 h-8 rounded-lg bg-slate-700 flex items-center justify-center text-blue-400 group-hover:bg-blue-500/20 group-hover:text-blue-300 transition-colors">
                        <MapPin className="w-4 h-4" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-white font-medium truncate">{location.name}</p>
                        <p className="text-xs text-slate-400 capitalize">{location.type}</p>
                      </div>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          ) : (
            <div className="px-4 py-8 text-center text-slate-400">
              <MapPin className="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">No locations found for "{query}"</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
