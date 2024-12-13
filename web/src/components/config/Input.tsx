import { ConfigOption } from "../../types/config";

interface ConfigInputProps {
  option: ConfigOption;
  value: any;
  onChange: (value: any) => void;
}

export function ConfigInput({ option, value, onChange }: ConfigInputProps) {
  const id = `${option.section}-${option.key}`;
  const baseClass = "w-full rounded-md border border-border bg-background p-2";

  switch (option.type) {
    case 'checkbox':
      return (
        <input
          type="checkbox"
          id={id}
          checked={Boolean(value)}
          onChange={e => onChange(e.target.checked)}
          className="h-4 w-4 rounded border-border bg-background"
        />
      );
    case 'select':
      return (
        <select
          id={id}
          value={String(value ?? '')}
          onChange={e => onChange(e.target.value)}
          className={baseClass}
        >
          {option.options?.map(opt => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      );
    case 'password':
      return (
        <input
          type="password"
          id={id}
          value={String(value ?? '')}
          onChange={e => onChange(e.target.value)}
          placeholder={option.placeholder}
          className={baseClass}
        />
      );
    default:
      return (
        <input
          type="text"
          id={id}
          value={String(value ?? '')}
          onChange={e => onChange(e.target.value)}
          placeholder={option.placeholder}
          className={baseClass}
        />
      );
  }
}