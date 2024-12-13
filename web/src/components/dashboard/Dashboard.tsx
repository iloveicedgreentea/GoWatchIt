import { Container } from "../layout/Container";
import { DashboardCard } from "./DashboardCard";
import { PageHeader } from "../layout/PageHeader";

export function Dashboard() {
    return (
        <Container>
            <PageHeader title="Dashboard" />

            <div className="grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                <DashboardCard title="Current BEQ Profile">
                    <p className="text-muted-foreground">Loading...</p>
                </DashboardCard>

                <DashboardCard title="Media Info">
                    <p className="text-muted-foreground">No media playing</p>
                </DashboardCard>

                <DashboardCard title="System Status">
                    <p className="text-muted-foreground">All systems operational</p>
                </DashboardCard>
            </div>
        </Container>
    );
}