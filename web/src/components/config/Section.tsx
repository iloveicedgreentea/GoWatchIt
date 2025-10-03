import { Card, CardHeader, CardTitle, CardContent } from "../ui/card";
import { ConfigOption, ConfigValue } from "../../types/config";
import { ConfigToggle } from "./Toggle";
import { ConfigField } from "./Field";

interface ConfigSectionProps {
    name: string;
    options: ConfigOption[];
    values: ConfigValue;
    onChange: (section: string, key: string, value: any) => void;
}

export function ConfigSection({ name, options, values, onChange }: ConfigSectionProps) {
    const isEnabled = Boolean(values[name]?.enabled);
    // Convert camelCase section name to readable title (e.g., hdmiSync -> Hdmi Sync)
    const displayName = name
        .replace(/([a-z])([A-Z])/g, '$1 $2')
        .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2')
        .replace(/^./, (str) => str.toUpperCase())
        .trim();

    return (
        <Card className="mb-6">
            <CardHeader>
                <CardTitle>{displayName}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <ConfigToggle
                    section={name}
                    enabled={isEnabled}
                    onChange={onChange}
                />

                {isEnabled && (
                    <div className="space-y-4 mt-4">
                        {options
                            .filter(opt => opt.key !== 'enabled')
                            .map(option => (
                                <ConfigField
                                    key={option.key}
                                    option={option}
                                    value={values[name]?.[option.key]}
                                    onChange={onChange}
                                />
                            ))}
                    </div>
                )}
            </CardContent>
        </Card>
    );
}