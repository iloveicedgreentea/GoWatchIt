import { DeviceProfilesGrid } from "../components/dashboard/DeviceProfile";

export function Dashboard() {
    return (
        <div className="space-y-8">
            {/* Header */}
            <div>
                <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                    Dashboard
                </h1>
                <p className="text-base-content/60 mt-2">
                    Connected Devices
                </p>
            </div>

            {/* Device Profiles */}
            <DeviceProfilesGrid />
        </div>
    );
}