/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
  	container: {
  		center: true,
  		padding: '2rem',
  		screens: {
  			'2xl': '1400px'
  		}
  	},
  	extend: {
  		colors: {
  			border: 'hsl(var(--border))',
  			input: 'hsl(var(--input))',
  			ring: 'hsl(var(--ring))',
  			background: 'hsl(var(--background))',
  			foreground: 'hsl(var(--foreground))',
  			primary: {
  				DEFAULT: 'hsl(var(--primary))',
  				foreground: 'hsl(var(--primary-foreground))'
  			},
  			secondary: {
  				DEFAULT: 'hsl(var(--secondary))',
  				foreground: 'hsl(var(--secondary-foreground))'
  			},
  			destructive: {
  				DEFAULT: 'hsl(var(--destructive))',
  				foreground: 'hsl(var(--destructive-foreground))'
  			},
  			muted: {
  				DEFAULT: 'hsl(var(--muted))',
  				foreground: 'hsl(var(--muted-foreground))'
  			},
  			accent: {
  				DEFAULT: 'hsl(var(--accent))',
  				foreground: 'hsl(var(--accent-foreground))'
  			},
  			popover: {
  				DEFAULT: 'hsl(var(--popover))',
  				foreground: 'hsl(var(--popover-foreground))'
  			},
  			card: {
  				DEFAULT: 'hsl(var(--card))',
  				foreground: 'hsl(var(--card-foreground))'
  			},
  			chart: {
  				'1': 'hsl(var(--chart-1))',
  				'2': 'hsl(var(--chart-2))',
  				'3': 'hsl(var(--chart-3))',
  				'4': 'hsl(var(--chart-4))',
  				'5': 'hsl(var(--chart-5))'
  			}
  		},
  		borderRadius: {
  			lg: 'var(--radius)',
  			md: 'calc(var(--radius) - 2px)',
  			sm: 'calc(var(--radius) - 4px)'
  		}
  	}
  },
  plugins: [
  	require("tailwindcss-animate"),
  	require("daisyui")
  ],
  daisyui: {
  	themes: [
  		{
  			dark: {
  				"primary": "#a855f7",        // Purple 500
  				"secondary": "#7c3aed",      // Violet 600
  				"accent": "#c084fc",         // Purple 400
  				"neutral": "#18181b",        // Zinc 900
  				"base-100": "#09090b",       // Zinc 950 - Deep black
  				"base-200": "#18181b",       // Zinc 900
  				"base-300": "#27272a",       // Zinc 800
  				"base-content": "#fafafa",   // Zinc 50
  				"info": "#8b5cf6",           // Violet 500
  				"success": "#10b981",        // Emerald 500
  				"warning": "#f59e0b",        // Amber 500
  				"error": "#ef4444",          // Red 500
  			},
  		},
  	],
  	darkTheme: "dark",
  	base: true,
  	styled: true,
  	utils: true,
  },
}