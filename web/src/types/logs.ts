// types/logs.ts
export interface LogEntry {
    time: string;
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