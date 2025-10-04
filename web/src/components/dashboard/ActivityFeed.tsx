import { Activity, Settings, Play } from 'lucide-react';

interface ActivityItem {
  id: string;
  type: 'profile_change' | 'playback' | 'config';
  message: string;
  timestamp: Date;
}

interface ActivityFeedProps {
  activities?: ActivityItem[];
}

export function ActivityFeed({ activities = [] }: ActivityFeedProps) {
  const getIcon = (type: ActivityItem['type']) => {
    switch (type) {
      case 'profile_change':
        return <Activity className="w-4 h-4" />;
      case 'playback':
        return <Play className="w-4 h-4" />;
      case 'config':
        return <Settings className="w-4 h-4" />;
    }
  };

  const getColor = (type: ActivityItem['type']) => {
    switch (type) {
      case 'profile_change':
        return 'badge-primary';
      case 'playback':
        return 'badge-secondary';
      case 'config':
        return 'badge-accent';
    }
  };

  return (
    <div className="card bg-base-200 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">
          <Activity className="w-5 h-5" />
          Recent Activity
        </h2>
        {activities.length === 0 ? (
          <p className="text-base-content/60">No recent activity</p>
        ) : (
          <ul className="space-y-3 mt-4">
            {activities.slice(0, 5).map((activity) => (
              <li key={activity.id} className="flex items-start gap-3">
                <div className={`badge ${getColor(activity.type)} gap-2`}>
                  {getIcon(activity.type)}
                </div>
                <div className="flex-1">
                  <p className="text-sm">{activity.message}</p>
                  <p className="text-xs text-base-content/50">
                    {activity.timestamp.toLocaleString()}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
