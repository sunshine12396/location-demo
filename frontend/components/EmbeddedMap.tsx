'use client';

import { MapPin, Maximize2, ExternalLink } from 'lucide-react';

interface EmbeddedMapProps {
  lat: number;
  lng: number;
  name: string;
}

export default function EmbeddedMap({ lat, lng, name }: EmbeddedMapProps) {
  const mapUrl = `https://maps.google.com/maps?q=${lat},${lng}&z=15&output=embed`;
  const externalUrl = `https://www.google.com/maps/place/${lat},${lng}`;

  return (
    <div className="group relative w-full h-[400px] bg-slate-900 rounded-[2.5rem] overflow-hidden border border-slate-700/50 shadow-2xl">
      {/* Map Iframe */}
      <iframe
        title={`Map of ${name}`}
        width="100%"
        height="100%"
        style={{ border: 0, filter: 'grayscale(0.2) contrast(1.1) brightness(0.9)' }}
        src={mapUrl}
        allowFullScreen
        loading="lazy"
        referrerPolicy="no-referrer-when-downgrade"
        className="opacity-90 group-hover:opacity-100 transition-opacity duration-700"
      />

      {/* Overlay UI */}
      <div className="absolute top-6 left-6 right-6 flex items-start justify-between pointer-events-none">
        <div className="flex items-center gap-3 px-5 py-2.5 bg-slate-900/80 backdrop-blur-xl rounded-2xl border border-slate-700/50 shadow-2xl pointer-events-auto">
          <div className="w-8 h-8 rounded-xl bg-blue-600 flex items-center justify-center shadow-lg shadow-blue-600/40">
            <MapPin className="w-4 h-4 text-white" />
          </div>
          <div>
            <p className="text-[10px] font-black text-blue-400 uppercase tracking-widest leading-none mb-1">Geological Fix</p>
            <p className="text-sm font-black text-white tracking-tight">{lat.toFixed(4)}, {lng.toFixed(4)}</p>
          </div>
        </div>

        <a
          href={externalUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="p-3 bg-slate-900/80 backdrop-blur-xl rounded-2xl border border-slate-700/50 shadow-2xl hover:bg-slate-800 transition-all pointer-events-auto group/btn"
        >
          <ExternalLink className="w-4 h-4 text-slate-400 group-hover/btn:text-blue-400 transition-colors" />
        </a>
      </div>

      {/* Bottom Label */}
      <div className="absolute bottom-6 left-6 pointer-events-none">
         <div className="px-4 py-1.5 bg-blue-600 rounded-full shadow-lg shadow-blue-600/20 pointer-events-auto flex items-center gap-2">
            <Maximize2 className="w-3 h-3 text-white" />
            <span className="text-[9px] font-black text-white uppercase tracking-widest">Interactive Terminal</span>
         </div>
      </div>

      {/* Gradient Vignette */}
      <div className="absolute inset-0 pointer-events-none ring-1 ring-inset ring-white/10 rounded-[2.5rem]" />
    </div>
  );
}
