import { useState, useEffect } from 'react';
import { cn } from "../../lib/utils";
import { Alert } from '../ui/alert';
import { Container } from '../layout/Container';
import { ConfigSection } from './Section';
import { CONFIG_SCHEMA, ConfigValue, Notification } from '../../types/config';
import { Form, FloatingButton, SaveButton } from '../ui/form';
import { PageHeader } from '../layout/PageHeader';

const API_BASE_URL = 'http://localhost:9999';

export default function ConfigurationPage() {
  const [config, setConfig] = useState<ConfigValue>({});
  const [notification, setNotification] = useState<Notification | null>(null);

  useEffect(() => {
    fetch(`${API_BASE_URL}/config`)
      .then(res => res.json())
      .then(setConfig)
      .catch(error => {
        console.error('Error loading config:', error);
        setNotification({
          message: 'Failed to load configuration: ' + error.message,
          type: 'error'
        });
      });
  }, []);

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

      setNotification({
        message: 'Configuration saved successfully',
        type: 'success'
      });
    } catch (error) {
      setNotification({
        message: 'Failed to save configuration',
        type: 'error'
      });
    }
  };

  return (
    <Container>
      <PageHeader title="Configuration" />

      {notification && (
        <Alert 
          className={cn(
            "mb-6",
            notification.type === 'success' ? "bg-primary/20" : "bg-destructive/20"
          )}
        >
          {notification.message}
        </Alert>
      )}

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