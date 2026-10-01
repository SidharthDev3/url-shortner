import { useState } from 'react';
import { QRCodeSVG } from 'qrcode.react';

function UrlShortener() {
  const [url, setUrl] = useState('https://teamcode.com/projects/exclusive-promo-bundle-2025');
  const [shortUrl, setShortUrl] = useState('shorturl.at/tm-x89k');
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [showQR, setShowQR] = useState(false);
  const [highlightResult, setHighlightResult] = useState(false);

  const handleShorten = (e) => {
    e.preventDefault();
    if (!url.trim()) return;

    setLoading(true);

    // Simulate link generation (ready for API hookup)
    setTimeout(() => {
      const hash = Math.random().toString(36).substring(2, 7);
      setShortUrl(`shorturl.at/tm-${hash}`);
      setLoading(false);
      setHighlightResult(true);
      setTimeout(() => setHighlightResult(false), 1200);
    }, 400);
  };

  const handleCopy = () => {
    if (!shortUrl) return;
    const fullUrl = shortUrl.startsWith('http') ? shortUrl : `https://${shortUrl}`;
    
    if (navigator.clipboard) {
      navigator.clipboard.writeText(fullUrl);
    }
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <main className="flex-1 flex flex-col items-center justify-center px-4 w-full max-w-4xl mx-auto z-10 py-1">
      {/* Top Pill Badge & Title */}
      <div className="w-full text-center space-y-2 mb-5 sm:mb-6">
        <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-sky-50 border border-sky-200/80 text-sky-700 text-xs font-semibold shadow-xs backdrop-blur-xs">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-sky-500"></span>
          </span>
          <span>Free &amp; Instant Link Shortener</span>
        </div>

        <h1 className="text-3xl sm:text-5xl font-extrabold text-slate-900 tracking-tight leading-tight">
          Shorten your link{' '}
          <span className="bg-gradient-to-r from-sky-600 via-sky-500 to-indigo-600 bg-clip-text text-transparent">
            instantly
          </span>
        </h1>

        <p className="text-xs sm:text-sm text-slate-500 max-w-lg mx-auto font-normal">
          Paste any long URL to generate a fast, secure, and clean redirect link in seconds.
        </p>
      </div>

      {/* The Interactive Console / Elevated Card */}
      <div className="w-full max-w-2xl bg-white/95 backdrop-blur-md rounded-2xl border border-slate-200/90 shadow-subtle-card p-5 sm:p-6 transition-all">
        <form className="flex flex-col gap-3" onSubmit={handleShorten}>
          {/* Input Row */}
          <div className="flex flex-col sm:flex-row items-stretch gap-2.5">
            <div className="relative flex-1 group">
              <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-sky-600 transition-colors">
                <span className="material-symbols-outlined text-[19px]">link</span>
              </div>
              <input
                className="w-full pl-10 pr-4 py-2.5 sm:py-3 rounded-xl border border-slate-200 bg-slate-50/60 text-slate-800 placeholder-slate-400 text-sm focus:bg-white focus:border-sky-500 focus:ring-4 focus:ring-sky-100 outline-none transition-all"
                placeholder="https://example.com/very-long-url-to-shorten..."
                required
                type="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
              />
            </div>

            <button
              className="inline-flex items-center justify-center gap-1.5 px-6 py-2.5 sm:py-3 rounded-xl bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 active:scale-[0.98] text-white font-semibold text-sm shadow-md shadow-sky-500/25 transition-all cursor-pointer whitespace-nowrap disabled:opacity-70"
              type="submit"
              disabled={loading}
            >
              {loading ? (
                <>
                  <span className="material-symbols-outlined text-[18px] animate-spin">
                    progress_activity
                  </span>
                  <span>Shortening...</span>
                </>
              ) : (
                <>
                  <span>Shorten URL</span>
                  <span className="material-symbols-outlined text-[18px]">bolt</span>
                </>
              )}
            </button>
          </div>

          {/* Meta Sub-Badges */}
          <div className="flex items-center justify-between text-[11.5px] text-slate-400 px-1 pt-0.5">
            <span className="flex items-center gap-1.5 text-emerald-600 font-medium">
              <span className="material-symbols-outlined text-[15px] font-bold">verified</span>
              <span>SSL Secure &amp; Permanent Redirect</span>
            </span>
            <span className="text-slate-400 font-medium">No expiration • Unlimited clicks</span>
          </div>
        </form>

        {/* Short Link Output Card */}
        {shortUrl && (
          <div
            className={`mt-4 p-3.5 rounded-xl bg-slate-50/90 border border-slate-200/90 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 text-left transition-all ${
              highlightResult ? 'ring-2 ring-sky-500 bg-sky-50/60' : ''
            }`}
          >
            <div className="flex items-center gap-3 min-w-0">
              <div className="w-9 h-9 rounded-lg bg-white border border-slate-200/90 flex items-center justify-center flex-shrink-0 text-sky-600 shadow-xs">
                <span className="material-symbols-outlined text-[19px]">link</span>
              </div>
              <div className="min-w-0">
                <div className="flex items-center gap-1.5">
                  <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
                    Short Link
                  </span>
                  <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-bold bg-emerald-100/80 text-emerald-700 leading-none">
                    Ready
                  </span>
                </div>
                <div className="flex items-center gap-2 mt-0.5">
                  <a
                    className="font-mono font-semibold text-[14px] text-sky-600 hover:text-sky-700 hover:underline truncate"
                    href={shortUrl.startsWith('http') ? shortUrl : `https://${shortUrl}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {shortUrl}
                  </a>
                </div>
              </div>
            </div>

            {/* Output Action Buttons */}
            <div className="flex items-center gap-2 w-full sm:w-auto justify-end flex-shrink-0">
              <button
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white hover:bg-slate-100/90 border border-slate-200/90 text-slate-700 text-xs font-semibold transition-all shadow-xs active:scale-95 cursor-pointer"
                onClick={handleCopy}
                type="button"
              >
                <span className="material-symbols-outlined text-[15px]">
                  {copied ? 'check' : 'content_copy'}
                </span>
                <span>{copied ? 'Copied!' : 'Copy Link'}</span>
              </button>

              <button
                className={`inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg border text-xs font-semibold transition-all shadow-xs active:scale-95 cursor-pointer ${
                  showQR
                    ? 'bg-sky-50 border-sky-300 text-sky-700'
                    : 'bg-white hover:bg-slate-100/90 border-slate-200/90 text-slate-700'
                }`}
                onClick={() => setShowQR((prev) => !prev)}
                title="Show QR Code"
                type="button"
              >
                <span className="material-symbols-outlined text-[15px]">qr_code_2</span>
                <span>QR</span>
              </button>
            </div>
          </div>
        )}

        {/* QR Code Expansion Drawer */}
        {shortUrl && showQR && (
          <div className="mt-3.5 p-4 rounded-xl bg-slate-50 border border-slate-200/90 flex flex-col items-center justify-center gap-2 animate-in fade-in duration-200">
            <div className="p-3 bg-white rounded-xl shadow-xs border border-slate-100">
              <QRCodeSVG
                value={shortUrl.startsWith('http') ? shortUrl : `https://${shortUrl}`}
                size={130}
              />
            </div>
            <span className="text-[11px] font-medium text-slate-400">
              Scan to open short link
            </span>
          </div>
        )}

        {/* Card Bottom Bar: Custom Domain Link */}
        <div className="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-500">
          <div className="flex items-center gap-1.5">
            <span className="material-symbols-outlined text-[16px] text-sky-500">language</span>
            <span>Need a branded custom domain?</span>
          </div>
          <a
            className="font-semibold text-sky-600 hover:text-sky-700 hover:underline inline-flex items-center gap-1"
            href="https://teamcode.com"
            target="_blank"
            rel="noreferrer"
          >
            Explore domains on TeamCode
            <span className="material-symbols-outlined text-[13px]">arrow_forward</span>
          </a>
        </div>
      </div>
    </main>
  );
}

export default UrlShortener;
