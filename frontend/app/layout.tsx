import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Location Alpha - Global Geological Node Alpha",
  description: "High-performance geological infrastructure supporting recursive hierarchies and realtime engagement telemetry.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} antialiased bg-slate-950 text-slate-200 selection:bg-blue-500/30`}>
        <div className="flex flex-col min-h-screen">
          {/* Main Content Area */}
          <main className="flex-grow">
            {children}
          </main>
          
          {/* Footer Component */}
          <footer className="py-12 border-t border-slate-900 bg-slate-950/50">
            <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6 opacity-40">
               <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center font-black text-xs text-slate-500">A</div>
                  <span className="text-[10px] font-black uppercase tracking-[0.3em]">Location Alpha System</span>
               </div>
               <p className="text-[10px] font-black uppercase tracking-widest leading-none">Designed for Geological Node Excellence • 2026</p>
               <div className="flex gap-4">
                  <span className="text-[9px] font-black uppercase tracking-widest text-slate-700">Internal Reference: DC-0402</span>
               </div>
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}
