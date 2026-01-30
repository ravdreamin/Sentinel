import { useState } from 'react';
import { useNavigate, Link, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { Loader2, Chrome } from 'lucide-react';
import api from '../lib/api';

export default function Login() {
    const navigate = useNavigate();
    const { login } = useAuth();
    const [searchParams] = useSearchParams();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const isGuestMode = searchParams.get('guest') === 'true';

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const { data } = await api.post('/login', { email, password });
            login(data.user, data.token, false);
            navigate('/dashboard');
        } catch (err) {
            setError(err.response?.data?.error || 'Login failed');
        } finally {
            setLoading(false);
        }
    };

    const handleGuestLogin = () => {
        login({ id: 0 }, 'guest_token', true);
        navigate('/dashboard');
    };

    const handleGoogleLogin = () => {
        window.location.href = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/auth/google`;
    };

    if (isGuestMode) {
        handleGuestLogin();
        return null;
    }

    return (
        <div className="min-h-screen bg-neutral-950 text-white font-sans flex flex-col items-center justify-center px-8">
            <div className="w-full max-w-sm space-y-8">
                <div className="text-center">
                    <Link to="/" className="text-2xl font-semibold tracking-tight hover:opacity-80 transition-opacity">
                        Sentinel
                    </Link>
                    <p className="text-neutral-400 mt-2">Welcome back</p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <input
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full h-12 px-4 bg-neutral-900 border border-neutral-800 rounded-xl text-white placeholder:text-neutral-500 focus:outline-none focus:border-neutral-600 transition-colors"
                            required
                        />
                    </div>
                    <div>
                        <input
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="w-full h-12 px-4 bg-neutral-900 border border-neutral-800 rounded-xl text-white placeholder:text-neutral-500 focus:outline-none focus:border-neutral-600 transition-colors"
                            required
                        />
                    </div>

                    {error && <p className="text-red-400 text-sm">{error}</p>}

                    <Button
                        type="submit"
                        disabled={loading}
                        className="w-full h-12 bg-white text-black hover:bg-neutral-200 rounded-xl font-medium"
                    >
                        {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Log in'}
                    </Button>
                </form>

                <div className="relative">
                    <div className="absolute inset-0 flex items-center">
                        <div className="w-full border-t border-neutral-800"></div>
                    </div>
                    <div className="relative flex justify-center text-sm">
                        <span className="bg-neutral-950 px-4 text-neutral-500">or</span>
                    </div>
                </div>

                <div className="space-y-3">
                    <Button
                        onClick={handleGoogleLogin}
                        variant="outline"
                        className="w-full h-12 border border-neutral-700 hover:border-neutral-500 rounded-xl text-neutral-300 hover:text-white"
                    >
                        <Chrome className="w-5 h-5 mr-2" />
                        Continue with Google
                    </Button>

                    <Button
                        onClick={handleGuestLogin}
                        variant="outline"
                        className="w-full h-12 border border-neutral-700 hover:border-neutral-500 rounded-xl text-neutral-300 hover:text-white"
                    >
                        Continue as guest
                    </Button>
                </div>

                <p className="text-center text-sm text-neutral-500">
                    Don't have an account?{' '}
                    <Link to="/register" className="text-white hover:underline">Sign up</Link>
                </p>
            </div>
        </div>
    );
}
