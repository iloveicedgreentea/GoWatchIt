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
        <div className="collapse bg-base-200 shadow-xl transition-all border-2 border-base-300 group">
            <input type="checkbox" defaultChecked className="peer" />
            <div className="collapse-title text-xl font-semibold flex items-center gap-3 pr-12 relative">
                <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                    <svg
                        className="w-5 h-5 text-primary transition-transform duration-200 peer-checked:group-[]:rotate-180"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                </div>
                {displayName}
                
            </div>
            <div className="collapse-content">
                <div className="space-y-4 pt-4">
                    <ConfigToggle
                        section={name}
                        enabled={isEnabled}
                        onChange={onChange}
                    />

                    {isEnabled && (
                        <div className="divider"></div>
                    )}

                    {isEnabled && (
                        <div className="space-y-4">
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
                </div>
            </div>
        </div>
    );
}