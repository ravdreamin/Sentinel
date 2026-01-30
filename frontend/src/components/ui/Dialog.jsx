import * as React from 'react';
import { X } from 'lucide-react';

export function Dialog({ open, onOpenChange, children }) {
    if (!open) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-black/60 backdrop-blur-sm"
                onClick={() => onOpenChange(false)}
            />
            {/* Content */}
            <div className="relative z-10 w-full max-w-md mx-4 animate-in fade-in zoom-in-95 duration-200">
                {children}
            </div>
        </div>
    );
}

export function DialogContent({ children, className = '' }) {
    return (
        <div className={`bg-neutral-900 border border-neutral-800 rounded-2xl p-6 shadow-2xl ${className}`}>
            {children}
        </div>
    );
}

export function DialogHeader({ children }) {
    return <div className="mb-4">{children}</div>;
}

export function DialogTitle({ children }) {
    return <h2 className="text-xl font-semibold text-white">{children}</h2>;
}

export function DialogDescription({ children }) {
    return <p className="text-neutral-400 mt-1">{children}</p>;
}

export function DialogFooter({ children }) {
    return <div className="flex justify-end gap-3 mt-6">{children}</div>;
}
