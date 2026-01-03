import React, { useEffect } from 'react';
import { CheckCircleIcon, XCircleIcon, XMarkIcon } from '@heroicons/react/24/solid';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastProps {
    id: string;
    message: string;
    type: ToastType;
    onClose: (id: string) => void;
    duration?: number;
}

const Toast: React.FC<ToastProps> = ({ id, message, type, onClose, duration = 3000 }) => {
    useEffect(() => {
        const timer = setTimeout(() => {
            onClose(id);
        }, duration);

        return () => clearTimeout(timer);
    }, [id, duration, onClose]);

    const styles = {
        success: 'bg-green-500/10 border-green-500/20 text-green-200',
        error: 'bg-red-500/10 border-red-500/20 text-red-200',
        info: 'bg-indigo-500/10 border-indigo-500/20 text-indigo-200',
    };

    const icons = {
        success: <CheckCircleIcon className="w-5 h-5 text-green-500" />,
        error: <XCircleIcon className="w-5 h-5 text-red-500" />,
        info: <CheckCircleIcon className="w-5 h-5 text-indigo-500" />,
    };

    return (
        <div className={`flex items-center gap-3 px-4 py-3 rounded-xl border backdrop-blur-md shadow-2xl animate-slide-in mb-3 w-80 ${styles[type]}`}>
            <div className="flex-shrink-0">
                {icons[type]}
            </div>
            <p className="flex-1 text-sm font-medium">{message}</p>
            <button
                onClick={() => onClose(id)}
                className="p-1 rounded-full hover:bg-white/10 transition-colors"
            >
                <XMarkIcon className="w-4 h-4 opacity-60 hover:opacity-100" />
            </button>
        </div>
    );
};

export default Toast;
