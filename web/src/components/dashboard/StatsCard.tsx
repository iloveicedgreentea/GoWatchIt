import { LucideIcon } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  description?: string;
}

export function StatsCard({ title, value, icon: Icon, trend, description }: StatsCardProps) {
  return (
    <div className="stats shadow bg-base-200">
      <div className="stat">
        <div className="stat-figure text-primary">
          <Icon className="w-8 h-8" />
        </div>
        <div className="stat-title">{title}</div>
        <div className="stat-value text-primary">{value}</div>
        {description && <div className="stat-desc">{description}</div>}
        {trend && (
          <div className={`stat-desc ${trend.isPositive ? 'text-success' : 'text-error'}`}>
            {trend.isPositive ? '↗︎' : '↘︎'} {Math.abs(trend.value)}%
          </div>
        )}
      </div>
    </div>
  );
}
