'use client';

import { Post } from "@/lib/api";
import { MessageSquare, MapPin, Clock, Video, Image, FileText } from 'lucide-react';
import Link from "next/link";
import { cn } from "@/lib/utils";

interface PostFeedProps {
  posts: Post[];
  language?: string;
  className?: string;
}

export default function PostFeed({ posts, language = "en", className }: PostFeedProps) {
  if (!posts || posts.length === 0) {
    return (
      <div className={cn("text-center py-20 bg-slate-800/10 rounded-3xl border-2 border-dashed border-slate-700/50 flex flex-col items-center justify-center opacity-60", className)}>
        <MessageSquare className="w-16 h-16 text-slate-700 mb-4 animate-pulse" />
        <p className="text-sm font-black uppercase tracking-[0.3em] text-slate-600">ZERO DATA IN PULSE</p>
        <p className="text-xs text-slate-500 font-medium mt-1">Be the first to share a moment at this node.</p>
      </div>
    );
  }

  const getMediaIcon = (type: string) => {
    switch (type) {
      case 'photo': return <Image className="w-4 h-4 text-purple-400" />;
      case 'video': return <Video className="w-4 h-4 text-pink-400" />;
      default: return <FileText className="w-4 h-4 text-blue-400" />;
    }
  };

  return (
    <div className={cn("space-y-4", className)}>
      {posts.map((post, idx) => (
        <div 
          key={post.id} 
          className="group relative transition-all duration-300 transform hover:translate-y-[-4px]"
        >
          {/* Subtle Accent Glow */}
          <div className="absolute inset-0 bg-gradient-to-r from-blue-500/10 via-purple-500/5 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity blur-xl -z-10" />
          
          <div className="bg-slate-800/30 backdrop-blur-xl rounded-2xl p-6 border border-slate-700/40 group-hover:border-slate-500/50 transition-all shadow-2xl shadow-indigo-950/20">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="relative">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-500 flex items-center justify-center font-black text-xl text-white shadow-lg group-hover:scale-105 transition-transform ring-4 ring-slate-800/50">
                    U
                    <div className="absolute -top-1 -right-1 w-4 h-4 rounded-full bg-emerald-500 border-2 border-slate-800 flex items-center justify-center">
                       <div className="w-1.5 h-1.5 rounded-full bg-white animate-ping" />
                    </div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-white font-black text-base transition-colors group-hover:text-blue-400">User {post.user_id}</span>
                  </div>
                  <div className="flex items-center gap-1.5 text-slate-500 text-[9px] font-black uppercase tracking-widest mt-0.5">
                    <Clock className="w-3 h-3" />
                    {new Date(post.created_at).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })}
                    <span className="opacity-20">•</span>
                    {new Date(post.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
                  </div>
                </div>
              </div>
              
              <div className="flex items-center gap-3">
                {post.location_name && (
                  <Link 
                    href={`/location/${post.location_id}?lang=${language}`} 
                    className="flex items-center gap-1.5 px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-wider bg-slate-700/30 text-slate-400 hover:bg-blue-500/20 hover:text-blue-300 transition-all border border-slate-600/50"
                  >
                    <MapPin className="w-3 h-3" />
                    {post.location_name}
                  </Link>
                )}
                <div className="p-2 rounded-lg bg-slate-800/40 border border-slate-700/50">
                   {getMediaIcon(post.media_type)}
                </div>
              </div>
            </div>
            
            <p className="text-slate-200 text-lg leading-relaxed font-medium selection:bg-blue-500/30">{post.content}</p>
            
            <div className="mt-6 pt-4 border-t border-slate-700/30 flex items-center justify-between opacity-0 group-hover:opacity-100 transition-opacity">
               <div className="flex gap-4">
                  <button className="text-[10px] font-black uppercase tracking-widest text-slate-500 hover:text-blue-400">Like</button>
                  <button className="text-[10px] font-black uppercase tracking-widest text-slate-500 hover:text-blue-400">Comment</button>
                  <button className="text-[10px] font-black uppercase tracking-widest text-slate-500 hover:text-blue-400">Share</button>
               </div>
               <span className="text-[9px] font-black uppercase tracking-widest text-slate-600">ID: {post.id}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
