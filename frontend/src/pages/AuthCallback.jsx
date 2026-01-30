import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Loader2 } from 'lucide-react';

export default function AuthCallback() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const { checkAuth } = useAuth(); // We might need to expose checkAuth or just rely on window.location.reload() for simplicity

    useEffect(() => {
        const token = searchParams.get('token');
        if (token) {
            localStorage.setItem('token', token);
            localStorage.removeItem('isGuest');

            // Force a reload to ensure AuthContext picks up the new token
            // Or if using checkAuth, call it. But reload is safer to reset all state.
            window.location.href = '/dashboard';
        } else {
            // No token, go back to login
            navigate('/login');
        }
    }, [searchParams, navigate]);

    return (
        <div className="flex h-screen items-center justify-center bg-background">
            <div className="text-center space-y-4">
                <Loader2 className="h-8 w-8 animate-spin text-primary mx-auto" />
                <p className="text-muted-foreground">Completing login...</p>
            </div>
        </div>
    );
}
