import Navbar from "../components/Navbar";
import UrlShortener from "../components/UrlShortener";

function Home() {
  return (
    <div className="home">
      <Navbar />

      <main>
        <UrlShortener />
      </main>
    </div>
  );
}

export default Home;
