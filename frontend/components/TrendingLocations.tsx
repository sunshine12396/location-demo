'use client';

import { useState, useEffect } from 'react';
import { TrendingUp, MapPin, Loader2 } from 'lucide-react';
import { getTrending, TrendingLocation } from '@/lib/api';
import Link from 'next/link';
import { cn } from '@/lib/utils';


interface TrendingLocationsProps {
  language: string;
  customOnSelect?: (id: number) => void;
}

export default function TrendingLocations({
  language,
  customOnSelect,
}: TrendingLocationsProps) {

  const [trending, setTrending] = useState<TrendingLocation[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTrending = async () => {
      try {
        const data = await getTrending(language, 5);
        setTrending(data);
      } catch (err) {
        console.error('Failed to fetch trending:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchTrending();
  }, [language]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-12 gap-3 text-slate-500">
        <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
        <p className="text-sm font-medium animate-pulse uppercase tracking-widest text-[10px]">Loading Trending Alpha...</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {trending.map((location, idx) => (
        <div
          key={location.location_id}
          onClick={() => {
            if (customOnSelect) {
              customOnSelect(location.location_id);
            }
          }}
          className={cn(
            "block group relative overflow-hidden cursor-pointer",
            !customOnSelect && "pointer-events-none" // Fallback if no select provided but rendered
          )}
        >
          {/* If no customOnSelect, wrap in Link (original behavior) */}
          {!customOnSelect ? (
            <Link href={`/location/${location.location_id}?lang=${language}`} className="absolute inset-0 z-20" />
          ) : null}

          <div className="relative z-10 w-full text-left bg-slate-800/20 hover:bg-slate-800/40 border border-slate-700/50 hover:border-blue-500/30 rounded-2xl p-4 transition-all duration-300">

            <div className="flex items-start justify-between gap-4">
              <div className="flex items-start gap-4 flex-1 min-w-0">
                {/* Rank Badge */}
                <div
                  className={`flex-shrink-0 w-10 h-10 rounded-xl flex flex-col items-center justify-center font-black text-sm shadow-lg ${
                    idx === 0
                      ? 'bg-gradient-to-br from-amber-400 to-orange-500 text-white shadow-amber-900/20'
                      : idx === 1
                        ? 'bg-gradient-to-br from-slate-300 to-slate-400 text-slate-900 shadow-slate-900/20'
                        : idx === 2
                          ? 'bg-gradient-to-br from-orange-600 to-orange-700 text-white shadow-orange-900/20'
                          : 'bg-slate-700 text-slate-200 border border-slate-600 border-dashed animate-none'
                  }`}
                >
                  <span className="text-[10px] leading-none mb-0.5 opacity-60">#</span>
                  <span className="leading-none">{idx + 1}</span>
                </div>

                <div className="flex-1 min-w-0 flex flex-col justify-center h-10">
                  <div className="flex items-center gap-2">
                    <h3 className="font-bold text-white group-hover:text-blue-400 transition-colors truncate text-base">
                      {location.name}
                    </h3>
                  </div>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="text-[8px] font-black tracking-widest text-blue-500 uppercase px-1.5 py-0.5 bg-blue-500/10 rounded-sm">
                      {location.type}
                    </span>
                    <span className="text-[10px] text-slate-500 flex items-center gap-1 font-medium">
                       <div className="w-1 h-1 rounded-full bg-slate-600" />
                       REALTIME PULSE
                    </span>
                  </div>
                </div>
              </div>

              {/* Trending Score */}
              <div className="flex-shrink-0 text-right h-10 flex flex-col justify-center">
                <div className="flex items-center gap-1.5 justify-end mb-0.5">
                  <TrendingUp className={`w-3.5 h-3.5 ${idx === 0 ? 'text-amber-400 animate-bounce' : 'text-slate-500'}`} />
                  <span className={`text-xl font-black ${idx === 0 ? 'text-amber-400' : 'text-slate-200'}`}>{location.score.toFixed(1)}</span>
                </div>
                <p className="text-[9px] font-black uppercase tracking-widest text-slate-600 leading-none">HOT SCORE</p>
              </div>
            </div>
          </div>
          {/* Subtle Hover Glow */}
          <div className="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
        </div>
      ))}
    </div>
  );
}
