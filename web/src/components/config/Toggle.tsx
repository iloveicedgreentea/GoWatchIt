interface ConfigToggleProps {
    section: string;
    enabled: boolean;
    onChange: (section: string, key: string, value: boolean) => void;
  }

  export function ConfigToggle({ section, enabled, onChange }: ConfigToggleProps) {
    const displayName = section
      .replace(/([a-z])([A-Z])/g, '$1 $2')
      .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2')
      .replace(/^./, (str) => str.toUpperCase())
      .trim();

    return (
      <div className="form-control bg-base-300 p-4 rounded-lg">
        <label className="label cursor-pointer justify-start gap-4">
          <input
            type="checkbox"
            id={`${section}-enabled`}
            checked={enabled}
            onChange={e => onChange(section, 'enabled', e.target.checked)}
            className="toggle toggle-primary toggle-lg border-2 border-base-content/20"
          />
          <div className="flex flex-col">
            <span className="label-text font-bold text-base">
              Enable {displayName}
            </span>
            <span className="label-text-alt opacity-70">
              {enabled ? 'Currently enabled' : 'Currently disabled'}
            </span>
          </div>
          <span className={`ml-auto badge ${enabled ? 'badge-success' : 'badge-ghost'}`}>
            {enabled ? 'ON' : 'OFF'}
          </span>
        </label>
      </div>
    );
  }