import { getPostsByLocation, getLocation, Post } from "@/lib/api";
import PostCreator from "@/components/PostCreator";
import PostFeed from "@/components/PostFeed";
import { MessageSquare, Share2, MessageCircle } from 'lucide-react';
import LocationNav from "@/components/LocationNav";

interface PageProps {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ lang?: string }>;
}

export default async function CommunityPage({
  params,
  searchParams,
}: PageProps) {
  const { id } = await params;
  const { lang = "en" } = await searchParams;

  const location = await getLocation(Number(id), lang);
  const posts: Post[] = await getPostsByLocation(Number(id), lang);

  return (
    <div className="min-h-screen pb-24">
      <div className="max-w-7xl mx-auto px-6 pt-12">
        <LocationNav id={id} lang={lang} active="community" locationName={location.name} />

        <div className="mt-12 grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">
           
           {/* Feed Column (8 cols) */}
           <div className="lg:col-span-8 space-y-12">
              <div className="flex items-center justify-between mb-8 px-4">
                 <div className="flex items-center gap-4">
                    <MessageSquare className="w-8 h-8 text-blue-500" />
                    <h2 className="text-3xl font-black text-white tracking-tighter uppercase">Community Stream</h2>
                 </div>
                 <div className="flex items-center gap-3 p-1.5 bg-slate-800/40 rounded-xl border border-slate-700/50">
                    <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse ml-1" />
                    <span className="text-[10px] font-black tracking-widest text-slate-500 mr-2">{posts.length} TRANSMISSIONS</span>
                 </div>
              </div>
              
              <PostFeed posts={posts} language={lang} />
           </div>

           {/* Engagement Sidebar (4 cols) */}
           <div className="lg:col-span-4 lg:sticky lg:top-12 space-y-8">
              <div className="bg-slate-800/10 rounded-[2.5rem] p-10 border border-slate-700/30 shadow-2xl relative overflow-hidden backdrop-blur-sm">
                  <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-blue-500/50 to-transparent" />
                  <h2 className="text-xl font-black text-white tracking-widest mb-8 flex items-center gap-3">
                     <Share2 className="w-6 h-6 text-blue-500" />
                     NEXUS FEED
                  </h2>
                  <PostCreator 
                    locationId={Number(id)} 
                    locationName={location.name} 
                    language={lang} 
                  />
                  <p className="text-[10px] text-slate-500 font-bold uppercase tracking-[0.2em] mt-8 text-center leading-relaxed">
                     Broadcasting from coordinate <span className="text-blue-400 font-black">{location.name}</span>
                  </p>
              </div>

              <div className="p-8 rounded-3xl bg-slate-800/20 border border-slate-700/30">
                  <h3 className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2 mb-6">
                     <MessageCircle className="w-4 h-4" />
                     Social Protocol
                  </h3>
                  <div className="space-y-4">
                     {[
                        { title: 'Public Ledger', desc: 'All posts are recorded to the node permanently.' },
                        { title: 'Global Reach', desc: 'Translations are synthesized automatically.' },
                        { title: 'Direct Access', desc: 'Connect directly to regional geological nodes.' }
                     ].map((p, i) => (
                        <div key={i} className="flex flex-col gap-1 px-2 border-l-2 border-slate-700 hover:border-blue-500 transition-colors pl-4">
                           <h4 className="text-xs font-black text-white uppercase tracking-tight">{p.title}</h4>
                           <p className="text-[10px] text-slate-500 font-medium leading-relaxed">{p.desc}</p>
                        </div>
                     ))}
                  </div>
              </div>
           </div>

        </div>
      </div>
    </div>
  );
}
