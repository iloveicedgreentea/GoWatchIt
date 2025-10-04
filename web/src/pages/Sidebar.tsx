import { NavLink } from "react-router-dom";
import { Layout, LayoutDashboard, Settings, FileText } from "lucide-react";

const navItems = [
  {
    title: "Configuration",
    href: "/configuration",
    icon: Settings
  },
  {
    title: "Dashboard",
    href: "/",
    icon: LayoutDashboard
  },
  {
    title: "Logs",
    href: "/logs",
    icon: FileText
  }
] as const;

export function Sidebar() {
  return (
    <div className="flex h-screen w-64 flex-col fixed border-r border-base-300 bg-base-200">
      {/* Header */}
      <div className="p-4 border-b border-base-300">
        <div className="flex items-center gap-2 mb-6">
          <Layout className="h-6 w-6 text-primary" />
          <span className="text-xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            GoWatchIt
          </span>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4">
        <ul className="menu menu-lg gap-2">
          {navItems.map((item) => (
            <li key={item.href}>
              <NavLink
                to={item.href}
                end={item.href === "/"}
                className={({ isActive }) =>
                  isActive ? "active" : ""
                }
              >
                <item.icon className="h-5 w-5" />
                {item.title}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>

      {/* Footer */}
      <div className="p-4 border-t border-base-300">
        <div className="text-xs text-base-content/50 text-center">
          v0.1.0
        </div>
      </div>
    </div>
  );
}