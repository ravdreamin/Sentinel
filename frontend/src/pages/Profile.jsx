import { useState, useEffect } from 'react';
import api from '../lib/api';
import { Button } from '../components/ui/Button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '../components/ui/Dialog';
import { FileText, Trash2, Download, AlertTriangle } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';

export default function Profile() {
    const { user, isGuest } = useAuth();
    const [jobs, setJobs] = useState([]);
    const [deleteDialog, setDeleteDialog] = useState({ open: false, filename: null });
    const [deleting, setDeleting] = useState(false);

    useEffect(() => {
        if (!isGuest && user) fetchJobs();
    }, [user, isGuest]);

    const fetchJobs = async () => {
        try {
            const { data } = await api.get('/api/jobs');
            setJobs(data.jobs || []);
        } catch (error) {
            console.error("Failed to fetch jobs", error);
        }
    };

    const handleDelete = async () => {
        if (!deleteDialog.filename) return;
        setDeleting(true);
        try {
            await api.delete(`/api/jobs/${deleteDialog.filename}`);
            setDeleteDialog({ open: false, filename: null });
            fetchJobs();
        } catch (error) {
            console.error("Delete error:", error);
        } finally {
            setDeleting(false);
        }
    };

    const handleDownload = async (filename) => {
        const response = await api.get(`/api/jobs/${filename}/download`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `${filename}_results.json`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    };

    if (isGuest) {
        return (
            <div className="min-h-screen bg-neutral-950 text-white font-sans flex items-center justify-center">
                <div className="text-center space-y-4">
                    <p className="text-neutral-400">Guest users don't have saved history.</p>
                    <Link to="/register">
                        <Button className="bg-white text-black hover:bg-neutral-200 rounded-full px-6">
                            Create an account
                        </Button>
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-neutral-950 text-white font-sans">
            {/* Header */}
            <header className="px-8 py-6 flex justify-between items-center border-b border-neutral-800">
                <div className="flex items-center gap-8">
                    <Link to="/" className="text-xl font-semibold tracking-tight hover:opacity-80 transition-opacity">
                        Sentinel
                    </Link>
                    <nav className="flex gap-6 text-sm text-neutral-400">
                        <Link to="/dashboard" className="hover:text-white transition-colors">Dashboard</Link>
                        <Link to="/profile" className="text-white">History</Link>
                    </nav>
                </div>
            </header>

            <main className="max-w-2xl mx-auto px-8 py-16">
                <div className="space-y-8">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight mb-2">Extraction history</h1>
                        <p className="text-neutral-400">Your saved extractions ({jobs.length}/10)</p>
                    </div>

                    {jobs.length === 0 ? (
                        <div className="text-center py-16 text-neutral-500">
                            <FileText className="w-12 h-12 mx-auto mb-4 opacity-50" />
                            <p>No extractions yet</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {jobs.map((job) => (
                                <div key={job} className="bg-neutral-900 rounded-xl p-5 flex items-center justify-between">
                                    <div className="flex items-center gap-4">
                                        <FileText className="w-5 h-5 text-neutral-500" />
                                        <span className="font-medium truncate max-w-xs">{job}</span>
                                    </div>
                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => handleDownload(job)}
                                            className="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg transition-colors"
                                            title="Download"
                                        >
                                            <Download className="w-4 h-4" />
                                        </button>
                                        <button
                                            onClick={() => setDeleteDialog({ open: true, filename: job })}
                                            className="p-2 text-neutral-400 hover:text-red-400 hover:bg-neutral-800 rounded-lg transition-colors"
                                            title="Delete"
                                        >
                                            <Trash2 className="w-4 h-4" />
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    <Link to="/dashboard">
                        <Button className="w-full h-14 bg-white text-black hover:bg-neutral-200 rounded-xl text-base font-medium">
                            New extraction
                        </Button>
                    </Link>
                </div>
            </main>

            {/* Delete Confirmation Dialog */}
            <Dialog open={deleteDialog.open} onOpenChange={(open) => setDeleteDialog({ ...deleteDialog, open })}>
                <DialogContent>
                    <DialogHeader>
                        <div className="flex items-center gap-3 mb-2">
                            <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                                <AlertTriangle className="w-5 h-5 text-red-500" />
                            </div>
                            <DialogTitle>Delete extraction</DialogTitle>
                        </div>
                        <DialogDescription>
                            Are you sure you want to delete this extraction? This action cannot be undone.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="bg-neutral-800 rounded-lg p-3 text-sm text-neutral-300 font-mono truncate">
                        {deleteDialog.filename}
                    </div>
                    <DialogFooter>
                        <Button
                            variant="ghost"
                            onClick={() => setDeleteDialog({ open: false, filename: null })}
                            className="text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg"
                        >
                            Cancel
                        </Button>
                        <Button
                            onClick={handleDelete}
                            disabled={deleting}
                            className="bg-red-600 hover:bg-red-700 text-white rounded-lg"
                        >
                            {deleting ? 'Deleting...' : 'Delete'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}
