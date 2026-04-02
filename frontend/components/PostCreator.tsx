'use client';

import { useState } from 'react';
import { Send, AlertCircle, Loader2 } from 'lucide-react';
import { createPost, LocationDetail } from '@/lib/api';
import { useRouter } from 'next/navigation';

interface PostCreatorProps {
  locationId: number;
  locationName: string;
  language: string;
}

export default function PostCreator({
  locationId,
  locationName,
  language,
}: PostCreatorProps) {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!content.trim()) {
      return;
    }

    setIsSubmitting(true);
    try {
      await createPost(content, 'text', locationId);
      setContent('');
      router.refresh();
      // Optional: scroll to first post?
    } catch (err) {
      console.error('Failed to create post:', err);
      alert('Failed to publish post. Try again?');
    } finally {
      setIsSubmitting(false);
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
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Location Indicator */}
      <div className="flex items-center gap-2 bg-emerald-500/10 rounded-xl p-3 border border-emerald-500/20 shadow-inner">
        <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse ring-4 ring-emerald-500/20" />
        <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">Live Feed Activation: <strong className="text-emerald-400 font-black">{locationName}</strong></span>
      </div>

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
            <span className={`text-[11px] font-black tracking-widest uppercase transition-colors ${
              remainingChars < 0
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
              - {remainingChars < 0 ? `LIMIT OVERFLOWED BY ${Math.abs(remainingChars)}` : `${remainingChars} REMAINING` }
            </span>
          )}
        </div>

        <button
          type="submit"
          disabled={!content.trim() || isSubmitting || remainingChars < 0}
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
  );
}
