import Navbar from "../components/Navbar";
import UrlShortener from "../components/UrlShortener";
import ShortUrlResult from "../components/ShortUrlResult";

function Home() {
  return (
    <div className="home">
      <Navbar />

      <main>
        <UrlShortener />
        <ShortUrlResult />
      </main>
    </div>
  );
}

export default Home;
