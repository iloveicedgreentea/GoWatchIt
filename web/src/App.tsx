// src/App.tsx
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { ThemeProvider } from "./components/ThemeProvider"
import { AppLayout } from "./components/layout/AppLayout"
import { Dashboard } from "./pages/Dashboard"
import ConfigurationPage from "./pages/ConfigurationPage"
import Logs from "./pages/Logs"
import { ToastContextProvider } from "./components/providers/toast";

function App() {
  return (
    <ThemeProvider defaultTheme="dark">
      <Router>
        <ToastContextProvider>
          <AppLayout>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/configuration" element={<ConfigurationPage />} />
              <Route path="/logs" element={<Logs />} />
            </Routes>
          </AppLayout>
        </ToastContextProvider>
      </Router>
    </ThemeProvider>
  )
}

export default App