import { useState } from "react";

function UrlShortener() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");

  const handleShorten = () => {
    if (!url.trim()) {
      alert("Please enter a URL");
      return;
    }

    // Temporary short URL for testing the frontend
    setShortUrl("shorturl.at/abc123");
  };

  return (
    <section className="url-shortener">
      <div className="badge">
        🔵 Free & Instant Link Shortener
      </div>

      <h1>
        Shorten your link <span>instantly</span>
      </h1>

      <p>
        Paste any long URL to generate a fast, secure, and clean
        redirect link in seconds.
      </p>

      <div className="url-form">
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com/your-long-url"
        />

        <button onClick={handleShorten}>
          Shorten URL ⚡
        </button>
      </div>

      {shortUrl && (
        <div className="result-card">
          <div className="result-header">
            <span>SHORT LINK</span>
            <span className="ready">Ready</span>
          </div>

          <div className="result-content">
            <span className="short-url">{shortUrl}</span>

            <div className="result-actions">
              <button
                onClick={() =>
                  navigator.clipboard.writeText(shortUrl)
                }
              >
                Copy Link
              </button>

              <button>QR</button>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}

export default UrlShortener;

