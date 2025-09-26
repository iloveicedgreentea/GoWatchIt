import { ConfigOption } from "../../types/config";
import { ConfigInput } from "./Input";
import { Separator } from "../ui/seperator";

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
      <div>
        <div className="space-y-2 py-4">
          <div className="flex justify-between items-center">
            <div>
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
            </div>
            <ConfigInput
              option={option}
              value={value}
              onChange={handleChange}
            />
          </div>
        </div>
        <Separator className="mt-2" />
      </div>
    );
  }