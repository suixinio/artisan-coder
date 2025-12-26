import { StrictMode, useEffect } from 'react'
import { createRoot } from 'react-dom/client'
import '@/styles/index.css'
import { Router } from '@/router'
import { ToastProvider } from '@/components/layout/ToastProvider'
import { useUserStore } from '@/stores/userStore'
import { useThemeStore } from '@/stores/themeStore'

// Initialize theme immediately to prevent flash
const themeStore = useThemeStore.getState()
const theme = themeStore.theme
const root = window.document.documentElement

if (theme === 'dark') {
  root.classList.add('dark')
} else if (theme === 'light') {
  root.classList.add('light')
} else {
  // System theme
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    root.classList.add('dark')
  } else {
    root.classList.add('light')
  }
}

function App() {
  const initialize = useUserStore((state) => state.initialize)

  useEffect(() => {
    initialize()
  }, [initialize])

  return <Router />
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ToastProvider />
    <App />
  </StrictMode>,
)
