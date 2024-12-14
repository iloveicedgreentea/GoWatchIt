import { NavLink } from "react-router-dom";
import { Layout, LayoutDashboard, Settings, Logs } from "lucide-react";
import { cn } from "../lib/utils";

// TODO: move this to a separate file
const navItems = [
  {
    title: "Dashboard",
    href: "/",
    icon: LayoutDashboard
  },
  //   TODO: add logs page
  {
    title: "Configuration",
    href: "/configuration",
    icon: Settings
  },
  {
    title: "Logs",
    href: "/logs",
    icon: Logs
  }
] as const;

export function Sidebar() {
  return (
    <div className="flex h-screen w-64 flex-col fixed border-r border-border bg-background">
      <div className="p-4">
        <div className="flex items-center gap-2 mb-6">
          <Layout className="h-6 w-6 text-primary" />
          <span className="text-lg font-semibold text-foreground">
            GoWatchIt
          </span>
        </div>

        <nav className="space-y-2">
          {navItems.map((item) => (
            <NavLink
              key={item.href}
              to={item.href}
              end={item.href === "/"}
              className={({ isActive }) =>
                cn(
                  "flex items-center w-full rounded-lg px-3 py-2 transition-colors",
                  isActive
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                )
              }
            >
              <item.icon className="mr-2 h-4 w-4" />
              {item.title}
            </NavLink>
          ))}
        </nav>
      </div>
    </div>
  );
}