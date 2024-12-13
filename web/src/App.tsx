// src/App.tsx
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { ThemeProvider } from "./components/ThemeProvider"
import { AppLayout } from "./components/layout/AppLayout"
import { Dashboard } from "./components/dashboard/Dashboard"
import  ConfigurationPage  from "./components/config/ConfigurationPage"

function App() {
  return (
    <ThemeProvider defaultTheme="dark">
      <Router>
        <AppLayout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/configuration" element={<ConfigurationPage />} />
          </Routes>
        </AppLayout>
      </Router>
    </ThemeProvider>
  )
}

export default App