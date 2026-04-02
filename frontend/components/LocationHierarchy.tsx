'use client';

import { ChevronRight, MapPin } from 'lucide-react';
import { LocationDetail } from '@/lib/api';
import Link from 'next/link';

interface LocationHierarchyProps {
  location: LocationDetail;
  language: string;
}

export default function LocationHierarchy({
  location,
  language,
}: LocationHierarchyProps) {
  // Build breadcrumb from parents
  const breadcrumb: { id: number; name: string; type: string }[] = [];
  let current: any = location.parent;
  while (current) {
    breadcrumb.unshift({ id: current.id, name: current.name, type: current.type });
    current = current.parent;
  }
  // Add self
  breadcrumb.push({ id: location.id, name: location.name, type: location.type });

  const getTypeColor = (type: string) => {
    const colors: { [key: string]: string } = {
      country: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
      city: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
      district: 'bg-purple-500/10 text-purple-400 border-purple-500/20',
      landmark: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
      venue: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
    };
    return colors[type] || 'bg-slate-500/10 text-slate-400 border-slate-500/20';
  };

  const typeIcon: Record<string, string> = {
    country: "🌍",
    city: "🏙️",
    district: "📍",
    landmark: "🏛️",
    venue: "📌",
  };

  return (
    <div className="space-y-4">
      {/* Visual Hierarchy Cards */}
      <div className="flex items-center flex-wrap gap-2">
        {breadcrumb.map((item, index) => (
          <div key={item.id} className="flex items-center gap-2">
            <Link
              href={`/location/${item.id}?lang=${language}`}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium border transition-all hover:scale-105 active:scale-95 ${getTypeColor(
                item.type
              )} ${item.id === location.id ? 'ring-2 ring-blue-500 ring-offset-2 ring-offset-slate-900 border-transparent shadow-[0_0_15px_rgba(59,130,246,0.3)]' : ''}`}
            >
              <div className="flex items-center gap-1.5">
                <span className="text-base leading-none">
                  {typeIcon[item.type] || "📍"}
                </span>
                {item.name}
              </div>
            </Link>
            {index < breadcrumb.length - 1 && (
              <ChevronRight className="w-4 h-4 text-slate-600" />
            )}
          </div>
        ))}
      </div>

      {/* Structural Visual Layout */}
      <div className="bg-slate-800/20 rounded-xl p-5 border border-slate-700/30">
        <p className="text-xs text-slate-500 mb-4 font-semibold uppercase tracking-wider">
          Hierarchical Path Visualization
        </p>
        <div className="space-y-3 relative">
          {breadcrumb.map((item, index) => (
            <div key={item.id} className="relative">
              {/* Connection Line */}
              {index < breadcrumb.length - 1 && (
                <div 
                  className="absolute left-[7px] w-px bg-gradient-to-b from-blue-500/30 to-blue-500/10" 
                  style={{ 
                    top: '20px', 
                    height: '24px', 
                    marginLeft: `${index * 24}px` 
                  }} 
                />
              )}
              
              <div
                style={{ marginLeft: `${index * 24}px` }}
                className="flex items-center gap-3"
              >
                <div className={`w-[14px] h-[14px] rounded-full flex-shrink-0 z-10 border-2 border-slate-900 shadow-lg ${item.id === location.id ? 'bg-blue-500 scale-125' : 'bg-slate-600 border-slate-500'}`} />
                <div className={`px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-tight border flex-shrink-0 ${getTypeColor(item.type)}`}>
                  {item.type}
                </div>
                <span className={`text-sm font-medium truncate ${item.id === location.id ? 'text-blue-400' : 'text-slate-400'}`}>
                  {item.name}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
