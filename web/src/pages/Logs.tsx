import { useState, useEffect } from 'react';
import { useToast } from '../components/providers/toast';
import { AlertCircle, FileText } from 'lucide-react';
import type { LogEntry } from '../types/logs';
import { API_BASE_URL, API_ENDPOINTS } from '../lib/const';

const REFRESH_INTERVAL = 1000; // 1 second

export default function Logs() {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const { addToast, isConnected, handleConnectionChange } = useToast();
    const [lastError, setLastError] = useState<string | null>(null);

    const fetchLogs = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.LOGS}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            setLogs(data);

            // Reset error state on successful fetch
            if (lastError) {
                setLastError(null);
            }

            // Signal successful connection
            handleConnectionChange(false);

        } catch (error) {
            const errorMessage = (error as Error).message;

            // Only show error toast if it's a new error
            if (errorMessage !== lastError) {
                setLastError(errorMessage);
                addToast({
                    title: 'Error',
                    description: 'Failed to load logs: ' + errorMessage,
                    variant: 'destructive',
                });
            }

            // Signal connection error
            handleConnectionChange(true);
        }
    };

    useEffect(() => {
        fetchLogs();

        const intervalId = setInterval(fetchLogs, isConnected ? REFRESH_INTERVAL : 30000);

        return () => clearInterval(intervalId);
    }, [isConnected]); // Add isConnected to dependency array

    const getLevelClass = (level: string): string => {
        switch (level) {
            case 'ERROR': return 'badge-error';
            case 'WARN': return 'badge-warning';
            case 'INFO': return 'badge-info';
            case 'DEBUG': return 'badge-ghost';
            default: return 'badge-ghost';
        }
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                        Logs
                    </h1>
                    <p className="text-base-content/60 mt-2">
                        Real-time application logs and events
                    </p>
                </div>
                {!isConnected && (
                    <div className="badge badge-error gap-2">
                        <AlertCircle className="w-4 h-4" />
                        Disconnected
                    </div>
                )}
            </div>

            {/* Logs Container */}
            <div className="card bg-base-200 shadow-xl">
                <div className="card-body p-0">
                    <div className="overflow-x-auto max-h-[calc(100vh-16rem)]">
                        <table className="table table-zebra table-pin-rows">
                            <thead>
                                <tr>
                                    <th className="w-24">Level</th>
                                    <th className="w-48">Time</th>
                                    <th>Message</th>
                                    <th className="w-64">Source</th>
                                </tr>
                            </thead>
                            <tbody>
                                {logs.length === 0 ? (
                                    <tr>
                                        <td colSpan={4} className="text-center py-8">
                                            <div className="flex flex-col items-center gap-2 opacity-60">
                                                <FileText className="w-12 h-12" />
                                                <p>No logs available</p>
                                            </div>
                                        </td>
                                    </tr>
                                ) : (
                                    logs.map((log, index) => {
                                        // Extract source info from Extra if it exists
                                        let sourceFromExtra = '';
                                        if (log.Extra?.source) {
                                            sourceFromExtra = typeof log.Extra.source === 'string'
                                                ? log.Extra.source
                                                : JSON.stringify(log.Extra.source);
                                        } else if (log.Extra?.caller) {
                                            sourceFromExtra = typeof log.Extra.caller === 'string'
                                                ? log.Extra.caller
                                                : JSON.stringify(log.Extra.caller);
                                        } else if (log.Extra?.file) {
                                            const file = log.Extra.file;
                                            const line = log.Extra.line || '';
                                            sourceFromExtra = `${file}${line ? ':' + line : ''}`;
                                        }

                                        // Determine which source to display
                                        const hasValidSource = log.source?.file && log.source?.file !== '';
                                        const displaySource = hasValidSource
                                            ? `${log.source.file}:${log.source.line}`
                                            : sourceFromExtra;

                                        // Filter out source-related fields from Extra for message display
                                        const filteredExtra = log.Extra ? Object.entries(log.Extra)
                                            .filter(([key]) => !['source', 'caller', 'file', 'line'].includes(key))
                                            : [];

                                        return (
                                            <tr key={index} className="hover">
                                                <td>
                                                    <div className={`badge ${getLevelClass(log.level)} badge-sm`}>
                                                        {log.level}
                                                    </div>
                                                </td>
                                                <td className="text-sm opacity-70">
                                                    {new Date(log.time).toLocaleString()}
                                                </td>
                                                <td>
                                                    <div className="text-sm">
                                                        {log.msg}
                                                        {filteredExtra.length > 0 && (
                                                            <div className="mt-1 space-y-1">
                                                                {filteredExtra.map(([key, value]) => (
                                                                    <div key={key} className="text-xs opacity-60">
                                                                        <span className="font-semibold">{key}:</span> {JSON.stringify(value)}
                                                                    </div>
                                                                ))}
                                                            </div>
                                                        )}
                                                    </div>
                                                </td>
                                                <td>
                                                    {displaySource && (
                                                        <code className="text-xs bg-base-300 px-2 py-1 rounded block w-fit">
                                                            {displaySource}
                                                        </code>
                                                    )}
                                                </td>
                                            </tr>
                                        );
                                    })
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
}