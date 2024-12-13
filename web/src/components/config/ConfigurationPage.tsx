import { useState, useEffect } from 'react';
import { Container } from '../layout/Container';
import { ConfigSection } from './Section';
import { CONFIG_SCHEMA, ConfigValue } from '../../types/config';
import { Form, FloatingButton, SaveButton } from '../ui/form';
import { PageHeader } from '../layout/PageHeader';
import { useToast } from '../providers/toast';

const API_BASE_URL = 'http://localhost:9999';

export default function ConfigurationPage() {
  const [config, setConfig] = useState<ConfigValue>({});
  const { addToast } = useToast();

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
      <PageHeader title="Configuration" />
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
          <SaveButton>Save Configuration</SaveButton>
        </FloatingButton>
      </Form>
    </Container>
  );
}