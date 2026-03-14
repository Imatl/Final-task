import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { User, Lock, Eye, EyeOff } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/cn';
import { useAuthStore } from '@/store/auth';
import api from '@/api/client';

function GoogleIcon({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none">
      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
    </svg>
  );
}

export function AuthPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { login, loginWithGoogle } = useAuthStore();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleMessage = useCallback((event: MessageEvent) => {
    if (event.data?.type === 'google-auth-success') {
      console.log('[Google Auth] Received success from popup:', event.data.user?.email);
      loginWithGoogle(event.data.user);
      navigate('/dashboard');
    } else if (event.data?.type === 'google-auth-error') {
      console.error('[Google Auth] Error from popup:', event.data.error);
      setError(event.data.error || 'Google login failed');
    }
  }, [loginWithGoogle, navigate]);

  useEffect(() => {
    window.addEventListener('message', handleMessage);
    return () => window.removeEventListener('message', handleMessage);
  }, [handleMessage]);

  const handleGoogleLogin = async () => {
    setError('');
    try {
      console.log('[Google Auth] Requesting auth URL...');
      const { data } = await api.get<{ auth_url: string; state: string }>('/auth/google');
      console.log('[Google Auth] Opening popup...');

      const width = 500;
      const height = 600;
      const left = window.screenX + (window.outerWidth - width) / 2;
      const top = window.screenY + (window.outerHeight - height) / 2;

      window.open(
        data.auth_url,
        'Google Sign In',
        `width=${width},height=${height},left=${left},top=${top}`
      );
    } catch (err) {
      console.error('[Google Auth] Failed to get auth URL:', err);
      setError('Google login failed');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    const ok = await login(email, password);
    setLoading(false);
    if (ok) {
      navigate('/dashboard');
    } else {
      setError(t('auth.invalidCredentials'));
    }
  };

  const inputBase = cn(
    'w-full bg-cosmic-800/60 border border-neon-violet/30 text-white rounded-xl pl-11 pr-4 py-3 text-sm',
    'placeholder-gray-500 transition-all duration-200',
    'focus:outline-none focus:border-neon-violet/60 focus:ring-2 focus:ring-neon-violet/20'
  );

  return (
    <div className="min-h-screen bg-cosmic-900 flex items-center justify-center p-4 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-1/4 -left-1/4 w-1/2 h-1/2 bg-velvet-600/20 rounded-full blur-[120px]" />
        <div className="absolute bottom-1/4 -right-1/4 w-1/2 h-1/2 bg-neon-cyan/10 rounded-full blur-[120px]" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-neon-pink/5 rounded-full blur-[100px]" />
      </div>

      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage: `linear-gradient(rgba(168, 85, 247, 0.5) 1px, transparent 1px),
                           linear-gradient(90deg, rgba(168, 85, 247, 0.5) 1px, transparent 1px)`,
          backgroundSize: '50px 50px',
        }}
      />

      <div className="relative z-10 w-full max-w-[470px]">
        <h1 className="text-3xl font-bold text-white text-center mb-6">Kairon</h1>

        <div className="bg-cosmic-800/50 backdrop-blur-sm border border-neon-violet/50 rounded-2xl p-8 shadow-[0_0_13px_rgba(168,85,247,0.5),0_0_30px_rgba(168,85,247,0.3),inset_0_0_20px_rgba(168,85,247,0.1)]">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="px-4 py-3 bg-red-500/5 border border-red-500/20 rounded-xl text-sm text-red-300 text-center">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <label className="text-[16px] text-gray-300">{t('auth.email')}</label>
              <div className="relative">
                <User className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-neon-violet/60" />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  autoComplete="email"
                  className={inputBase}
                />
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-[16px] text-gray-300">{t('auth.password')}</label>
              <div className="relative">
                <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-neon-violet/60" />
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  autoComplete="current-password"
                  className={cn(inputBase, 'pr-11')}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  tabIndex={-1}
                  className="absolute right-3.5 top-1/2 -translate-y-1/2"
                >
                  {showPassword
                    ? <EyeOff className="w-5 h-5 text-neon-violet/60" />
                    : <Eye className="w-5 h-5 text-neon-violet/60" />
                  }
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className={cn(
                'w-full py-3 rounded-xl text-sm font-semibold transition-all duration-200',
                'bg-neon-violet/80 text-white',
                'hover:bg-neon-violet/70 active:scale-[0.98]',
                'disabled:opacity-60 disabled:cursor-not-allowed',
                loading && 'animate-pulse'
              )}
            >
              {loading ? '...' : t('auth.signIn')}
            </button>

            <div className="w-full border-t border-cosmic-700/50" />

            <button
              type="button"
              onClick={handleGoogleLogin}
              className={cn(
                'w-full relative overflow-hidden group flex items-center justify-center gap-3 py-3 rounded-xl text-sm font-medium transition-all duration-200',
                'bg-cosmic-800/80 border border-cosmic-700/50 text-white',
                'hover:border-cosmic-600/60 active:scale-[0.98]'
              )}
            >
              <GoogleIcon className="w-5 h-5" />
              <span className="relative z-10">{t('auth.continueGoogle')}</span>
              <div className="absolute inset-0 bg-gradient-to-r from-blue-600/10 via-red-600/10 to-yellow-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
