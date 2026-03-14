import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Headphones, Eye, EyeOff } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/cn';
import { useAuthStore } from '@/store/auth';
import { AmbientTwinkle } from '@/components/effects';

export function AuthPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    await new Promise((r) => setTimeout(r, 600));
    const ok = login(email, password);
    setLoading(false);
    if (ok) {
      navigate('/dashboard');
    } else {
      setError(t('auth.invalidCredentials'));
    }
  };

  return (
    <div className="min-h-screen bg-cosmic-950 flex items-center justify-center relative overflow-hidden">
      <div className="fixed inset-0 pointer-events-none z-0">
        <AmbientTwinkle starCount={50} />
      </div>

      <div className="relative z-10 w-full max-w-md px-4">
        <div className="bg-cosmic-900/90 backdrop-blur-sm border border-cosmic-700/50 rounded-2xl shadow-xl overflow-hidden">
          <div className="px-8 pt-8 pb-6">
            <div className="text-center mb-7">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center mx-auto mb-3 shadow-[0_0_20px_rgba(168,85,247,0.3)]">
                <Headphones className="w-6 h-6 text-white" />
              </div>
              <h1 className="text-2xl font-bold text-white mb-1">{t('auth.welcomeBack')}</h1>
              <p className="text-sm text-gray-500">{t('auth.tagline')}</p>
            </div>

            {error && (
              <div className="mb-4 px-4 py-2.5 bg-red-500/10 border border-red-500/30 rounded-xl text-sm text-red-400">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.email')}
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  autoComplete="email"
                  placeholder="you@company.com"
                  className={cn(
                    'w-full bg-cosmic-800/80 border border-cosmic-700/80 text-white rounded-xl px-4 py-2.5 text-sm',
                    'placeholder-gray-600 transition-all duration-200',
                    'focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/10'
                  )}
                />
              </div>

              <div>
                <div className="flex items-center justify-between mb-1.5">
                  <label className="block text-xs font-medium text-gray-400">
                    {t('auth.password')}
                  </label>
                  <button type="button" className="text-xs text-neon-violet hover:text-neon-purple transition-colors">
                    {t('auth.forgotPassword')}
                  </button>
                </div>
                <div className="relative">
                  <input
                    type={showPass ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    autoComplete="current-password"
                    placeholder="••••••••"
                    className={cn(
                      'w-full bg-cosmic-800/80 border border-cosmic-700/80 text-white rounded-xl px-4 py-2.5 text-sm pr-10',
                      'placeholder-gray-600 transition-all duration-200',
                      'focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/10'
                    )}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPass(!showPass)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
                  >
                    {showPass ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className={cn(
                  'w-full py-2.5 rounded-xl text-sm font-semibold transition-all duration-200',
                  'bg-gradient-to-r from-velvet-600 to-neon-violet text-white',
                  'hover:opacity-90 active:scale-[0.98]',
                  'shadow-[0_0_20px_rgba(168,85,247,0.25)] hover:shadow-[0_0_30px_rgba(168,85,247,0.4)]',
                  'disabled:opacity-60 disabled:cursor-not-allowed',
                  loading && 'animate-pulse'
                )}
              >
                {loading ? '...' : t('auth.signIn')}
              </button>
            </form>
          </div>

          <div className="px-8 py-4 bg-cosmic-950/40 border-t border-cosmic-700/30 text-center">
            <span className="text-sm text-gray-500">
              {t('auth.noAccount')}{' '}
              <Link
                to="/register"
                className="text-neon-violet hover:text-neon-purple font-medium transition-colors"
              >
                {t('auth.toSignUp')}
              </Link>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
