import { useState } from "react";
import { QRCodeSVG } from "qrcode.react";
import { shortenUrl } from "../services/api";

function UrlShortener() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [showQR, setShowQR] = useState(false);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState("");

  const handleShorten = async (e) => {
    if (e) e.preventDefault();
    setError("");

    if (!url.trim()) {
      setError("Please enter a URL to shorten");
      return;
    }

    setLoading(true);
    try {
      const data = await shortenUrl(url);
      setShortUrl(data.short_url);
      setShowQR(false);
    } catch (err) {
      setError(err.message || "Failed to shorten URL. Make sure the backend server is running.");
      setShortUrl("");
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = () => {
    if (!shortUrl) return;
    navigator.clipboard.writeText(shortUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="url-shortener">
      <div className="badge">
        🔵 Free &amp; Instant Link Shortener
      </div>

      <h1>
        Shorten your link <span>instantly</span>
      </h1>

      <p>
        Paste any long URL to generate a fast, secure, and clean
        redirect link in seconds.
      </p>

      <form className="url-form" onSubmit={handleShorten}>
        <input
          type="text"
          value={url}
          onChange={(e) => {
            setUrl(e.target.value);
            if (error) setError("");
          }}
          placeholder="https://example.com/your-long-url"
          disabled={loading}
        />

        <button type="submit" disabled={loading}>
          {loading ? "Shortening..." : "Shorten URL ⚡"}
        </button>
      </form>

      {/* Validation / Server Error Message */}
      {error && (
        <div className="error-message">
          ⚠️ {error}
        </div>
      )}

      {/* Short Link Output Card */}
      {shortUrl && (
        <div className="result-card">
          <div className="result-header">
            <span>SHORT LINK</span>
            <span className="ready">Ready</span>
          </div>

          <div className="result-content">
            <a
              href={shortUrl}
              target="_blank"
              rel="noreferrer"
              className="short-url"
              title="Click to visit shortened URL"
            >
              {shortUrl}
            </a>

            <div className="result-actions">
              <button type="button" onClick={handleCopy}>
                {copied ? "✓ Copied!" : "Copy Link"}
              </button>

              <button
                type="button"
                onClick={() => setShowQR(!showQR)}
                className={showQR ? "active-qr" : ""}
              >
                {showQR ? "Hide QR" : "QR"}
              </button>
            </div>
          </div>

          {showQR && (
            <div className="qr-code">
              <QRCodeSVG value={shortUrl} size={130} />
            </div>
          )}
        </div>
      )}
    </section>
  );
}

export default UrlShortener;
