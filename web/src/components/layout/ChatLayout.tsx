import { Outlet, Link } from 'react-router-dom';
import { Headphones, Shield } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { changeLanguage } from '@/i18n';
import { cn } from '@/lib/cn';
import { AmbientTwinkle, FloatingParticleField } from '../effects';

function LanguageToggle() {
  const { i18n } = useTranslation();

  return (
    <div className="relative flex p-1 bg-cosmic-900/50 rounded-lg border border-cosmic-700">
      <div
        className={cn(
          'absolute inset-y-1 w-[calc(50%-4px)] bg-neon-violet/20 border border-neon-violet/30 rounded-md transition-all duration-300 ease-in-out',
          i18n.language === 'uk' ? 'left-[calc(48%+1px)]' : 'left-1'
        )}
      />
      <button
        onClick={() => changeLanguage('en')}
        className={cn(
          'relative z-10 flex-1 px-2 py-1 text-[10px] font-bold tracking-wider transition-colors duration-200',
          i18n.language === 'en' ? 'text-white' : 'text-gray-500 hover:text-gray-400'
        )}
      >
        EN
      </button>
      <button
        onClick={() => changeLanguage('uk')}
        className={cn(
          'relative z-10 flex-1 px-2 py-1 text-[10px] font-bold tracking-wider transition-colors duration-200',
          i18n.language === 'uk' ? 'text-white' : 'text-gray-500 hover:text-gray-400'
        )}
      >
        UA
      </button>
    </div>
  );
}

export function ChatLayout() {
  return (
    <div className="min-h-screen bg-cosmic-950 relative overflow-hidden">
      <div className="fixed inset-0 pointer-events-none z-0">
        <FloatingParticleField />
        <AmbientTwinkle starCount={40} />
      </div>

      <div className="fixed inset-0 pointer-events-none z-0">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-neon-violet/5 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-neon-cyan/5 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 flex flex-col h-screen">
        <header className="bg-cosmic-900/80 backdrop-blur-sm border-b border-cosmic-700/50 px-6 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center neon-glow-sm">
              <Headphones className="w-5 h-5 text-white" />
            </div>
            <div>
              <span className="font-bold text-white tracking-wide">SupportFlow</span>
              <p className="text-xs text-gray-500">AI Support</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <LanguageToggle />
            <Link
              to="/agent/dashboard"
              className="flex items-center gap-1.5 px-3 py-1.5 text-xs text-gray-400 hover:text-neon-purple bg-cosmic-800/80 hover:bg-cosmic-700/80 border border-cosmic-700/50 rounded-lg transition-all duration-200"
            >
              <Shield className="w-3.5 h-3.5" />
              Agent
            </Link>
          </div>
        </header>

        <main className="flex-1 overflow-hidden">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
