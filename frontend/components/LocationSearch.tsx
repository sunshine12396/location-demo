'use client';

import { useState, useRef, useEffect, KeyboardEvent } from 'react';
import { Search, MapPin, Loader2 } from 'lucide-react';
import { searchLocations, autocompleteLocations, SearchResult } from '@/lib/api';
import { useRouter } from 'next/navigation';
import { cn } from '@/lib/utils';

interface LocationSearchProps {
  language: string;
  useAutocomplete?: boolean;
  customOnSelect?: (id: number, externalId: string | undefined, name: string, type: string) => void;
  className?: string;
  placeholder?: string;
  filterType?: string; // added
}

export default function LocationSearch({ language, useAutocomplete = false, customOnSelect, className, placeholder, filterType }: LocationSearchProps) {

  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<SearchResult[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const router = useRouter();

  // Handle click outside to close dropdown
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const triggerSearch = async (searchQuery: string) => {
    if (searchQuery.length < 2) {
      setResults([]);
      setIsOpen(false);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = useAutocomplete 
        ? await autocompleteLocations(searchQuery, language)
        : await searchLocations(searchQuery, language, filterType);
      setResults(data || []);
      setIsOpen(true);
      setSelectedIndex(-1);
    } catch (err: any) {
      console.error('Search failed:', err);
      setError(err instanceof Error ? err.message : String(err));
      setResults([]);
      setIsOpen(true);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (isOpen && selectedIndex >= 0 && selectedIndex < results.length) {
        handleSelect(results[selectedIndex]);
      } else {
        triggerSearch(query);
      }
    } else if (e.key === 'ArrowDown') {
      if (!isOpen || results.length === 0) return;
      e.preventDefault();
      setSelectedIndex((prev) => (prev < results.length - 1 ? prev + 1 : prev));
    } else if (e.key === 'ArrowUp') {
      if (!isOpen || results.length === 0) return;
      e.preventDefault();
      setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
    } else if (e.key === 'Escape') {
      setIsOpen(false);
    }
  };

  const handleSelect = (result: SearchResult) => {
    if (customOnSelect) {
      customOnSelect(result.id, result.external_id, result.name, result.type);
    } else {
      if (result.id > 0) {
        router.push(`/location/${result.id}?lang=${language}`);
      } else if (result.external_id) {
        // For home search, if it's external-only, we can't redirect to /location/0.
        // We might want a "preview" page or just force selection in parent components.
        console.warn('Selected external-only location:', result.external_id);
      }
    }
    setQuery('');
    setIsOpen(false);
    setSelectedIndex(-1);
  };


  return (
    <div className="relative" ref={wrapperRef}>
      <div className="relative">
        <Search className={cn(
          "absolute top-1/2 -translate-y-1/2 text-slate-500",
          className ? "left-2 w-3.5 h-3.5" : "left-3 w-5 h-5"
        )} />
        <input
          ref={inputRef}
          type="text"
          placeholder={placeholder || "Search locations... (try 'sai gon' or 'tokyo')"}
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            if (e.target.value.length < 2) {
              setIsOpen(false);
              setResults([]);
            }
          }}
          onFocus={() => query.length >= 2 && setIsOpen(true)}
          onKeyDown={handleKeyDown}
          className={cn(
            "w-full bg-slate-800/50 border border-slate-700/50 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all",
            className ? className : "pl-10 pr-4 py-3" // Apply default padding if no custom class is provided
          )}
        />
        {loading && (
          <Loader2 className={cn(
            "absolute top-1/2 -translate-y-1/2 w-4 h-4 text-blue-500 animate-spin",
            className ? "right-2" : "right-3 w-5 h-5"
          )} />
        )}
        {!className && (
          <div className="absolute right-12 top-1/2 -translate-y-1/2 flex items-center gap-2">
              <span className="text-[9px] font-black uppercase text-slate-600 tracking-tighter bg-slate-900 border border-slate-700 px-1.5 py-0.5 rounded-sm">
                  Engine: {useAutocomplete ? 'Google' : 'Local'}
              </span>
          </div>
        )}
      </div>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-slate-800 border border-slate-700 rounded-xl shadow-2xl z-[100] max-h-96 overflow-hidden flex flex-col ring-1 ring-white/10">
          {error ? (
            <div className="px-4 py-8 text-center bg-red-500/10 border-t border-red-500/20">
              <div className="w-10 h-10 rounded-full bg-red-500/20 flex items-center justify-center mx-auto mb-3">
                <Search className="w-5 h-5 text-red-500" />
              </div>
              <p className="text-sm text-red-400 font-medium mb-1">Search Error</p>
              <p className="text-xs text-red-500/70 max-w-[240px] mx-auto leading-relaxed">{error}</p>
              <button 
                onClick={() => triggerSearch(query)}
                className="mt-4 text-xs bg-red-500/20 hover:bg-red-500/30 text-red-400 px-3 py-1.5 rounded-lg border border-red-500/30 transition-colors"
               >
                Try Again
              </button>
            </div>
          ) : results.length > 0 ? (
            <ul className="divide-y divide-slate-700/50 overflow-y-auto">
              {results.map((location, index) => (
                <li key={location.id || location.external_id} className="scroll-m-2">
                  <button
                    onClick={() => handleSelect(location)}
                    className={cn(
                      "w-full text-left transition-colors flex items-center justify-between group",
                      index === selectedIndex ? 'bg-slate-700/80' : 'hover:bg-slate-700/80',
                      className ? "py-2 px-3" : "py-3.5 px-4"
                    )}
                  >
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <div className={cn(
                        "rounded-lg bg-slate-700 flex items-center justify-center text-blue-400 group-hover:bg-blue-500/20 group-hover:text-blue-300 transition-colors",
                        className ? "w-6 h-6" : "w-8 h-8"
                      )}>
                        <MapPin className={className ? "w-3 h-3" : "w-4 h-4"} />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className={cn("text-white font-medium line-clamp-2 leading-snug", className ? "text-xs" : "text-sm")}>{location.name}</p>
                        <p className="text-xs text-slate-400 capitalize mt-1">{location.type}</p>
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
