import { useState } from "react";
import { ConfigOption } from "../../types/config";

interface ConfigInputProps {
    option: ConfigOption;
    value: any;
    onChange: (value: any) => void;
}

export function ConfigInput({ option, value, onChange }: ConfigInputProps) {
    const id = `${option.section}-${option.key}`;
    const [error, setError] = useState<string | null>(null);
    const baseClass = "rounded-md border border-border bg-background p-2";
    const errorClass = error ? "border-red-500" : "";

    // Validate input and update state
    const handleTextChange = (newValue: string) => {
        if (option.validator) {
            const result = option.validator(newValue);
            if (result === true) {
                setError(null);
            } else {
                setError(result);
            }
        }
        onChange(newValue);
    };

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

        case 'stringArray':
            // Handle string arrays as checkboxes
            const selectedItems = Array.isArray(value) ? value : [];
            if (!option.options) return null;
            return (
                <div className="flex flex-col gap-2">
                    {option.options.map(opt => (
                        <label key={opt.value} className="flex items-center gap-2">
                            <input
                                type="checkbox"
                                checked={selectedItems.includes(opt.value)}
                                onChange={e => {
                                    const newItems = e.target.checked
                                        ? [...selectedItems, opt.value]
                                        : selectedItems.filter(item => item !== opt.value);
                                    onChange(newItems);
                                }}
                                className="h-4 w-4 rounded border-border bg-background"
                            />
                            <span className="text-sm">{opt.label}</span>
                        </label>
                    ))}
                </div>
            );

        case 'password':
            return (
                <div className="flex flex-col gap-1">
                    <input
                        type="password"
                        id={id}
                        value={String(value ?? '')}
                        onChange={e => handleTextChange(e.target.value)}
                        placeholder={option.placeholder}
                        className={`${baseClass} ${errorClass}`}
                    />
                    {error && <span className="text-xs text-red-500">{error}</span>}
                </div>
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
                <div className="flex flex-col gap-1">
                    <input
                        type="text"
                        id={id}
                        value={String(value ?? '')}
                        onChange={e => handleTextChange(e.target.value)}
                        placeholder={option.placeholder}
                        className={`${baseClass} ${errorClass}`}
                    />
                    {error && <span className="text-xs text-red-500">{error}</span>}
                </div>
            );
    }
}