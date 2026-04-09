'use client';

import { useState, useEffect } from 'react';
import { Search, MapPin, TrendingUp, Globe, Zap, Compass, Share2, MessageSquare, Loader2, Layers, ChevronRight, LayoutGrid, Filter, X } from 'lucide-react';
import { hydrateLocation, getPosts, getPostsByLocation, Post } from '@/lib/api';
import LocationSearch from '@/components/LocationSearch';
import TrendingLocations from '@/components/TrendingLocations';
import PostCreator from '@/components/PostCreator';
import PostFeed from '@/components/PostFeed';
import { cn } from '@/lib/utils';

type Language = 'en' | 'vi' | 'ja';

export default function Home() {
  const [language, setLanguage] = useState<Language>('en');
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoadingPosts, setIsLoadingPosts] = useState(true);
  const [filterLoc, setFilterLoc] = useState<{ id: number, name: string } | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  useEffect(() => {
    fetchPosts();
  }, [language, filterLoc]);

  const fetchPosts = async () => {
    setIsLoadingPosts(true);
    try {
      const data = await getPosts(language, 20, 0, filterLoc?.id);
      setPosts(data);
    } catch (err) {
      console.error('Failed to fetch posts:', err);
    } finally {
      setIsLoadingPosts(false);
    }
  };

  const handlePostSuccess = () => {
    fetchPosts();
    setRefreshTrigger(prev => prev + 1); // Triggers Trending widget reload
  };

  const handleSelectLocation = async (id: number, externalId?: string) => {
    if (id > 0) {
      window.location.href = `/location/${id}?lang=${language}`;
    } else if (externalId) {
      try {
        const localId = await hydrateLocation(externalId, language);
        window.location.href = `/location/${localId}?lang=${language}`;
      } catch (err) {
        console.error('Immediate hydration failed:', err);
      }
    }
  };

  const languages: { code: Language; label: string; icon: string }[] = [
    { code: 'en', label: 'English', icon: '🇺🇸' },
    { code: 'vi', label: 'Tiếng Việt', icon: '🇻🇳' },
    { code: 'ja', label: '日本語', icon: '🇯🇵' },
  ];

  return (
    <main className="min-h-screen bg-slate-950 pb-24">
      {/* Background Decorative Elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none -z-10">
        <div className="absolute top-0 left-1/4 w-[500px] h-[500px] bg-blue-600/10 blur-[120px] rounded-full opacity-30 animate-pulse" />
        <div className="absolute bottom-0 right-1/4 w-[400px] h-[400px] bg-purple-600/10 blur-[120px] rounded-full opacity-20" />
      </div>

      <div className="max-w-6xl mx-auto px-6 pt-24">
        {/* Header Section */}
        <div className="text-center mb-16">
          <h1 className="text-5xl md:text-7xl font-black text-white tracking-tight mb-6">
            Discovery <span className="text-blue-500">Explorer</span>
          </h1>
          <p className="text-slate-400 text-lg font-medium max-w-xl mx-auto opacity-80">
            Navigate and discover locations with recursive hierarchies and cross-location social feeds.
          </p>
        </div>

        {/* Action Bar */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">

          {/* Left Column: Search & Trending */}
          <div className="lg:col-span-5 space-y-12">
            <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/80 rounded-[2.5rem] p-8 shadow-2xl">
              <div className="flex items-center gap-3 mb-6 px-2">
                <Compass className="w-5 h-5 text-blue-500" />
                <h2 className="text-sm font-black text-slate-200 tracking-widest uppercase">Global Navigation</h2>
              </div>

              <LocationSearch
                language={language}
                useAutocomplete={true}
                customOnSelect={handleSelectLocation}
              />

              <div className="flex flex-wrap items-center justify-center gap-2 mt-8">
                {languages.map((lang) => (
                  <button
                    key={lang.code}
                    onClick={() => setLanguage(lang.code)}
                    className={cn(
                      "px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all",
                      language === lang.code
                        ? "bg-blue-600 text-white shadow-lg shadow-blue-500/20"
                        : "text-slate-500 hover:text-slate-300 hover:bg-slate-800/50 border border-slate-700/50"
                    )}
                  >
                    {lang.code}
                  </button>
                ))}
              </div>
            </div>

            <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/80 rounded-[2.5rem] p-8 shadow-2xl">
              <div className="flex items-center gap-3 mb-8 px-2">
                <TrendingUp className="w-5 h-5 text-amber-500" />
                <h2 className="text-sm font-black text-slate-200 tracking-widest uppercase">Popular Nodes</h2>
              </div>
              <TrendingLocations language={language} refreshTrigger={refreshTrigger} customOnSelect={handleSelectLocation} />
            </div>
          </div>

          {/* Right Column: Feed & Post Creator */}
          <div className="lg:col-span-7 space-y-8">
            {/* Create Post Card */}
            <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/80 rounded-[2.5rem] p-8 shadow-2xl relative">
              <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500/5 blur-3xl rounded-full pointer-events-none" />
              <div className="flex items-center gap-3 mb-6 px-2">
                <MessageSquare className="w-5 h-5 text-emerald-500" />
                <h2 className="text-sm font-black text-slate-200 tracking-widest uppercase">Community Pulse</h2>
              </div>
              <PostCreator language={language} onSuccess={handlePostSuccess} />
            </div>

            {/* Feed Section */}
            <div className="space-y-6">
              <div className="flex items-center justify-between px-4">
                <div className="flex items-center gap-3">
                  <Zap className="w-4 h-4 text-blue-400" />
                  <h3 className="text-xs font-black text-slate-400 tracking-[0.2em] uppercase">
                    {filterLoc ? `Feed: ${filterLoc.name}` : 'Global Live Stream'}
                  </h3>
                </div>
                {filterLoc && (
                  <button
                    onClick={() => setFilterLoc(null)}
                    className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-red-500/10 text-red-400 text-[10px] font-black uppercase tracking-widest hover:bg-red-500/20 transition-all border border-red-500/20"
                  >
                    <X className="w-3 h-3" />
                    Reset
                  </button>
                )}
              </div>

              {/* Feed Filter (Quick Pick) */}
              <div className="flex items-center gap-2 px-4 pb-2">
                <Filter className="w-3.5 h-3.5 text-slate-600" />
                <span className="text-[10px] font-black uppercase tracking-widest text-slate-600 mr-2">Filter Node:</span>
                <div className="flex-1 max-w-[200px]">
                  <LocationSearch
                    language={language}
                    useAutocomplete={false}
                    filterType="city"
                    placeholder="Search city..."
                    className="py-1.5! pr-3! pl-8! text-[10px]! rounded-lg!"
                    customOnSelect={(id, ext, name) => id > 0 && setFilterLoc({ id, name: name || "Unknown" })}
                  />
                </div>
              </div>

              {isLoadingPosts ? (
                <div className="flex flex-col items-center justify-center py-20 opacity-40">
                  <Loader2 className="w-8 h-8 animate-spin text-blue-500 mb-4" />
                  <span className="text-[10px] font-black uppercase tracking-widest text-slate-600">Syncing Stream...</span>
                </div>
              ) : (
                <PostFeed posts={posts} language={language} />
              )}
            </div>
          </div>

        </div>
      </div>
    </main>
  );
}
