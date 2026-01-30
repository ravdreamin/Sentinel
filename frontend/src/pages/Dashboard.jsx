import { useState, useEffect, useRef } from 'react';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { Loader2, Upload, Download, Check, AlertCircle, LogOut, Clock } from 'lucide-react';
import { Link } from 'react-router-dom';
import api from '../lib/api';

export default function Dashboard() {
    const { user, isGuest, logout } = useAuth();
    const [file, setFile] = useState(null);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState(null);
    const [currentJob, setCurrentJob] = useState(null);
    const [metrics, setMetrics] = useState(null);
    const [startTime, setStartTime] = useState(null);
    const [elapsedTime, setElapsedTime] = useState(0);
    const timerRef = useRef(null);

    // Timer for elapsed time
    useEffect(() => {
        if (startTime && currentJob && currentJob.status !== 'completed') {
            timerRef.current = setInterval(() => {
                setElapsedTime(Math.floor((Date.now() - startTime) / 1000));
            }, 1000);
        } else if (currentJob?.status === 'completed' && timerRef.current) {
            clearInterval(timerRef.current);
        }
        return () => clearInterval(timerRef.current);
    }, [startTime, currentJob?.status]);

    // Polling for job status
    useEffect(() => {
        let interval;
        if (currentJob && currentJob.status !== 'completed') {
            interval = setInterval(async () => {
                try {
                    const statusRes = await api.get(`/api/jobs/${currentJob.filename}/status`);
                    setCurrentJob(prev => ({ ...prev, ...statusRes.data }));

                    const metricsRes = await api.get(`/api/jobs/${currentJob.filename}/metrics`);
                    setMetrics(metricsRes.data);

                    if (statusRes.data.status === 'completed') clearInterval(interval);
                } catch (err) {
                    console.error("Polling error", err);
                }
            }, 1500);
        }
        return () => clearInterval(interval);
    }, [currentJob]);

    const handleUpload = async () => {
        if (!file) return;
        setUploading(true);
        setError(null);
        const formData = new FormData();
        formData.append('document', file);

        try {
            const { data } = await api.post('/api/upload', formData, {
                headers: { 'Content-Type': 'multipart/form-data' },
            });
            setStartTime(Date.now());
            setElapsedTime(0);
            setCurrentJob({
                filename: data.filename,
                total: data.total_found,
                completed: 0,
                failed: 0,
                status: 'processing'
            });
            setMetrics(null);
            setFile(null);
        } catch (err) {
            setError(err.response?.data?.error || 'Upload failed');
        } finally {
            setUploading(false);
        }
    };

    const downloadResults = async () => {
        if (!currentJob) return;
        const response = await api.get(`/api/jobs/${currentJob.filename}/download`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `${currentJob.filename}_results.json`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    };

    const formatTime = (seconds) => {
        if (seconds < 60) return `${seconds}s`;
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        if (mins < 60) return `${mins}m ${secs}s`;
        const hours = Math.floor(mins / 60);
        const remainMins = mins % 60;
        return `${hours}h ${remainMins}m`;
    };

    const resetJob = () => {
        setCurrentJob(null);
        setMetrics(null);
        setStartTime(null);
        setElapsedTime(0);
    };

    const progress = currentJob ? Math.round(((currentJob.completed + currentJob.failed) / (currentJob.total || 1)) * 100) : 0;
    const isComplete = currentJob?.status === 'completed';

    return (
        <div className="min-h-screen bg-neutral-950 text-white font-sans">
            {/* Header */}
            <header className="px-8 py-6 flex justify-between items-center border-b border-neutral-800">
                <div className="flex items-center gap-8">
                    <Link to="/" className="text-xl font-semibold tracking-tight hover:opacity-80 transition-opacity">
                        Sentinel
                    </Link>
                    <nav className="flex gap-6 text-sm text-neutral-400">
                        <Link to="/dashboard" className="text-white">Dashboard</Link>
                        {!isGuest && <Link to="/profile" className="hover:text-white transition-colors">History</Link>}
                    </nav>
                </div>
                <div className="flex items-center gap-4">
                    <span className="text-sm text-neutral-500">
                        {isGuest ? 'Guest' : user?.email || `User ${user?.id}`}
                    </span>
                    <button
                        onClick={logout}
                        className="flex items-center gap-2 text-sm text-neutral-400 hover:text-white transition-colors"
                    >
                        <LogOut className="w-4 h-4" />
                        Log out
                    </button>
                </div>
            </header>

            <main className="max-w-2xl mx-auto px-8 py-16">
                {/* Upload Section */}
                {!currentJob && (
                    <div className="space-y-8">
                        <div>
                            <h1 className="text-3xl font-bold tracking-tight mb-2">New extraction</h1>
                            <p className="text-neutral-400">Upload a text file containing URLs to extract data from.</p>
                        </div>

                        {/* File Drop Zone */}
                        <div className="relative">
                            <input
                                type="file"
                                onChange={(e) => setFile(e.target.files?.[0])}
                                className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
                                accept=".txt,.pdf"
                            />
                            <div className={`border-2 border-dashed rounded-2xl p-12 text-center transition-colors ${file ? 'border-green-500 bg-green-500/5' : 'border-neutral-700 hover:border-neutral-500'}`}>
                                {file ? (
                                    <div className="flex items-center justify-center gap-3 text-green-400">
                                        <Check className="w-5 h-5" />
                                        <span className="font-medium">{file.name}</span>
                                    </div>
                                ) : (
                                    <div className="space-y-3">
                                        <Upload className="w-8 h-8 mx-auto text-neutral-500" />
                                        <p className="text-neutral-400">Drop your file here, or click to browse</p>
                                        <p className="text-xs text-neutral-600">Supports .txt and .pdf files</p>
                                    </div>
                                )}
                            </div>
                        </div>

                        {error && (
                            <div className="flex items-center gap-2 text-red-400 text-sm">
                                <AlertCircle className="w-4 h-4" />
                                {error}
                            </div>
                        )}

                        <Button
                            onClick={handleUpload}
                            disabled={!file || uploading}
                            className="w-full h-14 bg-white text-black hover:bg-neutral-200 rounded-xl text-base font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {uploading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Start extraction'}
                        </Button>
                    </div>
                )}

                {/* Progress Section */}
                {currentJob && (
                    <div className="space-y-8">
                        <div>
                            <h1 className="text-3xl font-bold tracking-tight mb-2">
                                {isComplete ? 'Extraction complete' : 'Extracting...'}
                            </h1>
                            <p className="text-neutral-400">{currentJob.filename}</p>
                        </div>

                        {/* Progress Bar */}
                        <div className="space-y-3">
                            <div className="flex justify-between text-sm">
                                <span className="text-neutral-400">Progress</span>
                                <span className="font-medium">{progress}%</span>
                            </div>
                            <div className="h-2 bg-neutral-800 rounded-full overflow-hidden">
                                <div
                                    className={`h-full transition-all duration-500 ${isComplete ? 'bg-green-500' : 'bg-white'}`}
                                    style={{ width: `${progress}%` }}
                                />
                            </div>
                        </div>

                        {/* Stats Grid */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <div className="bg-neutral-900 rounded-xl p-5">
                                <p className="text-sm text-neutral-500 mb-1">Processed</p>
                                <p className="text-2xl font-semibold">{currentJob.completed + currentJob.failed}/{currentJob.total}</p>
                            </div>
                            <div className="bg-neutral-900 rounded-xl p-5">
                                <p className="text-sm text-neutral-500 mb-1">Total time</p>
                                <p className="text-2xl font-semibold flex items-center gap-2">
                                    {!isComplete && <Clock className="w-4 h-4 text-neutral-500 animate-pulse" />}
                                    {formatTime(elapsedTime)}
                                </p>
                            </div>
                            <div className="bg-neutral-900 rounded-xl p-5">
                                <p className="text-sm text-neutral-500 mb-1">Avg. time</p>
                                <p className="text-2xl font-semibold">{metrics?.avg_response_time ? `${Math.round(metrics.avg_response_time)}ms` : '—'}</p>
                            </div>
                            <div className="bg-neutral-900 rounded-xl p-5">
                                <p className="text-sm text-neutral-500 mb-1">Data size</p>
                                <p className="text-2xl font-semibold">{metrics?.total_data_size ? `${(metrics.total_data_size / 1024).toFixed(1)}KB` : '—'}</p>
                            </div>
                        </div>

                        {/* Status Codes */}
                        {metrics?.status_codes && Object.keys(metrics.status_codes).length > 0 && (
                            <div className="space-y-3">
                                <p className="text-sm text-neutral-500">Response codes</p>
                                <div className="flex gap-3 flex-wrap">
                                    {Object.entries(metrics.status_codes).map(([code, count]) => (
                                        <div
                                            key={code}
                                            className={`px-3 py-1.5 rounded-lg text-sm font-medium ${code.startsWith('2') ? 'bg-green-500/10 text-green-400' :
                                                    code.startsWith('4') ? 'bg-yellow-500/10 text-yellow-400' :
                                                        'bg-red-500/10 text-red-400'
                                                }`}
                                        >
                                            {code}: {count}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Actions */}
                        <div className="flex gap-4 pt-4">
                            {isComplete && (
                                <Button
                                    onClick={downloadResults}
                                    className="flex-1 h-14 bg-white text-black hover:bg-neutral-200 rounded-xl text-base font-medium"
                                >
                                    <Download className="w-5 h-5 mr-2" />
                                    Download results
                                </Button>
                            )}
                            <Button
                                onClick={resetJob}
                                variant="ghost"
                                className="h-14 px-6 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-xl"
                            >
                                {isComplete ? 'New extraction' : 'Cancel'}
                            </Button>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}
