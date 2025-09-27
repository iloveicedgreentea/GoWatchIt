interface ConfigToggleProps {
    section: string;
    enabled: boolean;
    onChange: (section: string, key: string, value: boolean) => void;
  }
  
  export function ConfigToggle({ section, enabled, onChange }: ConfigToggleProps) {
    return (
      <div className="flex items-center gap-4">
        <input
          type="checkbox"
          id={`${section}-enabled`}
          checked={enabled}
          onChange={e => onChange(section, 'enabled', e.target.checked)}
          className="h-4 w-4 rounded border-border bg-background"
        />
        <label 
          htmlFor={`${section}-enabled`} 
          className="text-sm font-medium leading-none"
        >
          Enable {section.charAt(0).toUpperCase() + section.slice(1)}
        </label>
      </div>
    );
  }