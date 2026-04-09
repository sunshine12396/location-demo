'use client';

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip, ResponsiveContainer, Cell } from 'recharts';
import { TrendingUp, MessageSquare, Image, Video, MapPin, Search, Globe } from 'lucide-react';
import { LocationDetail } from '@/lib/api';

interface LocationStatsProps {
  location: LocationDetail;
  language: string;
}

export default function LocationStats({ location, language }: LocationStatsProps) {
  const stats = location.stats || { total_posts: 0, total_photos: 0, total_videos: 0, trending_score: 0 };

  const chartData = [
    {
      name: 'Posts',
      value: stats.total_posts,
      fill: '#3b82f6',
    },
    {
      name: 'Photos',
      value: stats.total_photos,
      fill: '#8b5cf6',
    },
    {
      name: 'Videos',
      value: stats.total_videos,
      fill: '#ec4899',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Visual Header */}
      <div className="bg-gradient-to-br from-blue-600/10 via-cyan-600/10 to-transparent p-5 rounded-2xl border border-blue-500/20 shadow-xl shadow-blue-900/10">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-blue-500 flex items-center justify-center text-white shadow-lg shadow-blue-500/30">
            <MapPin className="w-6 h-6" />
          </div>
          <div>
            <h3 className="text-2xl font-bold text-white tracking-tight leading-none mb-1">{location.name}</h3>
            <div className="flex items-center gap-2">
              <span className="px-2 py-0.5 rounded-md text-[10px] font-black uppercase tracking-widest bg-blue-500 text-white shadow-inner">
                {location.type}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-2 gap-3">
        {[
          { label: 'Total Posts', value: stats.total_posts, icon: MessageSquare, color: 'text-blue-400', bg: 'bg-blue-500/10', border: 'border-blue-500/20' },
          { label: 'Photos', value: stats.total_photos, icon: Image, color: 'text-purple-400', bg: 'bg-purple-500/10', border: 'border-purple-500/20' },
          { label: 'Videos', value: stats.total_videos, icon: Video, color: 'text-pink-400', bg: 'bg-pink-500/10', border: 'border-pink-500/20' },
          { label: 'Trending Score', value: stats.trending_score.toFixed(1), icon: TrendingUp, color: 'text-amber-400', bg: 'bg-amber-500/10', border: 'border-amber-500/20' },
        ].map((metric) => (
          <div key={metric.label} className={`flex flex-col p-4 rounded-xl border ${metric.bg} ${metric.border} shadow-inner`}>
            <div className="flex items-center gap-2 mb-2">
              <metric.icon className={`w-4 h-4 ${metric.color}`} />
              <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">{metric.label}</span>
            </div>
            <span className="text-2xl font-black text-white leading-tight">{metric.value}</span>
          </div>
        ))}
      </div>

      {/* Analytics Chart */}
      <div className="bg-slate-800/20 rounded-2xl p-5 border border-slate-700/30 overflow-hidden backdrop-blur-sm">
        <p className="text-[10px] font-black uppercase tracking-widest text-slate-500 mb-6 text-center">Visual Engagement Analysis</p>
        <div className="h-44 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} margin={{ top: 0, right: 0, left: -25, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#334155" vertical={false} />
              <XAxis dataKey="name" stroke="#64748b" style={{ fontSize: '10px' }} tickLine={false} axisLine={false} />
              <YAxis stroke="#64748b" style={{ fontSize: '10px' }} tickLine={false} axisLine={false} />
              <RechartsTooltip
                cursor={{ fill: 'rgba(255,255,255,0.05)' }}
                contentStyle={{
                  backgroundColor: '#0f172a',
                  border: '1px solid #334155',
                  borderRadius: '12px',
                  boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.5)',
                  fontSize: '12px'
                }}
                itemStyle={{ fontWeight: 'bold' }}
              />
              <Bar dataKey="value" radius={[6, 6, 0, 0]} barSize={40}>
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.fill} className="transition-all hover:opacity-80" />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}
