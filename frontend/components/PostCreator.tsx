'use client';

import { useState } from 'react';
import { Send, AlertCircle, Loader2, MapPin, X } from 'lucide-react';
import { createPost, SearchResult } from '@/lib/api';
import { useRouter } from 'next/navigation';
import LocationSearch from './LocationSearch';

interface PostCreatorProps {
  locationId?: number;
  locationName?: string;
  language: string;
  onSuccess?: () => void;
}

export default function PostCreator({
  locationId: initialLocationId,
  locationName: initialLocationName,
  language,
  onSuccess,
}: PostCreatorProps) {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Local language state specifically for searching and hydrating this post intent
  const [postLang, setPostLang] = useState(language);

  // New state for dynamic location selection
  const [selectedLoc, setSelectedLoc] = useState<{ id: number, externalId?: string, name: string } | null>(
    initialLocationId && initialLocationName ? { id: initialLocationId, name: initialLocationName } : null
  );
  const [showPicker, setShowPicker] = useState(false);

  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!content.trim() || !selectedLoc) {
      return;
    }

    setIsSubmitting(true);
    setError(null);
    try {
      await createPost(content, 'text', selectedLoc.id > 0 ? selectedLoc.id : undefined, selectedLoc.externalId, postLang);
      setContent('');
      // If we were using a picked location, clear it
      if (!initialLocationId) {
        setSelectedLoc(null);
      }
      if (onSuccess) onSuccess();
      router.refresh();
    } catch (err: any) {
      console.error('Failed to create post:', err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePickerSelect = (id: number, externalId: string | undefined, name: string) => {
    // We allow either an existing ID OR a fresh External ID (Save on Intent)
    if (id > 0 || externalId) {
      setSelectedLoc({ id, externalId, name });
      setShowPicker(false);
    } else {
      setError("Please select an existing location node to publish.");
    }
  };

  const charCount = content.length;
  const maxChars = 280;
  const remainingChars = maxChars - charCount;

  const placeholders: { [key: string]: string } = {
    en: 'Share your thoughts about this location...',
    vi: 'Chia sẻ suy nghĩ của bạn về địa điểm này...',
    ja: 'この場所についてあなたの考えを共有してください...',
  };

  return (
    <div className="space-y-4">
      {/* Error Message */}
      {error && (
        <div className="flex items-start gap-4 p-5 bg-red-500/10 border border-red-500/30 rounded-2xl animate-in fade-in slide-in-from-top-1">
          <div className="w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center text-red-400 shrink-0">
            <AlertCircle className="w-5 h-5" />
          </div>
          <div className="flex-1 space-y-1">
            <h4 className="text-sm font-black uppercase tracking-widest text-red-400">PUBLISH_FAILED_CORE_ERROR</h4>
            <p className="text-xs text-red-500/80 leading-relaxed font-medium">
              {error}
            </p>
          </div>
          <button
            type="button"
            onClick={() => setError(null)}
            className="text-red-400 hover:text-red-300 p-1"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Location Selector/Indicator */}
      {selectedLoc ? (
        <div className="flex items-center justify-between bg-emerald-500/10 rounded-xl p-3 border border-emerald-500/20 shadow-inner group/loc">
          <div className="flex items-center gap-2">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)] animate-pulse ring-4 ring-emerald-500/20" />
            <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">
              Live Feed Activation: <strong className="text-emerald-400 font-black">{selectedLoc.name}</strong>
            </span>
          </div>
          {!initialLocationId && (
            <button
              type="button"
              onClick={() => setSelectedLoc(null)}
              className="text-slate-500 hover:text-red-400 transition-colors"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      ) : showPicker ? (
        <div className="relative z-[60] bg-slate-900/50 backdrop-blur-md rounded-2xl p-4 border border-slate-700/50 space-y-3 shadow-2xl">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-3">
              <span className="text-[10px] font-black uppercase tracking-widest text-blue-400">Select Location Node</span>
              <div className="flex bg-slate-800 rounded p-0.5 border border-slate-700/50">
                {['en', 'vi', 'ja'].map((l) => (
                  <button 
                    key={l}
                    onClick={() => setPostLang(l)} 
                    className={`px-2 py-0.5 text-[9px] font-black uppercase rounded ${postLang === l ? 'bg-blue-500 text-white' : 'text-slate-500 hover:text-slate-300'}`}
                  >
                    {l}
                  </button>
                ))}
              </div>
            </div>
            <button onClick={() => setShowPicker(false)} className="text-slate-500 hover:text-white transition-colors">
              <X className="w-4 h-4" />
            </button>
          </div>
          <LocationSearch
            language={postLang}
            customOnSelect={(id, ext, name) => handlePickerSelect(id, ext, name || "Unknown")}
            useAutocomplete={true} // Allow Google results for detailed discovery/tagging
          />
        </div>
      ) : (
        <button
          type="button"
          onClick={() => setShowPicker(true)}
          className="w-full flex items-center justify-center gap-2 bg-slate-800/30 hover:bg-slate-800/50 border border-dashed border-slate-700 rounded-xl p-3 text-[10px] font-black uppercase tracking-widest text-slate-500 transition-all hover:text-blue-400 hover:border-blue-500/50"
        >
          <MapPin className="w-3.5 h-3.5" />
          Tag Location To Publish
        </button>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Content Area */}
        <div className="relative">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value.slice(0, maxChars))}
            placeholder={placeholders[language] || placeholders.en}
            disabled={isSubmitting}
            className="w-full px-5 py-4 bg-slate-800/30 border border-slate-700/50 rounded-2xl text-white placeholder-slate-600 focus:outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all resize-none disabled:opacity-50 min-h-[140px] shadow-2xl shadow-indigo-900/10 leading-relaxed font-medium text-lg"
            rows={4}
          />

          {/* Glow behind textarea */}
          <div className="absolute -inset-1 bg-gradient-to-r from-blue-500/10 via-cyan-500/5 to-purple-500/10 blur-xl opacity-20 pointer-events-none -z-10" />
        </div>

        {/* Character Count & Submit */}
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col">
            <div className="flex items-center gap-2">
              <span className={`text-[11px] font-black tracking-widest uppercase transition-colors ${remainingChars < 0
                  ? 'text-red-500 animate-pulse'
                  : remainingChars < 50
                    ? 'text-amber-500'
                    : 'text-slate-500 font-bold'
                }`}>
                {charCount} / {maxChars}
              </span>
            </div>
            {remainingChars < 50 && (
              <span className={`text-[8px] font-black uppercase transition-colors ${remainingChars < 0 ? 'text-red-600' : 'text-amber-600'}`}>
                - {remainingChars < 0 ? `LIMIT OVERFLOWED BY ${Math.abs(remainingChars)}` : `${remainingChars} REMAINING`}
              </span>
            )}
          </div>

          <button
            type="submit"
            disabled={!content.trim() || isSubmitting || remainingChars < 0 || !selectedLoc}
            className="relative group flex-1 max-w-[200px]"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-cyan-600 rounded-xl blur-lg transition hover:blur-xl opacity-40 group-hover:opacity-70 group-disabled:opacity-0" />
            <div className="relative w-full h-12 bg-gradient-to-r from-blue-600 to-cyan-600 disabled:from-slate-800 disabled:to-slate-800 text-white font-black uppercase tracking-[0.2em] text-xs rounded-xl transition-all flex items-center justify-center gap-3 disabled:opacity-50 disabled:cursor-not-allowed group-hover:translate-y-[-2px] group-active:translate-y-[1px]">
              {isSubmitting ? (
                <Loader2 className="w-4 h-4 animate-spin text-white shadow-lg" />
              ) : (
                <Send className="w-4 h-4 shadow-lg group-hover:rotate-12 transition-transform" />
              )}
              {isSubmitting ? 'UPDATING...' : 'PUBLISH HUD'}
            </div>
          </button>
        </div>
      </form>
    </div>
  );
}
