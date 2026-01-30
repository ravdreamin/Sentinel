import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { ArrowRight, FileText, Zap, Download, Globe, Shield, Clock, Code2, Github } from 'lucide-react';

export default function Landing() {
    const { user } = useAuth();

    return (
        <div className="min-h-screen bg-neutral-950 text-white font-sans">
            {/* Header */}
            <header className="px-8 py-6 flex justify-between items-center max-w-6xl mx-auto">
                <Link to="/" className="text-xl font-semibold tracking-tight hover:opacity-80 transition-opacity">
                    Sentinel
                </Link>
                <div className="flex gap-4 items-center">
                    {user ? (
                        <Link to="/dashboard">
                            <Button size="sm" className="bg-white text-black hover:bg-neutral-200 rounded-full px-5 h-9 text-sm font-medium">
                                Go to Dashboard
                            </Button>
                        </Link>
                    ) : (
                        <>
                            <Link to="/login" className="text-sm text-neutral-400 hover:text-white transition-colors">
                                Log in
                            </Link>
                            <Link to="/register">
                                <Button size="sm" className="bg-white text-black hover:bg-neutral-200 rounded-full px-5 h-9 text-sm font-medium">
                                    Sign up
                                </Button>
                            </Link>
                        </>
                    )}
                </div>
            </header>

            {/* Hero */}
            <section className="px-8 pt-24 pb-32 max-w-6xl mx-auto">
                <div className="max-w-3xl space-y-8">
                    <div className="inline-flex items-center gap-2 bg-neutral-900 border border-neutral-800 rounded-full px-4 py-2 text-sm">
                        <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                        <span className="text-neutral-400">Open source web scraper</span>
                    </div>

                    <h1 className="text-5xl md:text-7xl font-bold tracking-tight leading-[1.1]">
                        Extract web data<br />
                        <span className="text-neutral-500">in seconds.</span>
                    </h1>

                    <p className="text-xl text-neutral-400 max-w-xl leading-relaxed">
                        Sentinel is a high-performance web extraction tool. Upload a list of URLs,
                        and get structured data back — titles, metadata, links, all in clean JSON.
                    </p>

                    <div className="flex flex-col sm:flex-row gap-4 pt-4">
                        {user ? (
                            <Link to="/dashboard">
                                <Button size="lg" className="bg-white text-black hover:bg-neutral-200 rounded-full h-14 px-8 text-base font-medium group">
                                    Go to Dashboard <ArrowRight className="ml-2 w-4 h-4 group-hover:translate-x-1 transition-transform" />
                                </Button>
                            </Link>
                        ) : (
                            <Link to="/login?guest=true">
                                <Button size="lg" className="bg-white text-black hover:bg-neutral-200 rounded-full h-14 px-8 text-base font-medium group">
                                    Try it free <ArrowRight className="ml-2 w-4 h-4 group-hover:translate-x-1 transition-transform" />
                                </Button>
                            </Link>
                        )}
                        <a href="https://github.com/ravdreamin/Sentinel" target="_blank" rel="noopener noreferrer">
                            <Button variant="outline" size="lg" className="border-neutral-700 hover:border-neutral-500 rounded-full h-14 px-8 text-base">
                                <Code2 className="w-4 h-4 mr-2" /> View on GitHub
                            </Button>
                        </a>
                    </div>

                    {!user && (
                        <p className="text-sm text-neutral-600">
                            No account required • Guest access available
                        </p>
                    )}
                </div>
            </section>

            {/* How it works */}
            <section className="px-8 py-24 border-t border-neutral-800 bg-neutral-900/30">
                <div className="max-w-6xl mx-auto">
                    <div className="text-center mb-16">
                        <h2 className="text-3xl font-bold mb-4">How it works</h2>
                        <p className="text-neutral-400 max-w-lg mx-auto">
                            Three simple steps to extract data from any website
                        </p>
                    </div>

                    <div className="grid md:grid-cols-3 gap-8">
                        <div className="bg-neutral-900 border border-neutral-800 rounded-2xl p-8 space-y-4">
                            <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center">
                                <FileText className="w-6 h-6 text-blue-400" />
                            </div>
                            <h3 className="text-xl font-semibold">1. Upload URLs</h3>
                            <p className="text-neutral-400 leading-relaxed">
                                Upload a .txt file with your target URLs, one per line. We'll process them all in parallel.
                            </p>
                        </div>

                        <div className="bg-neutral-900 border border-neutral-800 rounded-2xl p-8 space-y-4">
                            <div className="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center">
                                <Zap className="w-6 h-6 text-purple-400" />
                            </div>
                            <h3 className="text-xl font-semibold">2. We extract</h3>
                            <p className="text-neutral-400 leading-relaxed">
                                Our Go-powered workers fetch each page and extract titles, meta tags, headers, and all links.
                            </p>
                        </div>

                        <div className="bg-neutral-900 border border-neutral-800 rounded-2xl p-8 space-y-4">
                            <div className="w-12 h-12 rounded-xl bg-green-500/10 flex items-center justify-center">
                                <Download className="w-6 h-6 text-green-400" />
                            </div>
                            <h3 className="text-xl font-semibold">3. Download JSON</h3>
                            <p className="text-neutral-400 leading-relaxed">
                                Get your data as a clean, structured JSON file. Ready for analysis or integration.
                            </p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Features */}
            <section className="px-8 py-24 border-t border-neutral-800">
                <div className="max-w-6xl mx-auto">
                    <div className="grid md:grid-cols-2 gap-16 items-center">
                        <div className="space-y-8">
                            <h2 className="text-4xl font-bold tracking-tight">
                                Built for speed & reliability
                            </h2>
                            <p className="text-neutral-400 text-lg leading-relaxed">
                                Sentinel uses a distributed worker pool architecture written in Go.
                                Fast, concurrent, and built to handle thousands of URLs without breaking a sweat.
                            </p>

                            <div className="grid grid-cols-2 gap-6">
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-white">
                                        <Globe className="w-5 h-5 text-blue-400" />
                                        <span className="font-medium">Parallel Processing</span>
                                    </div>
                                    <p className="text-sm text-neutral-500">Multiple workers fetching simultaneously</p>
                                </div>
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-white">
                                        <Shield className="w-5 h-5 text-green-400" />
                                        <span className="font-medium">SHA-256 Hashing</span>
                                    </div>
                                    <p className="text-sm text-neutral-500">Content integrity verification</p>
                                </div>
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-white">
                                        <Clock className="w-5 h-5 text-yellow-400" />
                                        <span className="font-medium">Real-time Metrics</span>
                                    </div>
                                    <p className="text-sm text-neutral-500">Track progress as it happens</p>
                                </div>
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-white">
                                        <Code2 className="w-5 h-5 text-purple-400" />
                                        <span className="font-medium">Clean JSON Output</span>
                                    </div>
                                    <p className="text-sm text-neutral-500">Structured, ready-to-use format</p>
                                </div>
                            </div>
                        </div>

                        {/* Code Preview */}
                        <div className="bg-neutral-900 rounded-2xl border border-neutral-800 overflow-hidden">
                            <div className="flex items-center gap-2 px-4 py-3 border-b border-neutral-800 bg-neutral-900/50">
                                <span className="w-3 h-3 rounded-full bg-red-500/60"></span>
                                <span className="w-3 h-3 rounded-full bg-yellow-500/60"></span>
                                <span className="w-3 h-3 rounded-full bg-green-500/60"></span>
                                <span className="ml-4 text-xs text-neutral-500 font-mono">results.json</span>
                            </div>
                            <pre className="p-6 text-sm text-neutral-300 font-mono overflow-x-auto">
                                {`[
  {
    "url": "https://example.com",
    "title": "Example Domain",
    "h1": "Example Domain",
    "meta_description": "...",
    "status_code": 200,
    "response_time": 145,
    "content_hash": "2cf24dba5...",
    "links": [
      "https://iana.org/domains"
    ]
  }
]`}
                            </pre>
                        </div>
                    </div>
                </div>
            </section>

            {/* Tech Stack */}
            <section className="px-8 py-24 border-t border-neutral-800 bg-neutral-900/30">
                <div className="max-w-6xl mx-auto text-center">
                    <h2 className="text-3xl font-bold mb-4">Tech Stack</h2>
                    <p className="text-neutral-400 mb-12 max-w-lg mx-auto">
                        Built with modern, battle-tested technologies
                    </p>

                    <div className="flex flex-wrap justify-center gap-4">
                        {['Go 1.25', 'Gin Framework', 'PostgreSQL', 'React', 'Vite', 'Tailwind CSS'].map((tech) => (
                            <div key={tech} className="px-6 py-3 bg-neutral-800 border border-neutral-700 rounded-full text-sm font-medium">
                                {tech}
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* CTA */}
            <section className="px-8 py-32 border-t border-neutral-800 text-center">
                <div className="max-w-2xl mx-auto space-y-8">
                    <h2 className="text-4xl md:text-5xl font-bold tracking-tight">
                        Ready to extract?
                    </h2>
                    <p className="text-xl text-neutral-400">
                        Start for free. No account required.
                    </p>
                    {user ? (
                        <Link to="/dashboard">
                            <Button size="lg" className="bg-white text-black hover:bg-neutral-200 rounded-full h-14 px-10 text-lg font-medium">
                                Go to Dashboard <ArrowRight className="ml-2 w-5 h-5" />
                            </Button>
                        </Link>
                    ) : (
                        <Link to="/login?guest=true">
                            <Button size="lg" className="bg-white text-black hover:bg-neutral-200 rounded-full h-14 px-10 text-lg font-medium">
                                Get started <ArrowRight className="ml-2 w-5 h-5" />
                            </Button>
                        </Link>
                    )}
                </div>
            </section>

            {/* Footer */}
            <footer className="px-8 py-12 border-t border-neutral-800">
                <div className="max-w-6xl mx-auto flex flex-col md:flex-row justify-between items-center gap-6">
                    <span className="text-xl font-semibold tracking-tight">Sentinel</span>
                    <div className="flex gap-8 text-sm text-neutral-500">
                        <a href="https://github.com/ravdreamin/Sentinel" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">GitHub</a>
                        <Link to="/dashboard" className="hover:text-white transition-colors">Dashboard</Link>
                    </div>
                    <p className="text-sm text-neutral-600">© 2026 Sentinel</p>
                </div>
            </footer>
        </div>
    );
}
