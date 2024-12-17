import { useState, useEffect } from 'react';
import { Container } from '../components/layout/Container';
import { ConfigSection } from '../components/config/Section';
import { ConfigValue } from '../types/config';
import { CONFIG_SCHEMA } from '../types/configOptions';
import { Form, FloatingButton, SaveButton } from '../components/ui/form';
import { PageHeader } from '../components/layout/PageHeader';
import { useToast } from '../components/providers/toast';
import { API_BASE_URL } from '../lib/const';

// the backend API base URL
const TITLE = 'Configuration';
const SAVE_BUTTON_TEXT = 'Save Configuration';

export default function ConfigurationPage() {
  const [config, setConfig] = useState<ConfigValue>({});
  const { addToast } = useToast();

  // get config
  useEffect(() => {
    fetch(`${API_BASE_URL}/config`)
      .then(res => res.json())
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
      const response = await fetch(`${API_BASE_URL}/config`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
      });

      if (!response.ok) throw new Error(response.statusText);

      addToast({
        title: 'Success',
        description: 'Configuration saved successfully',
        variant: 'success',
      });
    } catch (error) {
      addToast({
        title: 'Error',
        description: 'Failed to save configuration',
        variant: 'destructive',
      });
    }
  };

  return (
    <Container>
      <PageHeader title={TITLE} />
      <Form onSubmit={handleSubmit}>
        {CONFIG_SCHEMA.map(section => (
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