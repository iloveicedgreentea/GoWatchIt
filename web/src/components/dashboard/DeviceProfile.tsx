import { useState, useEffect } from 'react';
import { DashboardCard } from "./DashboardCard";
import { Badge } from "@/components/ui/badge";
import { API_BASE_URL } from '../../lib/const';

interface DeviceProfiles {
    [device: string]: string;
}

export function DeviceProfilesGrid() {
    const [profiles, setProfiles] = useState<DeviceProfiles>({});
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchProfiles = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/currentprofile`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data = await response.json();
                setProfiles(data);
                setError(null);
            } catch (err) {
                setError((err as Error).message);
            } finally {
                setLoading(false);
            }
        };

        fetchProfiles();
        // Set up polling every 5 seconds
        const intervalId = setInterval(fetchProfiles, 5000);

        return () => clearInterval(intervalId);
    }, []);

    if (error) {
        return (
            <DashboardCard title="Device Profiles">
                <div className="flex items-center space-x-2">
                    <Badge variant="destructive">Error</Badge>
                    <p className="text-sm text-muted-foreground">{error}</p>
                </div>
            </DashboardCard>
        );
    }

    if (loading) {
        return (
            <DashboardCard title="Device Profiles">
                <p className="text-sm text-muted-foreground">Loading profiles...</p>
            </DashboardCard>
        );
    }

    if (Object.keys(profiles).length === 0) {
        return (
            <DashboardCard title="Device Profiles">
                <p className="text-sm text-muted-foreground">No devices found</p>
            </DashboardCard>
        );
    }

    return (
        <div className="grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
            {Object.entries(profiles).map(([device, profile]) => (
                <DashboardCard key={device} title={device}>
                    <div className="space-y-2">
                        <div className="flex items-center space-x-2">
                            <Badge variant="secondary">Current Profile</Badge>
                            <span className="text-sm font-medium">{profile}</span>
                            Mute button here TODO
                        </div>
                    </div>
                </DashboardCard>
            ))}
        </div>
    );
}