import {
    Toast,
    ToastClose,
    ToastDescription,
    ToastProvider,
    ToastTitle,
    ToastViewport,
} from "../ui/toast";
import { createContext, useContext, useState, useEffect, useCallback } from "react";

interface Toast {
    id: string;
    title?: string;
    description: string;
    variant?: "default" | "success" | "destructive";
    timestamp: number;
}

interface ToastContextType {
    addToast: (toast: Omit<Toast, "id" | "timestamp">) => void;
    isConnected: boolean;
    handleConnectionChange: (isError: boolean) => void;
}

const MAX_TOASTS = 3;
const INITIAL_RETRY_DELAY = 1000;
const MAX_RETRY_DELAY = 30000;

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastContextProvider({
    children,
}: {
    children: React.ReactNode;
}) {
    const [toasts, setToasts] = useState<Toast[]>([]);
    const [isConnected, setIsConnected] = useState(true);
    const [_, setRetryDelay] = useState<number>(INITIAL_RETRY_DELAY);

    // Remove old toasts when we exceed the limit
    useEffect(() => {
        if (toasts.length > MAX_TOASTS) {
            const sortedToasts = [...toasts].sort((a, b) => b.timestamp - a.timestamp);
            setToasts(sortedToasts.slice(0, MAX_TOASTS));
        }
    }, [toasts]);

    const addToast = useCallback((toast: Omit<Toast, "id" | "timestamp">) => {
        const id = Math.random().toString(36).slice(2);
        const timestamp = Date.now();

        setToasts((prev) => {
            const newToasts = [...prev, { ...toast, id, timestamp }];
            // Sort by timestamp so newest appears at the top
            return newToasts.sort((a, b) => b.timestamp - a.timestamp);
        });
    }, []);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
    }, []);

    // Handle connection status changes
    const handleConnectionChange = useCallback((isError: boolean) => {
        if (isError && isConnected) {
            setIsConnected(false);
            // Implement exponential backoff
            setRetryDelay((prev) => Math.min(prev * 2, MAX_RETRY_DELAY));
        } else if (!isError && !isConnected) {
            setIsConnected(true);
            setRetryDelay(INITIAL_RETRY_DELAY);
            // Show reconnected toast
            addToast({
                title: "Connected",
                description: "Successfully reconnected to the logs service",
                variant: "success",
            });
        }
    }, [isConnected, addToast]);

    return (
        <ToastContext.Provider value={{ addToast, isConnected, handleConnectionChange }}>
            <ToastProvider>
                {children}
                {toasts.map(({ id, title, description, variant }, index) => (
                    <Toast
                        key={id}
                        variant={variant}
                        className={`transition-all duration-200 ${
                            // Adjust positioning for stacking effect
                            index > 0 ? `translate-y-${index * 2}` : ""
                            }`}
                        onOpenChange={(open) => {
                            if (!open) removeToast(id);
                        }}
                    >
                        <div className="grid gap-1">
                            {title && <ToastTitle>{title}</ToastTitle>}
                            {description && <ToastDescription>{description}</ToastDescription>}
                        </div>
                        <ToastClose />
                    </Toast>
                ))}
                <ToastViewport />
            </ToastProvider>
        </ToastContext.Provider>
    );
}

export function useToast() {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error("useToast must be used within ToastContextProvider");
    }
    return context;
}