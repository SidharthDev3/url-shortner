function Navbar() {
  return (
    <header className="w-full flex-shrink-0 z-20 px-6 sm:px-12 py-3.5 sm:py-4">
      <div className="max-w-6xl mx-auto flex items-center justify-between">
        {/* Brand Logo */}
        <a className="flex items-center gap-2.5 group transition-transform active:scale-95" href="#">
          <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-sky-600 to-sky-500 text-white flex items-center justify-center shadow-md shadow-sky-500/20 ring-1 ring-white/60">
            <span className="material-symbols-outlined text-[20px] text-white">link</span>
          </div>
          <div className="flex items-baseline gap-1.5">
            <span className="text-xl font-bold tracking-tight text-slate-900">
              Short<span className="text-sky-600">URL</span>
            </span>
            <span className="text-xs font-semibold text-slate-400 tracking-tight">by TeamCode</span>
          </div>
        </a>

        {/* Powered Badge */}
        <a
          className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/80 hover:bg-white border border-slate-200/80 shadow-xs text-xs text-slate-600 transition-all hover:border-slate-300 active:scale-95"
          href="https://teamcode.com"
          target="_blank"
          rel="noreferrer"
        >
          <span className="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500 ring-2 ring-emerald-100"></span>
          <span className="text-[11px] font-medium tracking-wide text-slate-500">
            POWERED BY <strong className="font-bold text-slate-800">TeamCode</strong>
          </span>
          <span className="material-symbols-outlined text-[13px] text-slate-400">open_in_new</span>
        </a>
      </div>
    </header>
  );
}

export default Navbar;
