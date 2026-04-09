'use client';

import Link from "next/link";
import { ChevronLeft, Globe, Map, Compass, MessageSquare, Search } from 'lucide-react';
import { cn } from "@/lib/utils";
import LocationSearch from "./LocationSearch";
import { useState } from 'react';

interface LocationNavProps {
  id: string;
  lang: string;
  active: 'detail' | 'explore' | 'community';
  locationName: string;
}

export default function LocationNav({ id, lang, active, locationName }: LocationNavProps) {
  const [showSearch, setShowSearch] = useState(false);
  const languages = [
    { code: "en", label: "EN", icon: "🇺🇸" },
    { code: "vi", label: "VI", icon: "🇻🇳" },
    { code: "ja", label: "JA", icon: "🇯🇵" },
  ];

  const navItems = [
    { id: 'detail', label: 'Monitor HUD', icon: Map, path: `/location/${id}?lang=${lang}` },
    { id: 'explore', label: 'Node Discovery', icon: Compass, path: `/location/${id}/explore?lang=${lang}` },
    { id: 'community', label: 'Community Feed', icon: MessageSquare, path: `/location/${id}/community?lang=${lang}` },
  ];

  return (
    <div className="flex flex-col md:flex-row md:items-center justify-between gap-8 mb-12">
      <div className="flex items-center gap-6">
        <Link href="/" className="group flex items-center gap-3 px-5 py-2.5 bg-slate-800/40 hover:bg-slate-800/60 transition-all rounded-2xl border border-slate-700/50 shadow-2xl">
          <ChevronLeft className="w-4 h-4 text-slate-500 group-hover:text-blue-400 transition-all" />
          <span className="text-[10px] font-black uppercase tracking-widest text-slate-400 group-hover:text-white transition-colors">Nexus Surface</span>
        </Link>
        <div className="h-10 w-px bg-slate-700/50 hidden md:block" />
        <div>
          <h1 className="text-2xl font-black text-white tracking-widest uppercase truncate max-w-[300px] mb-1">{locationName}</h1>
          <div className="flex items-center gap-2">
             <div className="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse" />
             <span className="text-[9px] font-black uppercase tracking-widest text-blue-500/70">Coordinate Lock - Active</span>
          </div>
        </div>
      </div>

      <button 
        onClick={() => setShowSearch(!showSearch)}
        className={cn(
          "p-2.5 rounded-xl border transition-all md:hidden",
          showSearch ? "bg-blue-600 border-blue-500 text-white" : "bg-slate-800/40 border-slate-700/50 text-slate-500 hover:text-slate-300"
        )}
      >
        <Search className="w-4 h-4" />
      </button>

      <div className={cn(
        "flex-1 max-w-sm transition-all duration-300",
        showSearch ? "block opacity-100" : "hidden md:block opacity-100"
      )}>
        <LocationSearch 
          language={lang} 
          useAutocomplete={false} 
        />
      </div>

      <div className="flex items-center gap-6 overflow-x-auto pb-2 md:pb-0">
        <nav className="flex items-center gap-1.5 p-1.5 bg-slate-800/40 rounded-2xl border border-slate-700/50 backdrop-blur-md">
          {navItems.map((item) => (
            <Link
              key={item.id}
              href={item.path}
              className={cn(
                "flex items-center gap-2 px-5 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all",
                active === item.id 
                  ? "bg-blue-600 text-white shadow-lg shadow-blue-500/40 translate-y-[-1px]" 
                  : "text-slate-500 hover:text-slate-300 hover:bg-slate-700/50"
              )}
            >
              <item.icon className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">{item.label}</span>
            </Link>
          ))}
        </nav>

        <div className="h-10 w-px bg-slate-700/50 hidden lg:block" />

        <div className="flex items-center gap-1.5 p-1.5 bg-slate-800/40 rounded-2xl border border-slate-700/50 backdrop-blur-md">
          <Globe className="w-3.5 h-3.5 text-slate-500 ml-2" />
          {languages.map((l) => (
            <Link
              key={l.code}
              href={`/location/${id}${active !== 'detail' ? `/${active}` : ''}?lang=${l.code}`}
              className={cn(
                "px-4 py-2 rounded-xl text-[10px] font-black tracking-widest uppercase transition-all",
                lang === l.code
                  ? "bg-indigo-600 text-white shadow-lg shadow-indigo-600/40"
                  : "text-slate-500 hover:text-slate-300 hover:bg-slate-700/50"
              )}
            >
              {l.label}
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
