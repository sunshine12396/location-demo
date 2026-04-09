import { getChildren, getLocation } from "@/lib/api";
import Link from "next/link";
import { ChevronLeft, Compass, MapPin } from 'lucide-react';
import LocationNav from "@/components/LocationNav";

interface PageProps {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ lang?: string }>;
}

export default async function ExplorePage({
  params,
  searchParams,
}: PageProps) {
  const { id } = await params;
  const { lang = "en" } = await searchParams;

  const location = await getLocation(Number(id), lang);
  const children = await getChildren(Number(id), lang);

  return (
    <div className="min-h-screen pb-24">
      <div className="max-w-7xl mx-auto px-6 pt-12">
        <LocationNav id={id} lang={lang} active="explore" locationName={location.name} />

        <div className="mt-12">
          <div className="flex items-center gap-4 mb-10">
            <div className="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center border border-indigo-500/20">
               <Compass className="w-6 h-6 text-indigo-400" />
            </div>
            <div>
               <h2 className="text-3xl font-black text-white tracking-tighter uppercase">Sub-Coordinate Nodes</h2>
               <p className="text-slate-500 text-sm font-medium">Explore the recursive geological hierarchy of {location.name}</p>
            </div>
          </div>

          {children && children.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {children.map((child) => (
                <Link
                  key={child.id}
                  href={`/location/${child.id}?lang=${lang}`}
                  className="group p-8 rounded-[2rem] bg-slate-800/20 border border-slate-700/50 hover:border-indigo-500/30 transition-all flex flex-col gap-6 shadow-xl hover:shadow-indigo-500/5 relative overflow-hidden"
                >
                  <div className="absolute top-0 right-0 p-8 text-6xl font-black text-white/5 group-hover:scale-110 transition-transform pointer-events-none">
                    {child.type[0].toUpperCase()}
                  </div>
                  
                  <div className="w-14 h-14 rounded-2xl bg-slate-800 flex items-center justify-center group-hover:bg-indigo-500/20 group-hover:text-indigo-400 transition-colors border border-slate-700 group-hover:border-indigo-500/30">
                    <MapPin className="w-7 h-7" />
                  </div>
                  
                  <div>
                    <h4 className="text-xl font-bold text-white group-hover:text-white transition-colors mb-1">{child.name}</h4>
                    <span className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500 group-hover:text-indigo-400/70 transition-colors">{child.type}</span>
                  </div>

                  <div className="flex items-center gap-2 pt-4 border-t border-slate-700/30 mt-auto group-hover:border-indigo-500/20">
                     <span className="text-[9px] font-black uppercase tracking-widest text-slate-500">Initialize Connection</span>
                     <ChevronLeft className="w-3 h-3 text-slate-600 rotate-180 group-hover:text-indigo-400 group-hover:translate-x-1 transition-all" />
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="py-32 flex flex-col items-center justify-center bg-slate-800/10 rounded-[3rem] border border-slate-700/30 border-dashed">
               <Compass className="w-20 h-20 text-slate-700 mb-6 animate-pulse" />
               <h3 className="text-xl font-black text-slate-600 uppercase tracking-widest">No terminal nodes found</h3>
               <p className="text-slate-500 text-sm mt-2 max-w-sm text-center">This coordinate represents a leaf node in the global hierarchy matrix.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
