// src/components/ThemeProvider.tsx
import { createContext, useState } from "react"

type Theme = "dark" | "light" | "system"

interface ThemeProviderProps {
    children: React.ReactNode
    defaultTheme?: Theme
    storageKey?: string
}

const ThemeProviderContext = createContext<{
    theme: Theme
    setTheme: (theme: Theme) => void
}>({
    theme: "dark",
    setTheme: () => null,
})

// Apply theme immediately before React hydration
if (typeof window !== 'undefined') {
    document.documentElement.classList.add('dark')
}

export function ThemeProvider({
    children,
    defaultTheme = "dark",
    storageKey = "ui-theme",
    ...props
}: ThemeProviderProps) {
    const [theme, setTheme] = useState<Theme>(defaultTheme)

    // useEffect no longer needed since we set it immediately above
    const value = {
        theme,
        setTheme: (theme: Theme) => setTheme(theme),
    }

    return (
        <ThemeProviderContext.Provider {...props} value={value}>
            {children}
        </ThemeProviderContext.Provider>
    )
}