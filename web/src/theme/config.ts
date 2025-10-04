// Theme configuration and constants
export type ThemeName = 'dark' | 'light' | 'sunset' | 'cyberpunk' | 'business' | 'night';

export interface ThemeConfig {
  name: ThemeName;
  displayName: string;
  description: string;
}

export const AVAILABLE_THEMES: ThemeConfig[] = [
  {
    name: 'dark',
    displayName: 'Modern Dark',
    description: 'Sleek dark theme with vibrant accents'
  },
  {
    name: 'night',
    displayName: 'Night',
    description: 'Deep dark theme for late-night coding'
  },
  {
    name: 'light',
    displayName: 'Professional Light',
    description: 'Clean and professional light theme'
  },
  {
    name: 'business',
    displayName: 'Business',
    description: 'Corporate-friendly neutral theme'
  },
  {
    name: 'sunset',
    displayName: 'Sunset',
    description: 'Warm and vibrant sunset colors'
  },
  {
    name: 'cyberpunk',
    displayName: 'Cyberpunk',
    description: 'Neon-inspired futuristic theme'
  }
];

export const DEFAULT_THEME: ThemeName = 'dark';
export const THEME_STORAGE_KEY = 'gowatchit-theme';
