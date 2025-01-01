import { ConfigOption } from "../../types/config";

interface ConfigInputProps {
    option: ConfigOption;
    value: any;
    onChange: (value: any) => void;
}

export function ConfigInput({ option, value, onChange }: ConfigInputProps) {
    const id = `${option.section}-${option.key}`;
    const baseClass = "rounded-md border border-border bg-background p-2";

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

        case 'numberArray':
            // Handle slots as checkboxes for numbers 1-4
            const slots = Array.isArray(value) ? value : [];
            return (
                <div className="flex gap-4">
                    {[1, 2, 3, 4].map(num => (
                        <label key={num} className="flex items-center gap-2">
                            <input
                                type="checkbox"
                                checked={slots.includes(num)}
                                onChange={e => {
                                    const newSlots = e.target.checked
                                        ? [...slots, num].sort((a, b) => a - b)
                                        : slots.filter(slot => slot !== num);
                                    onChange(newSlots);
                                }}
                                className="h-4 w-4 rounded border-border bg-background"
                            />
                            <span className="text-sm">{num}</span>
                        </label>
                    ))}
                </div>
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

        case 'number':
            return (
                <input
                    type="number"
                    id={id}
                    value={value ?? ''}
                    onChange={e => onChange(e.target.valueAsNumber)}
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