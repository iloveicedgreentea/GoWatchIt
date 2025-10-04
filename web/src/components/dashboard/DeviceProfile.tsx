import { useState, useEffect } from 'react';
import { Tv, AlertCircle, Loader2 } from 'lucide-react';
import { API_BASE_URL, API_ENDPOINTS } from '../../lib/const';

interface DeviceProfiles {
    [device: string]: string;
}

interface ProfileResponse {
    profile: DeviceProfiles;
}

export function DeviceProfilesGrid() {
    const [profiles, setProfiles] = useState<DeviceProfiles>({});
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchProfiles = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.PROFILE}`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data: ProfileResponse = await response.json();
                setProfiles(data.profile || {});
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
            <div className="alert alert-error">
                <AlertCircle className="w-5 h-5" />
                <span>Failed to load device profiles: {error}</span>
            </div>
        );
    }

    if (loading) {
        return (
            <div className="flex items-center justify-center p-12">
                <Loader2 className="w-8 h-8 animate-spin text-primary" />
                <span className="ml-3 text-lg">Loading profiles...</span>
            </div>
        );
    }

    if (Object.keys(profiles).length === 0) {
        return (
            <div className="alert alert-info">
                <AlertCircle className="w-5 h-5" />
                <span>No devices found. Make sure your devices are connected and configured.</span>
            </div>
        );
    }

    return (
        <div className="grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
            {Object.entries(profiles).map(([device, profile]) => (
                <div key={device} className="card bg-base-200 shadow-xl border border-base-300 transition-all duration-300">
                    <div className="card-body">
                        {/* Card Title with Icon */}
                        <div className="flex items-center gap-3 mb-4">
                            <div className="bg-gradient-to-br from-primary to-secondary p-3 rounded-lg">
                                <Tv className="w-6 h-6 text-white" />
                            </div>
                            <h2 className="card-title text-xl">{device}</h2>
                        </div>

                        <div className="divider my-0"></div>

                        {/* Active Profile Section */}
                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <span className="text-sm font-medium text-base-content/60">Active Profile</span>
                                <div className="badge badge-sm badge-success gap-1">
                                    <span className="w-1.5 h-1.5 rounded-full bg-white animate-pulse"></span>
                                    Online
                                </div>
                            </div>

                            <div className="bg-gradient-to-r from-primary/10 to-secondary/10 border border-primary/20 rounded-lg p-4">
                                <p className="text-lg font-semibold text-primary">
                                    {profile || 'Empty'}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    );
}