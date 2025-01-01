import { useState, useEffect } from 'react';
import { Container } from '../components/layout/Container';
import { PageHeader } from '../components/layout/PageHeader';
import { useToast } from '../components/providers/toast';
import { Card, CardContent } from '../components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Badge } from '@/components/ui/badge';
import type { LogEntry } from '../types/logs';
import { API_BASE_URL } from '../lib/const';

const TITLE = 'Logs';
const REFRESH_INTERVAL = 1000; // 1 second

export default function Logs() {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const { addToast, isConnected, handleConnectionChange } = useToast();
    const [lastError, setLastError] = useState<string | null>(null);

    const fetchLogs = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/logs`);
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

    const getLevelColor = (level: string): "destructive" | "default" | "secondary" | "outline" => {
        switch (level) {
            case 'ERROR': return 'destructive';
            case 'WARN': return 'secondary';
            case 'INFO': return 'default';
            case 'DEBUG': return 'secondary';
            default: return 'default';
        }
    };

    return (
        <Container>
            <PageHeader title={TITLE} />
            {!isConnected && (
                <Badge variant="destructive">Disconnected</Badge>
            )}
            <Card>
                <CardContent className="p-4">
                    <ScrollArea className="h-[calc(100vh-8rem)] w-full rounded-md border">
                        {logs.map((log, index) => (
                            <div key={index} className="p-4 border-b last:border-0">
                                <div className="flex items-center gap-2 mb-1">
                                    <Badge variant={getLevelColor(log.level)}>
                                        {log.level.toUpperCase()}
                                    </Badge>
                                    <span className="text-sm text-muted-foreground">
                                        {new Date(log.timestamp).toLocaleString()}
                                    </span>
                                </div>
                                <p className="text-sm">
                                    {log.msg}
                                    {log.Extra && Object.entries(log.Extra).map(([key, value]) => (
                                        <p key={key} className="">
                                            {key}: {JSON.stringify(value)}
                                        </p>
                                    ))}
                                </p>
                                {log.source && (
                                    <pre className="mt-2 text-xs bg-muted p-2 rounded-md overflow-x-auto">
                                        Source: {log.source.file}:{log.source.line}
                                    </pre>
                                )}
                            </div>
                        ))}
                    </ScrollArea>
                </CardContent>
            </Card>
        </Container>
    );
}