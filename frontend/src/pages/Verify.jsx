import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate, useLocation } from 'react-router-dom';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '../components/ui/Card';
import { Loader2, Mail } from 'lucide-react';

export default function Verify() {
    const location = useLocation();
    const [email, setEmail] = useState(location.state?.email || '');
    const [code, setCode] = useState('');
    const [loading, setLoading] = useState(false);
    const { verify } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            await verify(email, code);
            alert('Verification successful! Please login.');
            navigate('/login');
        } catch (error) {
            alert('Verification failed. Invalid code or email.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-background relative overflow-hidden">
            <div className="absolute top-[20%] left-[50%] translate-x-[-50%] w-[600px] h-[600px] bg-purple-500/10 rounded-full blur-[120px] opacity-40 animate-pulse"></div>

            <Card className="w-[400px] backdrop-blur-md bg-card/60 border-accent/20 shadow-2xl z-10">
                <CardHeader className="text-center">
                    <div className="flex justify-center mb-4">
                        <div className="p-3 rounded-full bg-purple-500/10 border border-purple-500/20">
                            <Mail className="w-8 h-8 text-purple-500" />
                        </div>
                    </div>
                    <CardTitle className="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent">Verify Email</CardTitle>
                    <CardDescription>Enter the 6-digit code sent to your email</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Input
                                type="email"
                                placeholder="Email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                className="bg-background/50 border-accent/20 transition-all"
                            />
                        </div>
                        <div className="space-y-2">
                            <Input
                                type="text"
                                placeholder="000000"
                                value={code}
                                onChange={(e) => setCode(e.target.value)}
                                required
                                maxLength={6}
                                className="text-center text-2xl tracking-widest bg-background/50 border-accent/20 focus:border-purple-500/50 transition-all font-mono"
                            />
                        </div>
                        <Button type="submit" className="w-full bg-purple-600 hover:bg-purple-700 shadow-lg shadow-purple-500/20" disabled={loading}>
                            {loading ? <Loader2 className="animate-spin mr-2" /> : null}
                            Verify
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}
