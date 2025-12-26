import { NavLink } from 'react-router-dom'
import { LayoutDashboard, FolderKanban, Settings } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useUserStore } from '@/stores/userStore'

const navItems = [
  { title: 'Dashboard', href: '/', icon: LayoutDashboard },
  { title: 'Projects', href: '/projects', icon: FolderKanban },
  { title: 'Settings', href: '/settings', icon: Settings },
]

export function Sidebar() {
  const { sidebarOpen, toggleSidebar } = useUserStore()

  const handleNavClick = () => {
    // Close sidebar on mobile after navigation
    if (window.innerWidth < 768) {
      toggleSidebar()
    }
  }

  return (
    <aside
      className={cn(
        'fixed inset-y-0 left-0 z-50 w-64 flex-col border-r bg-background transition-transform duration-200 ease-in-out flex',
        'md:relative md:translate-x-0',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      )}
    >
      {/* Logo */}
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-lg font-bold">Artisan Coder</h1>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 p-4">
        {navItems.map((item) => (
          <NavLink
            key={item.href}
            to={item.href}
            onClick={handleNavClick}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              )
            }
          >
            <item.icon className="h-4 w-4" />
            {item.title}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
