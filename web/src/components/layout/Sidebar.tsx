import { NavLink, useNavigate } from 'react-router-dom';
import { useState, useRef, useEffect } from 'react';
import {
  LayoutDashboard,
  BarChart3,
  Headphones,
  User,
  ChevronLeft,
  ChevronRight,
  Plug,
  Moon,
  Sun,
  Users,
  LogOut,
  SlidersHorizontal,
  Building2,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { changeLanguage } from '@/i18n';
import { cn } from '@/lib/cn';
import { useUIStore } from '@/store/ui';
import { useThemeStore } from '@/store/theme';
import { useAuthStore, ROLE_LABELS } from '@/store/auth';

export function Sidebar() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { sidebarCollapsed, toggleSidebar } = useUIStore();
  const { theme, toggleTheme } = useThemeStore();
  const { user, logout } = useAuthStore();
  const isDark = theme === 'dark';
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement>(null);

  const level = user?.level ?? 1;

  const navItems = [
    { to: '/dashboard', label: t('nav.dashboard'), icon: LayoutDashboard, minLevel: 1 },
    { to: '/analytics', label: t('nav.analytics'), icon: BarChart3, minLevel: 2 },
    { to: '/staff', label: t('nav.staffPanel'), icon: Users, minLevel: 3 },
    { to: '/integrations', label: t('nav.integrations'), icon: Plug, minLevel: 4 },
    { to: '/settings', label: t('nav.settings'), icon: SlidersHorizontal, minLevel: 4 },
    { to: '/companies', label: t('nav.companies'), icon: Building2, minLevel: 5 },
  ].filter((item) => level >= item.minLevel);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 bg-cosmic-900/95 border-r border-cosmic-700/50',
        'flex flex-col h-screen transition-all duration-300 z-40',
        sidebarCollapsed ? 'w-16' : 'w-64'
      )}
    >
      <div className="flex items-center gap-3 px-4 py-5 border-b border-cosmic-700/50">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center shadow-[0_0_10px_rgba(168,85,247,0.2),0_0_20px_rgba(168,85,247,0.1)]">
          <Headphones className="w-5 h-5 text-white" />
        </div>
        {!sidebarCollapsed && (
          <div className="flex flex-col">
            <span className="font-bold text-white tracking-wide">SupportFlow</span>
            <span className="text-xs text-gray-400">{t('sidebar.subtitle')}</span>
          </div>
        )}
      </div>

      <nav className="flex-1 px-2 py-4 space-y-1 overflow-y-auto">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3.5 px-4 py-3.5 rounded-xl text-sm font-medium transition-all duration-200',
                'hover:bg-cosmic-700/50 hover:text-white',
                isActive
                  ? 'bg-velvet-700/40 text-white border-l-[3px] border-neon-violet shadow-[inset_0_0_20px_rgba(168,85,247,0.08)]'
                  : 'text-gray-400',
                sidebarCollapsed && 'justify-center px-0'
              )
            }
          >
            <item.icon className="w-[22px] h-[22px] flex-shrink-0" />
            {!sidebarCollapsed && <span className="truncate">{item.label}</span>}
          </NavLink>
        ))}
      </nav>

      <div className="px-2 py-4 border-t border-cosmic-700/50 relative" ref={userMenuRef}>
        <button
          onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
          className={cn(
            'w-full flex items-center gap-2 px-2 py-2 rounded-lg transition-all duration-200',
            'hover:bg-cosmic-700/50',
            isUserMenuOpen && 'bg-cosmic-700/50',
            sidebarCollapsed && 'justify-center px-2'
          )}
        >
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center flex-shrink-0">
            <User className="w-4 h-4 text-white" />
          </div>
          {!sidebarCollapsed && (
            <div className="flex flex-col items-start min-w-0">
              <p className="text-sm font-medium text-white text-left truncate w-full">
                {user?.name ?? 'Agent'}
              </p>
              <p className="text-xs text-gray-400 text-left truncate w-full">
                {user ? ROLE_LABELS[user.level] : ''}
              </p>
            </div>
          )}
        </button>

        {isUserMenuOpen && (
          <div className={cn(
            'absolute bottom-full mb-1 bg-cosmic-800 border border-cosmic-700 rounded-xl shadow-xl overflow-hidden z-50',
            sidebarCollapsed ? 'left-full ml-2 w-52' : 'left-3 right-3'
          )}>
            <div className="px-4 py-3 space-y-1">
              <p className="text-[10px] font-semibold uppercase tracking-wider text-gray-500 mb-2">
                {t('sidebar.languageNote')}
              </p>
              <div className="flex gap-1.5">
                <button
                  onClick={() => changeLanguage('en')}
                  className={cn(
                    'flex-1 py-1.5 rounded-lg text-xs font-bold tracking-wider transition-all duration-200',
                    i18n.language === 'en'
                      ? 'bg-neon-violet/20 border border-neon-violet/40 text-white'
                      : 'text-gray-500 hover:text-gray-300 hover:bg-cosmic-700/60 border border-transparent'
                  )}
                >
                  EN
                </button>
                <button
                  onClick={() => changeLanguage('uk')}
                  className={cn(
                    'flex-1 py-1.5 rounded-lg text-xs font-bold tracking-wider transition-all duration-200',
                    i18n.language === 'uk'
                      ? 'bg-neon-violet/20 border border-neon-violet/40 text-white'
                      : 'text-gray-500 hover:text-gray-300 hover:bg-cosmic-700/60 border border-transparent'
                  )}
                >
                  UA
                </button>
              </div>
            </div>

            <div className="px-4 py-3 border-t border-cosmic-700/50 space-y-1">
              <p className="text-[10px] font-semibold uppercase tracking-wider text-gray-500 mb-2">
                {t('sidebar.themeLabel')}
              </p>
              <div className="flex gap-1.5">
                <button
                  onClick={() => { if (!isDark) toggleTheme(); }}
                  className={cn(
                    'flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg text-xs font-medium transition-all duration-200',
                    isDark
                      ? 'bg-neon-violet/20 border border-neon-violet/40 text-white'
                      : 'text-gray-500 hover:text-gray-300 hover:bg-cosmic-700/60 border border-transparent'
                  )}
                >
                  <Moon className="w-3 h-3" />
                  {t('sidebar.dark')}
                </button>
                <button
                  onClick={() => { if (isDark) toggleTheme(); }}
                  className={cn(
                    'flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-lg text-xs font-medium transition-all duration-200',
                    !isDark
                      ? 'bg-neon-violet/20 border border-neon-violet/40 text-white'
                      : 'text-gray-500 hover:text-gray-300 hover:bg-cosmic-700/60 border border-transparent'
                  )}
                >
                  <Sun className="w-3 h-3" />
                  {t('sidebar.light')}
                </button>
              </div>
            </div>

            <div className="px-4 py-2 border-t border-cosmic-700/50">
              <button
                onClick={handleLogout}
                className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm text-gray-400 hover:bg-rose-500/10 hover:text-rose-400 transition-all duration-200"
              >
                <LogOut className="w-4 h-4" />
                {t('sidebar.logout')}
              </button>
            </div>
          </div>
        )}
      </div>

      <button
        onClick={toggleSidebar}
        className="absolute -right-3 top-20 w-6 h-6 bg-cosmic-700 border border-cosmic-600 rounded-full flex items-center justify-center text-gray-400 hover:text-white hover:bg-cosmic-600 transition-colors"
      >
        {sidebarCollapsed ? (
          <ChevronRight className="w-4 h-4" />
        ) : (
          <ChevronLeft className="w-4 h-4" />
        )}
      </button>
    </aside>
  );
}
