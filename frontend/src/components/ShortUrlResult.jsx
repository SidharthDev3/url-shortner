
function ShortUrlResult() {
  return (
    <div className="result-card">
      <div className="result-header">
        <span>SHORT LINK</span>
        <span className="ready">Ready</span>
      </div>

      <div className="result-content">
        <span className="short-url">shorturl.at/abc123</span>

        <div className="result-actions">
          <button>Copy Link</button>
          <button>QR</button>
        </div>
      </div>
    </div>
  );
}

export default ShortUrlResult;
