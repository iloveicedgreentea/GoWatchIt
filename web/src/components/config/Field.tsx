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
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">{option.label}</span>
        </label>
        {option.description && (
          <p className="text-xs text-base-content/60 mb-2">
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