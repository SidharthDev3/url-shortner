import { useEffect, useState } from "react";
import { getRecentUrls } from "../services/api";

function HistoryModal({ isOpen, onClose }) {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(false);
  const [copiedCode, setCopiedCode] = useState("");

  useEffect(() => {
    if (isOpen) {
      setLoading(true);
      getRecentUrls()
        .then((data) => setHistory(data))
        .catch(() => setHistory([]))
        .finally(() => setLoading(false));
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleCopy = (shortUrl, code) => {
    navigator.clipboard.writeText(shortUrl);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(""), 2000);
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>📋 Shortened Links Dashboard</h3>
          <button className="modal-close" onClick={onClose}>✕</button>
        </div>

        <div className="modal-body">
          {loading ? (
            <p className="modal-empty">Loading history from database...</p>
          ) : history.length === 0 ? (
            <p className="modal-empty">No URLs shortened yet. Create one!</p>
          ) : (
            <div className="history-list">
              {history.map((item) => (
                <div key={item.short_code} className="history-item">
                  <div className="history-info">
                    <a
                      href={item.short_url}
                      target="_blank"
                      rel="noreferrer"
                      className="history-short"
                    >
                      {item.short_url}
                    </a>
                    <span className="history-original" title={item.original_url}>
                      {item.original_url}
                    </span>
                  </div>
                  <button
                    className="history-copy-btn"
                    onClick={() => handleCopy(item.short_url, item.short_code)}
                  >
                    {copiedCode === item.short_code ? "✓ Copied" : "Copy"}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default HistoryModal;
