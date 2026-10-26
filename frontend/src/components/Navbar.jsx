function Navbar({ onOpenHistory }) {
  return (
    <nav className="navbar">
      <div className="brand">
        <div className="logo">🔗</div>
        <div>
          <span className="brand-name">ShortURL</span>
          <span className="brand-by">by TeamCode</span>
        </div>
      </div>

      <div className="nav-actions">
        <button type="button" className="history-btn" onClick={onOpenHistory}>
          📋 Recent Links
        </button>

        <div className="powered-by">
          POWERED BY <strong>TeamCode</strong>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
