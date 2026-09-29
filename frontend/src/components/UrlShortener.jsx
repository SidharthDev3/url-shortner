function UrlShortener() {
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
          placeholder="https://example.com/your-long-url"
        />

        <button>
          Shorten URL ⚡
        </button>
      </div>
    </section>
  );
}

export default UrlShortener;
