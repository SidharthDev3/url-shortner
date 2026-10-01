function Footer() {
  return (
    <footer className="w-full flex-shrink-0 py-3 sm:py-3.5 border-t border-slate-200/70 bg-white/70 backdrop-blur-md text-xs text-slate-400 z-20">
      <div className="max-w-6xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-2">
        <div className="flex items-center gap-1.5 font-medium text-slate-400">
          <span>© 2025 ShortURL by TeamCode. All rights reserved.</span>
        </div>
        <div className="flex items-center gap-4 text-[11px] font-medium text-slate-500">
          <a className="hover:text-slate-800 transition-colors" href="#">Privacy</a>
          <span className="text-slate-300">•</span>
          <a className="hover:text-slate-800 transition-colors" href="#">Terms</a>
          <span className="text-slate-300">•</span>
          <a className="hover:text-slate-800 transition-colors" href="#">API Docs</a>
          <span className="text-slate-300">•</span>
          <a
            className="hover:text-sky-600 transition-colors flex items-center gap-0.5"
            href="https://teamcode.com"
            target="_blank"
            rel="noreferrer"
          >
            TeamCode.com
            <span className="material-symbols-outlined text-[11px]">arrow_outward</span>
          </a>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
