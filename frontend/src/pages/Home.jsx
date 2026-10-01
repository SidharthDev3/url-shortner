import Navbar from '../components/Navbar';
import UrlShortener from '../components/UrlShortener';
import Footer from '../components/Footer';

function Home() {
  return (
    <div className="h-screen w-screen bg-[#fafbfc] text-slate-900 font-sans flex flex-col justify-between overflow-hidden select-none relative">
      {/* Ambient Decorative Background Elements */}
      <div className="absolute inset-0 pointer-events-none grid-pattern -z-10" />
      <div className="absolute -top-32 left-1/2 -translate-x-1/2 w-[780px] h-[400px] bg-gradient-to-tr from-sky-200/40 via-blue-100/35 to-indigo-100/30 blur-3xl pointer-events-none -z-10 rounded-full" />

      {/* Header */}
      <Navbar />

      {/* Main Interactive Hero & Console */}
      <UrlShortener />

      {/* Pinned Clean Footer */}
      <Footer />
    </div>
  );
}

export default Home;
