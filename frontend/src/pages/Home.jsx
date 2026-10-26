import { useState } from "react";
import Navbar from "../components/Navbar";
import UrlShortener from "../components/UrlShortener";
import HistoryModal from "../components/HistoryModal";

function Home() {
  const [showHistory, setShowHistory] = useState(false);

  return (
    <div className="home">
      <Navbar onOpenHistory={() => setShowHistory(true)} />

      <main>
        <UrlShortener />
      </main>

      <HistoryModal
        isOpen={showHistory}
        onClose={() => setShowHistory(false)}
      />
    </div>
  );
}

export default Home;
