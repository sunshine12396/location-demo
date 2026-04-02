import { getLocation, getChildren, getPostsByLocation, Post } from "@/lib/api";
import Link from "next/link";
import LocationHierarchy from "@/components/LocationHierarchy";
import LocationStats from "@/components/LocationStats";
import PostCreator from "@/components/PostCreator";
import PostFeed from "@/components/PostFeed";
import { ChevronLeft, Globe, MapPin, Layers, MessageSquare, Compass } from 'lucide-react';

interface PageProps {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ lang?: string }>;
}

export default async function LocationDetailPage({
  params,
  searchParams,
}: PageProps) {
  const { id } = await params;
  const { lang = "en" } = await searchParams;

  let location;
  let children;
  let posts: Post[] = [];

  try {
    location = await getLocation(Number(id), lang);
    children = await getChildren(Number(id), lang);
    posts = await getPostsByLocation(Number(id), lang);
  } catch (err) {
    console.error(err);
    return (
      <div className="min-h-screen flex items-center justify-center px-6">
        <div className="bg-slate-800/50 backdrop-blur-xl rounded-[2.5rem] p-12 text-center max-w-lg border border-slate-700/50 shadow-2xl">
          <div className="text-7xl mb-6">🏝️</div>
          <h1 className="text-4xl font-black text-white tracking-tighter mb-4">Node Disconnected</h1>
          <p className="text-slate-400 mb-10 font-medium leading-relaxed">
            The requested geological coordinate (ID: {id}) is not responding or has been decommissioned from the global index.
          </p>
          <Link
            href="/"
            className="inline-flex items-center gap-3 px-8 py-4 bg-blue-600 hover:bg-blue-500 text-white rounded-2xl font-black uppercase tracking-widest text-xs transition-all shadow-lg shadow-blue-600/30 active:scale-95"
          >
            <ChevronLeft className="w-4 h-4" />
            Return to Nexus
          </Link>
        </div>
      </div>
    );
  }

  const languages = [
    { code: "en", label: "EN", icon: "🇺🇸" },
    { code: "vi", label: "VI", icon: "🇻🇳" },
    { code: "ja", label: "JA", icon: "🇯🇵" },
  ];

  return (
    <div className="min-h-screen pb-24">
      {/* Background Decor */}
      <div className="fixed inset-0 -z-50 pointer-events-none">
         <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[500px] bg-gradient-to-b from-blue-600/5 to-transparent blur-3xl opacity-50" />
      </div>

      <div className="max-w-7xl mx-auto px-6 pt-12">
        {/* Navigation / Actions Bar */}
        <div className="flex items-center justify-between mb-12">
           <Link href="/" className="group flex items-center gap-3 px-5 py-2.5 bg-slate-800/40 hover:bg-slate-800/60 transition-all rounded-2xl border border-slate-700/50">
              <ChevronLeft className="w-4 h-4 text-slate-500 group-hover:text-blue-400 transition-colors" />
              <span className="text-[10px] font-black uppercase tracking-widest text-slate-400 group-hover:text-white transition-colors">Surface Menu</span>
           </Link>

           <div className="flex items-center gap-2 p-1.5 bg-slate-800/40 rounded-2xl border border-slate-700/50 backdrop-blur-md">
              <Globe className="w-3.5 h-3.5 text-slate-500 ml-2" />
              {languages.map((l) => (
                <Link
                  key={l.code}
                  href={`/location/${id}?lang=${l.code}`}
                  className={`px-4 py-2 rounded-xl text-[10px] font-black tracking-widest uppercase transition-all ${
                    lang === l.code
                      ? "bg-blue-600 text-white shadow-lg shadow-blue-600/40"
                      : "text-slate-500 hover:text-slate-300 hover:bg-slate-700/50"
                  }`}
                >
                  {l.label}
                </Link>
              ))}
           </div>
        </div>

        {/* Main Grid Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          
          {/* Left Column: Stats & Meta (4 cols) */}
          <div className="lg:col-span-4 lg:sticky lg:top-12 space-y-8">
            <LocationStats location={location} language={lang} />
            
            {/* Hierarchy Insight */}
            <div className="bg-slate-800/10 rounded-[2rem] p-8 border border-slate-700/30 relative overflow-hidden group">
               <Layers className="absolute right-[-20px] top-[-20px] w-40 h-40 text-slate-500/5 group-hover:rotate-12 transition-transform duration-1000" />
               <div className="relative z-10">
                  <h3 className="text-xl font-black text-white tracking-widest uppercase mb-4 flex items-center gap-2">
                     <div className="w-2 h-2 rounded-full bg-blue-500" />
                     Recursive Map
                  </h3>
                  <LocationHierarchy location={location} language={lang} />
               </div>
            </div>
          </div>

          {/* Right Column: Interaction & Feed (8 cols) */}
          <div className="lg:col-span-8 space-y-12">
            
            {/* Create Post Section */}
            <div className="bg-slate-800/10 rounded-[2.5rem] p-8 md:p-12 border border-slate-700/30 shadow-2xl relative overflow-hidden">
                <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-blue-500/50 to-transparent" />
                <h2 className="text-3xl font-black text-white tracking-tighter mb-8 flex items-center gap-4">
                   <MessageSquare className="w-8 h-8 text-blue-500" />
                   ENGAGEMENT LAYER
                </h2>
                <PostCreator 
                  locationId={Number(id)} 
                  locationName={location.name} 
                  language={lang} 
                />
            </div>

            {/* Social Feed Section */}
            <div>
              <div className="flex items-center justify-between mb-8 px-4">
                 <h2 className="text-2xl font-black text-white tracking-tighter uppercase">Community Stream</h2>
                 <div className="flex items-center gap-2 text-slate-500">
                    <span className="text-[10px] font-black tracking-widest">{posts.length} TRANSMISSIONS</span>
                 </div>
              </div>
              <PostFeed posts={posts} language={lang} />
            </div>

            {/* Children / Discovery Section */}
            {children && children.length > 0 && (
              <div className="pt-12 border-t border-slate-700/30">
                <div className="flex items-center gap-4 mb-8">
                   <Compass className="w-6 h-6 text-indigo-400" />
                   <h2 className="text-2xl font-black text-white tracking-tighter uppercase">Sub-Coordinate Nodes</h2>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {children.map((child) => (
                    <Link
                      key={child.id}
                      href={`/location/${child.id}?lang=${lang}`}
                      className="group p-6 rounded-2xl bg-slate-800/20 border border-slate-700/50 hover:border-indigo-500/30 transition-all flex items-center justify-between shadow-lg hover:shadow-indigo-500/5"
                    >
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 rounded-xl bg-slate-700 flex items-center justify-center group-hover:bg-indigo-500/20 group-hover:text-indigo-400 transition-colors">
                           <MapPin className="w-5 h-5" />
                        </div>
                        <div>
                           <h4 className="font-bold text-white group-hover:text-indigo-400 transition-colors">{child.name}</h4>
                           <span className="text-[9px] font-black uppercase tracking-widest text-slate-500">{child.type}</span>
                        </div>
                      </div>
                      <ChevronLeft className="w-4 h-4 text-slate-600 rotate-180 group-hover:text-indigo-400 group-hover:translate-x-1 transition-all" />
                    </Link>
                  ))}
                </div>
              </div>
            )}

          </div>

        </div>
      </div>
    </div>
  );
}
