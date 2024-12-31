import { Container } from "../components/layout/Container";
import { PageHeader } from "../components/layout/PageHeader";
import { DeviceProfilesGrid } from "../components/dashboard/DeviceProfile";
import { DashboardCard } from "../components/dashboard/DashboardCard";

export function Dashboard() {
    return (
        <Container>
            <PageHeader title="Dashboard" />
            
            <div className="space-y-6">
                <DeviceProfilesGrid />

                <div className="grid gap-6 grid-cols-1 md:grid-cols-2">
                    <DashboardCard title="Media Info">
                        <p className="text-muted-foreground">No media playing</p>
                    </DashboardCard>

                    <DashboardCard title="System Status">
                        <p className="text-muted-foreground">All systems operational</p>
                    </DashboardCard>
                </div>
            </div>
        </Container>
    );
}