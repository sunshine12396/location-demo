import { getLocation, LocationDetail } from "@/lib/api";
import LocationStats from "@/components/LocationStats";
import LocationNav from "@/components/LocationNav";
import EmbeddedMap from "@/components/EmbeddedMap";

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

  let location: LocationDetail;
  try {
    location = await getLocation(Number(id), lang);
  } catch (err) {
    console.error(`Failed to fetch location ${id}:`, err);
    return (
      <div className="min-h-screen flex items-center justify-center px-6">
        <div className="bg-slate-900/50 backdrop-blur-xl rounded-[2.5rem] p-16 text-center max-w-lg border border-slate-700/50 shadow-2xl relative overflow-hidden">
          <div className="absolute top-0 left-0 w-full h-1 bg-linear-to-r from-transparent via-red-500/50 to-transparent" />
          <div className="text-7xl mb-8">🧭</div>
          <h1 className="text-4xl font-black text-white tracking-tight mb-6 uppercase">Node Missing</h1>
          <p className="text-slate-400 mb-12 font-medium leading-relaxed">
            The geological coordinate <span className="text-white font-bold">{id}</span> could not be synchronized with the global index. It may have been decommissioned or moved.
          </p>
          <a
            href="/"
            className="inline-flex items-center gap-3 px-8 py-4 bg-blue-600 hover:bg-blue-500 text-white rounded-2xl font-black uppercase tracking-widest text-xs transition-all shadow-lg shadow-blue-600/30"
          >
            Return to Explorer
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-24">

      {/* Simple Subtle Background */}
      <div className="fixed inset-0 -z-50 pointer-events-none">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[500px] bg-blue-600/5 blur-[120px] opacity-30" />
      </div>

      <div className="max-w-7xl mx-auto px-6 pt-12">
        <LocationNav id={id} lang={lang} active="detail" locationName={location.name} />

        <div className="mt-12 space-y-12">
          {/* Hero Map Section */}
          <EmbeddedMap
            lat={location.lat}
            lng={location.lng}
            name={location.name}
          />

          <div className="max-w-4xl mx-auto">
            <div className="grid grid-cols-1 gap-12">
              <LocationStats location={location} language={lang} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}


