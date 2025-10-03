import { useState, useEffect } from 'react';
import { Container } from '../components/layout/Container';
import { ConfigSection } from '../components/config/Section';
import { ConfigValue, ConfigSection as ConfigSectionType } from '../types/config';
import { generateConfigSchema } from '../lib/schema-loader';
import { Form, FloatingButton, SaveButton } from '../components/ui/form';
import { PageHeader } from '../components/layout/PageHeader';
import { useToast } from '../components/providers/toast';
import { API_BASE_URL, API_ENDPOINTS } from '../lib/const';

// the backend API base URL
const TITLE = 'Configuration';
const SAVE_BUTTON_TEXT = 'Save Configuration';

export default function ConfigurationPage() {
  const [config, setConfig] = useState<ConfigValue>({});
  const [configSchema, setConfigSchema] = useState<ConfigSectionType[]>([]);
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
    }
  };

  return (
    <Container>
      <PageHeader title={TITLE} />
      <Form onSubmit={handleSubmit}>
        {configSchema.map(section => (
          <ConfigSection
            key={section.name}
            name={section.name}
            options={section.options}
            values={config}
            onChange={handleChange}
          />
        ))}
        <FloatingButton>
          <SaveButton>{SAVE_BUTTON_TEXT}</SaveButton>
        </FloatingButton>
      </Form>
    </Container>
  );
}