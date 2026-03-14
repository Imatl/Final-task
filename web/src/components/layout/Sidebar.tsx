import { NavLink } from 'react-router-dom';
import {
  MessageSquare,
  LayoutDashboard,
  BarChart3,
  Settings,
  Headphones,
} from 'lucide-react';
import { cn } from '@/lib/cn';

const navItems = [
  { to: '/chat', label: 'Customer Chat', icon: MessageSquare },
  { to: '/dashboard', label: 'Agent Dashboard', icon: LayoutDashboard },
  { to: '/analytics', label: 'Analytics', icon: BarChart3 },
  { to: '/settings', label: 'Settings', icon: Settings },
];

export function Sidebar() {
  return (
    <aside className="w-64 bg-cosmic-900 border-r border-cosmic-700/50 flex flex-col h-screen fixed left-0 top-0 z-40">
      <div className="px-5 py-5 border-b border-cosmic-700/50 flex items-center gap-3">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center neon-glow-sm">
          <Headphones className="w-5 h-5 text-white" />
        </div>
        <div>
          <span className="text-base font-bold text-white tracking-wide">SupportFlow</span>
          <p className="text-xs text-gray-500">AI Assistant</p>
        </div>
      </div>

      <nav className="flex-1 px-3 py-4 space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm transition-all duration-200',
                isActive
                  ? 'bg-velvet-900/50 text-neon-purple border-l-2 border-neon-violet'
                  : 'text-gray-400 hover:text-gray-200 hover:bg-cosmic-800/50'
              )
            }
          >
            <item.icon className="w-5 h-5" />
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="px-4 py-4 border-t border-cosmic-700/50">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-neon-violet to-neon-pink flex items-center justify-center text-xs font-bold text-white">
            A
          </div>
          <div>
            <p className="text-sm text-gray-300">Agent</p>
            <p className="text-xs text-gray-500">Online</p>
          </div>
          <div className="ml-auto w-2 h-2 rounded-full bg-neon-green" />
        </div>
      </div>
    </aside>
  );
}
