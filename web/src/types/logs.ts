// types/logs.ts
export interface LogEntry {
    timestamp: string;
    level: 'INFO' | 'WARN' | 'ERROR' | 'DEBUG';
    msg: string;
    source: Source;
    error?: string;
    Extra: Record<string, unknown>; 
}

interface Source {
    function: string;
    line: number;
    file: string;
}