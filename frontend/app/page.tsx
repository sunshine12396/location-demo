'use client';

import { useState, useEffect } from 'react';
import { Search, MapPin, TrendingUp, Globe, Zap, Compass, Share2, MessageSquare, Loader2 } from 'lucide-react';
import LocationSearch from '@/components/LocationSearch';
import TrendingLocations from '@/components/TrendingLocations';
import LocationStats from '@/components/LocationStats';
import PostCreator from '@/components/PostCreator';
import PostFeed from '@/components/PostFeed';
import { getLocation, getPostsByLocation, LocationDetail, Post, SearchResult } from '@/lib/api';
import { cn } from '@/lib/utils';

type Language = 'en' | 'vi' | 'ja';

export default function Home() {
  const [language, setLanguage] = useState<Language>('en');
  const [selectedLocation, setSelectedLocation] = useState<LocationDetail | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // This handles selection from search or trending
  const handleSelectLocation = async (id: number) => {
    setIsLoading(true);
    try {
      const loc = await getLocation(id, language);
      const locPosts = await getPostsByLocation(id, language);
      setSelectedLocation(loc);
      setPosts(locPosts);
      // Optional: scroll to the interaction area
      window.scrollTo({ top: 400, behavior: 'smooth' });
    } catch (err) {
      console.error("Failed to load node detail:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const languages: { code: Language; label: string; icon: string }[] = [
    { code: 'en', label: 'English', icon: '🇺🇸' },
    { code: 'vi', label: 'Tiếng Việt', icon: '🇻🇳' },
    { code: 'ja', label: '日本語', icon: '🇯🇵' },
  ];

  return (
    <main className="min-h-screen bg-transparent">
      {/* Dynamic Animated Background */}
      <div className="fixed inset-0 -z-50 overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/10 blur-[120px] rounded-full animate-pulse" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-purple-600/10 blur-[120px] rounded-full animate-pulse delay-700" />
      </div>

      {/* Hero Section */}
      <div className="relative pt-20 pb-16 px-6 overflow-hidden">
        <div className="max-w-4xl mx-auto text-center relative z-10">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 mb-6 animate-active">
            <Zap className="w-3.5 h-3.5 text-blue-400 fill-blue-400" />
            <span className="text-[10px] font-black uppercase tracking-[0.2em] text-blue-400">Phase 2 Production Live</span>
          </div>
          
          <h1 className="text-6xl md:text-8xl font-black text-white tracking-tighter mb-6 selection:bg-blue-500">
            LOCATION <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-cyan-400 to-indigo-500">ALPHA</span>
          </h1>
          
          <p className="text-lg md:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed font-medium mb-10">
            A high-performance geological infrastructure supporting recursive hierarchies, 
            multilingual aliasing, and realtime engagement telemetry.
          </p>

          {/* Search Bar Integration */}
          <div className="max-w-2xl mx-auto transform hover:scale-[1.01] transition-transform duration-500 relative z-[100]">
            {/* Custom Search that triggers local state instead of redirect */}
            <div className="relative">
               <LocationSearch language={language} customOnSelect={handleSelectLocation} />
            </div>
          </div>

          {/* Language Selector Hub */}
          <div className="flex flex-wrap items-center justify-center gap-3 mt-8">
            <div className="flex items-center gap-1.5 p-1.5 bg-slate-800/40 backdrop-blur-xl rounded-2xl border border-slate-700/50 shadow-2xl">
              {languages.map((lang) => (
                <button
                  key={lang.code}
                  onClick={() => setLanguage(lang.code)}
                  className={cn(
                    "flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-black uppercase tracking-widest transition-all",
                    language === lang.code
                      ? "bg-blue-600 text-white shadow-lg shadow-blue-500/40 translate-y-[-1px]"
                      : "text-slate-500 hover:text-slate-300 hover:bg-slate-700/50"
                  )}
                >
                  <span className="text-base leading-none">{lang.icon}</span>
                  {lang.code}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Discovery Hub */}
      <div className="max-w-7xl mx-auto px-6 pb-24">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          
          {/* Trending Pulse Column */}
          <div className="lg:col-span-4 space-y-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-amber-500/10 rounded-lg border border-amber-500/20">
                  <TrendingUp className="w-5 h-5 text-amber-500" />
                </div>
                <h2 className="text-lg font-black text-white tracking-widest uppercase">Global Pulse</h2>
              </div>
            </div>
            
            <TrendingLocations language={language} customOnSelect={handleSelectLocation} />
            
            <div className="p-5 rounded-2xl bg-gradient-to-br from-indigo-600/20 to-purple-600/20 border border-indigo-500/20 relative group overflow-hidden">
                <div className="relative z-10">
                   <p className="text-xs font-black text-indigo-300 uppercase tracking-[0.2em] mb-2">Alpha Insight</p>
                   <p className="text-sm text-slate-300 leading-relaxed font-medium">Select a location node to activate the engagement layer and view statistics.</p>
                </div>
                <Compass className="absolute right-[-10px] bottom-[-10px] w-24 h-24 text-indigo-500/10 group-hover:rotate-45 transition-transform duration-700" />
            </div>
          </div>

          {/* Interaction Column (8 cols) */}
          <div className="lg:col-span-8">
            {isLoading ? (
               <div className="min-h-[400px] flex flex-col items-center justify-center bg-slate-800/10 rounded-[2.5rem] border border-slate-700/30 border-dashed animate-pulse">
                  <Loader2 className="w-12 h-12 text-blue-500 animate-spin mb-4" />
                  <p className="text-xs font-black uppercase tracking-[0.3em] text-slate-500">Retrieving Node Telemetry...</p>
               </div>
            ) : selectedLocation ? (
               <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
                  {/* Two columns within the main column */}
                  <div className="grid grid-cols-1 md:grid-cols-5 gap-8 items-start">
                     {/* Left: Stats */}
                     <div className="md:col-span-2">
                        <LocationStats location={selectedLocation} language={language} />
                     </div>
                     {/* Right: Post Creator */}
                     <div className="md:col-span-3 space-y-8">
                        <div className="bg-slate-800/10 rounded-[2rem] p-8 border border-slate-700/30 shadow-2xl relative overflow-hidden">
                           <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-blue-500/50 to-transparent" />
                           <h2 className="text-2xl font-black text-white tracking-tighter mb-6 flex items-center gap-3">
                              <MessageSquare className="w-6 h-6 text-blue-500" />
                              ENGAGEMENT
                           </h2>
                           <PostCreator 
                              locationId={selectedLocation.id} 
                              locationName={selectedLocation.name} 
                              language={language} 
                           />
                        </div>
                        
                        <div className="px-4">
                           <div className="flex items-center justify-between mb-6">
                              <h3 className="text-xl font-black text-white tracking-tighter uppercase">Local Stream</h3>
                              <span className="text-[10px] font-black tracking-widest text-slate-500">{posts.length} NODES</span>
                           </div>
                           <PostFeed posts={posts} language={language} />
                        </div>
                     </div>
                  </div>
               </div>
            ) : (
               <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {[
                    { title: 'Waterfall Search', desc: 'Alias → Localization → OSM Fallback', icon: Search, color: 'text-blue-500', bg: 'bg-blue-500/10 shadow-blue-500/5' },
                    { title: 'Hierarchy Traversal', desc: 'Recursive parent/child graph logic', icon: MapPin, color: 'text-cyan-500', bg: 'bg-cyan-500/10 shadow-cyan-500/5' },
                    { title: 'SEO Optimized', desc: 'Global slugs & verified data verification', icon: Globe, color: 'text-purple-500', bg: 'bg-purple-500/10 shadow-purple-500/5' },
                    { title: 'Social Integration', desc: 'Post engagement & stats aggregation', icon: Share2, color: 'text-emerald-500', bg: 'bg-emerald-500/10 shadow-emerald-500/5' },
                  ].map((feature, i) => (
                    <div key={i} className={cn("group p-6 rounded-3xl border border-slate-700/50 hover:border-slate-500/50 transition-all duration-300 flex flex-col gap-4 bg-slate-800/10 shadow-2xl relative overflow-hidden")}>
                       <div className={cn("w-12 h-12 rounded-2xl flex items-center justify-center transform group-hover:scale-110 group-hover:rotate-3 transition-transform duration-500", feature.bg)}>
                          <feature.icon className={cn("w-6 h-6", feature.color)} />
                       </div>
                       <div>
                          <h3 className="text-lg font-black text-white mb-1 group-hover:text-blue-400 transition-colors uppercase tracking-tight">{feature.title}</h3>
                          <p className="text-slate-500 text-sm font-medium leading-relaxed">{feature.desc}</p>
                       </div>
                       <span className="absolute right-6 top-6 text-6xl font-black text-white/5 pointer-events-none group-hover:scale-110 transition-transform">0{i+1}</span>
                    </div>
                  ))}
                  
                  {/* Tutorial Interaction Card */}
                  <div className="md:col-span-2 p-10 rounded-[2.5rem] bg-gradient-to-br from-blue-600/20 to-purple-600/20 border border-blue-500/20 flex flex-col items-center justify-center text-center animate-in fade-in duration-1000">
                     <div className="w-16 h-16 rounded-full bg-blue-500/20 flex items-center justify-center mb-6 ring-8 ring-blue-500/5">
                        <MapPin className="w-8 h-8 text-blue-400 animate-bounce" />
                     </div>
                     <h3 className="text-2xl font-black text-white tracking-tighter mb-4 uppercase">Node Activation Required</h3>
                     <p className="text-slate-400 max-w-sm font-medium leading-relaxed">
                        Start by searching for a location node in the nexus above or selecting from the global pulse on the left to activate the engagement matrix.
                     </p>
                  </div>
               </div>
            )}
          </div>

        </div>
      </div>
    </main>
  );
}
