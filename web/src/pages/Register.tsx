import { useState, useEffect, useCallback } from 'react';
import { Link, useParams, Navigate } from 'react-router-dom';
import { Eye, EyeOff, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/cn';
import api from '@/api/client';
import { inviteApi } from '@/api/client';
import { AxiosError } from 'axios';

export function RegisterPage() {
  const { t } = useTranslation();
  const { token } = useParams<{ token: string }>();

  const [fullName, setFullName] = useState('');
  const [company, setCompany] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [googleLinked, setGoogleLinked] = useState(false);

  const [validating, setValidating] = useState(true);
  const [tokenValid, setTokenValid] = useState(false);

  useEffect(() => {
    if (!token) return;
    inviteApi.validate(token)
      .then((res) => {
        if (res.data.valid) {
          setTokenValid(true);
        }
      })
      .catch(() => {
        setTokenValid(false);
      })
      .finally(() => setValidating(false));
  }, [token]);

  const handleMessage = useCallback((event: MessageEvent) => {
    if (event.data?.type === 'google-auth-success') {
      const user = event.data.user;
      setFullName(user?.name || '');
      setEmail(user?.email || '');
      setGoogleLinked(true);
      setError('');
    } else if (event.data?.type === 'google-auth-error') {
      setError(t('auth.googleFailed'));
    }
  }, [t]);

  useEffect(() => {
    window.addEventListener('message', handleMessage);
    return () => window.removeEventListener('message', handleMessage);
  }, [handleMessage]);

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  if (validating) {
    return (
      <div className="min-h-screen bg-cosmic-900 flex items-center justify-center">
        <Loader2 className="w-8 h-8 text-neon-violet animate-spin" />
      </div>
    );
  }

  if (!tokenValid) {
    return (
      <div className="min-h-screen bg-cosmic-900 flex items-center justify-center p-4 relative overflow-hidden">
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-1/4 -left-1/4 w-1/2 h-1/2 bg-velvet-600/20 rounded-full blur-[120px]" />
          <div className="absolute bottom-1/4 -right-1/4 w-1/2 h-1/2 bg-neon-cyan/10 rounded-full blur-[120px]" />
        </div>
        <div className="relative z-10 w-full max-w-md px-4">
          <div className="bg-cosmic-900/90 backdrop-blur-sm border border-cosmic-700/50 rounded-2xl shadow-xl p-8 text-center">
            <div className="w-12 h-12 rounded-full bg-red-500/20 border border-red-500/30 flex items-center justify-center mx-auto mb-4">
              <span className="text-red-400 text-xl">&#10007;</span>
            </div>
            <h2 className="text-xl font-bold text-white mb-2">{t('auth.invalidToken')}</h2>
            <Link
              to="/login"
              className="inline-block mt-4 px-6 py-2.5 rounded-xl text-sm font-semibold bg-gradient-to-r from-velvet-600 to-neon-violet text-white hover:opacity-90 transition-opacity"
            >
              {t('auth.toSignIn')}
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const handleGoogleLink = async () => {
    setError('');
    try {
      const { data } = await api.get<{ auth_url: string; state: string }>('/auth/google');

      const width = 500;
      const height = 600;
      const left = window.screenX + (window.outerWidth - width) / 2;
      const top = window.screenY + (window.outerHeight - height) / 2;

      window.open(
        data.auth_url,
        'Google Sign In',
        `width=${width},height=${height},left=${left},top=${top}`
      );
    } catch {
      setError(t('auth.googleFailed'));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirm) {
      setError(t('auth.passwordMismatch'));
      return;
    }
    setLoading(true);
    setError('');
    try {
      await inviteApi.register({ token: token!, name: fullName, company, email, password });
      setSuccess(true);
    } catch (err) {
      const axiosErr = err as AxiosError<{ error: string }>;
      if (axiosErr.response?.status === 409) {
        setError(t('auth.emailTaken'));
      } else {
        setError(t('auth.registrationFailed'));
      }
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="min-h-screen bg-cosmic-900 flex items-center justify-center p-4 relative overflow-hidden">
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-1/4 -left-1/4 w-1/2 h-1/2 bg-velvet-600/20 rounded-full blur-[120px]" />
          <div className="absolute bottom-1/4 -right-1/4 w-1/2 h-1/2 bg-neon-cyan/10 rounded-full blur-[120px]" />
        </div>
        <div className="relative z-10 w-full max-w-md px-4">
          <div className="bg-cosmic-900/90 backdrop-blur-sm border border-cosmic-700/50 rounded-2xl shadow-xl p-8 text-center">
            <div className="w-12 h-12 rounded-full bg-neon-green/20 border border-neon-green/30 flex items-center justify-center mx-auto mb-4">
              <span className="text-neon-green text-xl">&#10003;</span>
            </div>
            <h2 className="text-xl font-bold text-white mb-2">{t('auth.registrationSent')}</h2>
            <p className="text-sm text-gray-400 mb-6">{t('auth.registrationPending')}</p>
            <Link
              to="/login"
              className="inline-block px-6 py-2.5 rounded-xl text-sm font-semibold bg-gradient-to-r from-velvet-600 to-neon-violet text-white hover:opacity-90 transition-opacity"
            >
              {t('auth.toSignIn')}
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const inputClass = cn(
    'w-full bg-cosmic-800/80 border border-cosmic-700/80 text-white rounded-xl px-4 py-2.5 text-sm',
    'placeholder-gray-600 transition-all duration-200',
    'focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/10'
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

      <div className="relative z-10 w-full max-w-md px-4">
        <div className="bg-cosmic-800/50 backdrop-blur-sm border border-neon-violet/50 rounded-2xl shadow-[0_0_13px_rgba(168,85,247,0.5),0_0_30px_rgba(168,85,247,0.3),inset_0_0_20px_rgba(168,85,247,0.1)] overflow-hidden">
          <div className="px-8 pt-8 pb-6">
            <div className="text-center mb-7">
              <h1 className="text-2xl font-bold text-white mb-1">Kairon</h1>
              <p className="text-sm text-gray-500">{t('auth.tagline')}</p>
            </div>

            {error && (
              <div className="mb-4 px-4 py-2.5 bg-red-500/10 border border-red-500/30 rounded-xl text-sm text-red-400">
                {error}
              </div>
            )}

            {!googleLinked && (
              <>
                <button
                  type="button"
                  onClick={handleGoogleLink}
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
                  {t('auth.signUpGoogle')}
                </button>

                <div className="flex items-center gap-3 mt-4">
                  <div className="flex-1 h-px bg-cosmic-700/50" />
                  <span className="text-xs text-gray-600">{t('auth.orManual')}</span>
                  <div className="flex-1 h-px bg-cosmic-700/50" />
                </div>
              </>
            )}

            {googleLinked && (
              <div className="mb-4 px-4 py-2.5 bg-neon-green/10 border border-neon-green/30 rounded-xl text-sm text-neon-green flex items-center gap-2">
                <svg className="w-4 h-4 shrink-0" viewBox="0 0 24 24">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                {t('auth.googleLinked')}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4 mt-4">
              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.fullName')}
                </label>
                <input
                  type="text"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  required
                  autoComplete="name"
                  placeholder={t('auth.fullNamePlaceholder')}
                  className={inputClass}
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.company')}
                </label>
                <input
                  type="text"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                  required
                  autoComplete="organization"
                  placeholder={t('auth.companyPlaceholder')}
                  className={inputClass}
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.email')}
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  readOnly={googleLinked}
                  autoComplete="email"
                  placeholder="you@company.com"
                  className={cn(inputClass, googleLinked && 'opacity-60 cursor-not-allowed')}
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.password')}
                </label>
                <div className="relative">
                  <input
                    type={showPass ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    autoComplete="new-password"
                    placeholder="••••••••"
                    className={cn(inputClass, 'pr-10')}
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

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1.5">
                  {t('auth.confirmPassword')}
                </label>
                <div className="relative">
                  <input
                    type={showConfirm ? 'text' : 'password'}
                    value={confirm}
                    onChange={(e) => setConfirm(e.target.value)}
                    required
                    autoComplete="new-password"
                    placeholder="••••••••"
                    className={cn(inputClass, 'pr-10')}
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirm(!showConfirm)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 transition-colors"
                  >
                    {showConfirm ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
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
                {loading ? '...' : t('auth.signUp')}
              </button>
            </form>
          </div>

          <div className="px-8 py-4 bg-cosmic-950/40 border-t border-cosmic-700/30 text-center">
            <span className="text-sm text-gray-500">
              {t('auth.hasAccount')}{' '}
              <Link
                to="/login"
                className="text-neon-violet hover:text-neon-purple font-medium transition-colors"
              >
                {t('auth.toSignIn')}
              </Link>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
