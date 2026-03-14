import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Headphones, Eye, EyeOff } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useGoogleLogin } from '@react-oauth/google';
import { cn } from '@/lib/cn';
import { useAuthStore } from '@/store/auth';
import { AmbientTwinkle } from '@/components/effects';

const GOOGLE_CLIENT_ID: string = import.meta.env.VITE_GOOGLE_CLIENT_ID ?? '';

function GoogleButton({ onError }: { onError: (msg: string) => void }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { loginWithGoogle } = useAuthStore();

  const googleLogin = useGoogleLogin({
    onSuccess: async (tokenResponse) => {
      try {
        const res = await fetch('https://www.googleapis.com/oauth2/v3/userinfo', {
          headers: { Authorization: `Bearer ${tokenResponse.access_token}` },
        });
        const profile = await res.json();
        loginWithGoogle({ email: profile.email, name: profile.name, sub: profile.sub });
        navigate('/dashboard');
      } catch {
        onError('Google login failed');
      }
    },
    onError: () => onError('Google login failed'),
  });

  return (
    <>
      <button
        type="button"
        onClick={() => googleLogin()}
        className={cn(
          'w-full flex items-center justify-center gap-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200',
          'bg-cosmic-800/80 border border-cosmic-700/80 text-white',
          'hover:bg-cosmic-700/80 hover:border-cosmic-600/80 active:scale-[0.98]'
        )}
      >
        <svg className="w-4 h-4" viewBox="0 0 24 24">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"/>
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
        </svg>
        {t('auth.continueGoogle')}
      </button>

      <div className="flex items-center gap-3 mt-4">
        <div className="flex-1 h-px bg-cosmic-700/50" />
        <span className="text-xs text-gray-600">{t('auth.orEmail')}</span>
        <div className="flex-1 h-px bg-cosmic-700/50" />
      </div>
    </>
  );
}

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

            {GOOGLE_CLIENT_ID && <GoogleButton onError={setError} />}

            <form onSubmit={handleSubmit} className="space-y-4 mt-4">
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

        </div>
      </div>
    </div>
  );
}
