import { NavLink, useNavigate } from 'react-router-dom';
import { useState, useRef, useEffect } from 'react';
import {
  MessageSquare,
  LayoutDashboard,
  BarChart3,
  Headphones,
  Settings,
  Globe,
  User,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { changeLanguage } from '@/i18n';
import { cn } from '@/lib/cn';
import { useUIStore } from '@/store/ui';

export function Sidebar() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { sidebarCollapsed, toggleSidebar } = useUIStore();
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement>(null);

  const navItems = [
    { to: '/chat', label: t('nav.chat'), icon: MessageSquare },
    { to: '/dashboard', label: t('nav.dashboard'), icon: LayoutDashboard },
    { to: '/analytics', label: t('nav.analytics'), icon: BarChart3 },
  ];

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

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

      <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        {!sidebarCollapsed && (
          <div className="pb-2">
            <span className="px-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">
              {t('sidebar.main')}
            </span>
          </div>
        )}
        {sidebarCollapsed && <div className="border-t border-cosmic-700/50 mx-2 mb-2" />}

        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm transition-all duration-200',
                'hover:bg-cosmic-700/50 hover:text-gray-200',
                isActive && 'bg-velvet-900/50 text-neon-purple border-l-2 border-neon-violet',
                !isActive && 'text-gray-400',
                sidebarCollapsed && 'justify-center px-2'
              )
            }
          >
            <item.icon className="w-5 h-5 flex-shrink-0" />
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
              <p className="text-sm font-medium text-white text-left truncate w-full">Demo User</p>
              <p className="text-xs text-gray-400 text-left truncate w-full">support agent</p>
            </div>
          )}
        </button>

        {isUserMenuOpen && (
          <div className={cn(
            'absolute bottom-full mb-1 bg-cosmic-800 border border-cosmic-700 rounded-lg shadow-xl overflow-hidden z-50',
            sidebarCollapsed ? 'left-full ml-2 w-48' : 'left-3 right-3'
          )}>
            <button
              onClick={() => {
                navigate('/settings');
                setIsUserMenuOpen(false);
              }}
              className="w-full flex items-center gap-3 px-4 py-3 text-gray-300 hover:bg-cosmic-700 hover:text-white transition-colors"
            >
              <Settings className="w-4 h-4" />
              <span>{t('sidebar.settings')}</span>
            </button>

            <div className="px-4 py-3 border-t border-cosmic-700">
              <div className="flex items-center gap-2">
                <div className="flex items-center text-gray-400">
                  <Globe className="w-4 h-4" />
                </div>
                <div className="relative flex flex-1 max-w-[120px] p-1 bg-cosmic-900/50 rounded-lg border border-cosmic-700">
                  <div
                    className={cn(
                      'absolute inset-y-1 w-[calc(50%-4px)] bg-neon-violet/20 border border-neon-violet/30 rounded-md transition-all duration-300 ease-in-out',
                      i18n.language === 'uk' ? 'left-[calc(48%+1px)]' : 'left-1'
                    )}
                  />
                  <button
                    onClick={() => changeLanguage('en')}
                    className={cn(
                      'relative z-10 flex-1 py-1 text-[10px] font-bold tracking-wider transition-colors duration-200',
                      i18n.language === 'en' ? 'text-white' : 'text-gray-500 hover:text-gray-400'
                    )}
                  >
                    EN
                  </button>
                  <button
                    onClick={() => changeLanguage('uk')}
                    className={cn(
                      'relative z-10 flex-1 py-1 text-[10px] font-bold tracking-wider transition-colors duration-200',
                      i18n.language === 'uk' ? 'text-white' : 'text-gray-500 hover:text-gray-400'
                    )}
                  >
                    UA
                  </button>
                </div>
              </div>
              <p className="text-xs text-gray-500 mt-1">{t('sidebar.languageNote')}</p>
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
