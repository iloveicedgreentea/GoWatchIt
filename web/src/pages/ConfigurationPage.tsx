import { useState, useEffect } from 'react';
import { Settings, Save, Loader2 } from 'lucide-react';
import { ConfigSection } from '../components/config/Section';
import { ConfigValue, ConfigSection as ConfigSectionType } from '../types/config';
import { generateConfigSchema } from '../lib/schema-loader';
import { useToast } from '../components/providers/toast';
import { API_BASE_URL, API_ENDPOINTS } from '../lib/const';

export default function ConfigurationPage() {
  const [config, setConfig] = useState<ConfigValue>({});
  const [configSchema, setConfigSchema] = useState<ConfigSectionType[]>([]);
  const [isSaving, setIsSaving] = useState(false);
  const { addToast } = useToast();

  // Load config schema
  useEffect(() => {
    generateConfigSchema().then(setConfigSchema);
  }, []);

  // get config
  useEffect(() => {
    fetch(`${API_BASE_URL}${API_ENDPOINTS.CONFIG}`)
      .then(async res => {
        const text = await res.text();
        console.log('Raw response:', text);
        console.log('Response status:', res.status);
        console.log('Response headers:', Object.fromEntries(res.headers.entries()));

        try {
          return JSON.parse(text);
        } catch (parseError) {
          console.error('JSON parse error:', parseError);
          throw new Error(`Invalid JSON response: ${text}`);
        }
      })
      .then(setConfig)
      .catch(error => {
        console.error('Error loading config:', error);
        addToast({
          title: 'Error',
          description: 'Failed to load configuration: ' + error.message,
          variant: 'destructive',
        });
      });
  }, [addToast]);

  const handleChange = (section: string, key: string, value: any) => {
    setConfig(prev => ({
      ...prev,
      [section]: {
        ...prev[section],
        [key]: value
      }
    }));
  };

  // submit changes
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.CONFIG}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || response.statusText);
      }

      addToast({
        title: 'Success',
        description: 'Configuration saved successfully',
        variant: 'success',
      });
    } catch (error) {
      addToast({
        title: 'Error',
        description: error instanceof Error ? error.message : 'Failed to save configuration',
        variant: 'destructive',
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
          Configuration
        </h1>
        <p className="text-base-content/60 mt-2">
          Configure your application settings and integrations
        </p>
      </div>

      {/* Configuration Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {configSchema.map(section => (
          <ConfigSection
            key={section.name}
            name={section.name}
            options={section.options}
            values={config}
            onChange={handleChange}
          />
        ))}

        {/* Save Button */}
        <div className="sticky bottom-8 flex justify-end">
          <button
            type="submit"
            disabled={isSaving}
            className="btn btn-primary btn-lg gap-2 shadow-xl"
          >
            {isSaving ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Saving...
              </>
            ) : (
              <>
                <Save className="w-5 h-5" />
                Save Configuration
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
}