import { createContext, useContext, useState, useEffect } from 'react';
import api from '../lib/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isGuest, setIsGuest] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const guestFlag = localStorage.getItem('isGuest');

        if (guestFlag === 'true') {
            setUser({ id: 'guest', isGuest: true });
            setIsGuest(true);
            setLoading(false);
        } else if (token) {
            checkAuth();
        } else {
            setLoading(false);
        }
    }, []);

    const checkAuth = async () => {
        try {
            const { data } = await api.get('/api/profile');
            setUser({ id: data.user_id });
        } catch (error) {
            console.error('Auth check failed', error);
            logout();
        } finally {
            setLoading(false);
        }
    };

    const login = (userData, token, isGuest = false) => {
        localStorage.setItem('token', token);
        if (isGuest) {
            localStorage.setItem('isGuest', 'true');
            setIsGuest(true);
        } else {
            localStorage.removeItem('isGuest');
            setIsGuest(false);
        }
        setUser(userData);
    };

    const register = async (email, password) => {
        await api.post('/api/auth/register', { email, password });
    };

    const verify = async (email, code) => {
        await api.post('/api/auth/verify', { email, code });
    }

    const loginAsGuest = () => {
        localStorage.setItem('token', 'guest-session');
        localStorage.setItem('isGuest', 'true');
        setUser({ id: 'guest', isGuest: true });
        setIsGuest(true);
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('isGuest');
        setUser(null);
        setIsGuest(false);
    };

    return (
        <AuthContext.Provider value={{ user, isGuest, login, logout, register, verify, loginAsGuest, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
