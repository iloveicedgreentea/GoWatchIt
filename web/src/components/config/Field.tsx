import { ConfigOption } from "../../types/config";
import { ConfigInput } from "./Input";

interface ConfigFieldProps {
  option: ConfigOption;
  value: any;
  onChange: (section: string, key: string, value: any) => void;
}

export function ConfigField({ option, value, onChange }: ConfigFieldProps) {
  const handleChange = (newValue: any) => {
    onChange(option.section, option.key, newValue);
  };

  return (
    <div className="space-y-2">
      <label 
        htmlFor={`${option.section}-${option.key}`}
        className="block text-sm font-medium"
      >
        {option.label}
      </label>
      {option.description && (
        <p className="text-sm text-muted-foreground">
          {option.description}
        </p>
      )}
      <ConfigInput 
        option={option} 
        value={value} 
        onChange={handleChange} 
      />
    </div>
  );
}