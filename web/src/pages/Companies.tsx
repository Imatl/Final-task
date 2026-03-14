import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Building2, Users, TrendingUp, UserPlus, Copy, Check } from 'lucide-react';
import { companiesApi, inviteApi, type Company } from '@/api/client';
import { useAuthStore } from '@/store/auth';
import { cn } from '@/lib/cn';

function fmt(n: number) {
  return n.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function CompanyRow({ company }: { company: Company }) {
  return (
    <div className="flex items-center gap-4 px-5 py-4 hover:bg-cosmic-800/40 transition-colors">
      <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-velvet-700 to-neon-violet/60 flex items-center justify-center flex-shrink-0">
        <Building2 className="w-5 h-5 text-white" />
      </div>

      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-white truncate">{company.name}</p>
        <p className="text-xs text-gray-500 mt-0.5">
          {new Date(company.created_at).toLocaleDateString()}
        </p>
      </div>

      <div className="flex items-center gap-1.5 text-sm text-gray-300 w-24 justify-end flex-shrink-0">
        <Users className="w-3.5 h-3.5 text-neon-cyan flex-shrink-0" />
        <span>{company.staff_count}</span>
      </div>

      <div className="flex items-center gap-1 text-sm font-medium text-neon-green w-28 justify-end flex-shrink-0">
        <span>$</span>
        <span>{fmt(company.ai_spend_usd)}</span>
      </div>
    </div>
  );
}

function InvitePanel() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [generatedLink, setGeneratedLink] = useState('');
  const [copied, setCopied] = useState(false);

  const mutation = useMutation({
    mutationFn: () => inviteApi.generate(user!.id),
    onSuccess: (res) => {
      const fullLink = window.location.origin + res.data.link;
      setGeneratedLink(fullLink);
    },
  });

  const handleCopy = async () => {
    await navigator.clipboard.writeText(generatedLink);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const inputClass = cn(
    'w-full bg-cosmic-800/80 border border-cosmic-700/80 text-white rounded-xl px-4 py-2.5 text-sm',
    'placeholder-gray-600 transition-all duration-200',
    'focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/10'
  );

  return (
    <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-lg bg-neon-violet/10 border border-neon-violet/20 flex items-center justify-center">
            <UserPlus className="w-4.5 h-4.5 text-neon-violet" />
          </div>
          <h3 className="text-sm font-semibold text-white">{t('companies.generateInvite')}</h3>
        </div>
        <button
          onClick={() => mutation.mutate()}
          disabled={mutation.isPending}
          className={cn(
            'px-5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200',
            'bg-gradient-to-r from-velvet-600 to-neon-violet text-white',
            'hover:opacity-90 active:scale-[0.98]',
            'disabled:opacity-40 disabled:cursor-not-allowed',
            mutation.isPending && 'animate-pulse'
          )}
        >
          {t('companies.generateBtn')}
        </button>
      </div>

      {generatedLink && (
        <div className="space-y-1.5">
          <label className="block text-xs font-medium text-gray-400">
            {t('companies.inviteLinkLabel')}
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              value={generatedLink}
              readOnly
              className={cn(inputClass, 'flex-1 opacity-80 cursor-default')}
            />
            <button
              onClick={handleCopy}
              className={cn(
                'px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 flex items-center gap-1.5',
                copied
                  ? 'bg-neon-green/20 border border-neon-green/30 text-neon-green'
                  : 'bg-cosmic-800/80 border border-cosmic-700/80 text-gray-300 hover:text-white hover:border-cosmic-600/80'
              )}
            >
              {copied ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
              {copied ? t('companies.copied') : t('companies.copyLink')}
            </button>
          </div>
        </div>
      )}

      {mutation.isError && (
        <p className="text-sm text-red-400">{t('companies.error')}</p>
      )}
    </div>
  );
}

export function CompaniesPage() {
  const { t } = useTranslation();

  const { data, isLoading } = useQuery({
    queryKey: ['companies'],
    queryFn: () => companiesApi.list().then((r) => r.data),
    retry: false,
    throwOnError: false,
  });

  const companies: Company[] = data?.companies ?? [];
  const totalStaff = companies.reduce((s, c) => s + c.staff_count, 0);
  const totalSpend = companies.reduce((s, c) => s + c.ai_spend_usd, 0);

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">{t('companies.title')}</h1>
        <p className="text-sm text-gray-400 mt-1">{t('companies.subtitle')}</p>
      </div>

      <InvitePanel />

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-violet/10 border border-neon-violet/20 flex items-center justify-center mb-4">
            <Building2 className="w-6 h-6 text-neon-violet" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {isLoading ? '—' : companies.length}
          </p>
          <span className="text-sm text-gray-400">{t('companies.total')}</span>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-cyan/10 border border-neon-cyan/20 flex items-center justify-center mb-4">
            <Users className="w-6 h-6 text-neon-cyan" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {isLoading ? '—' : totalStaff.toLocaleString()}
          </p>
          <span className="text-sm text-gray-400">{t('companies.totalStaff')}</span>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-green/10 border border-neon-green/20 flex items-center justify-center mb-4">
            <TrendingUp className="w-6 h-6 text-neon-green" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {isLoading ? '—' : `$${fmt(totalSpend)}`}
          </p>
          <span className="text-sm text-gray-400">{t('companies.totalSpend')}</span>
        </div>
      </div>

      <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl overflow-hidden">
        <div className="hidden sm:flex items-center gap-4 px-5 py-3 border-b border-cosmic-700/50 text-[10px] font-semibold uppercase tracking-wider text-gray-500">
          <div className="w-10 flex-shrink-0" />
          <div className="flex-1">{t('companies.colName')}</div>
          <span className="w-24 text-right">{t('companies.colStaff')}</span>
          <span className="w-28 text-right">{t('companies.colAiSpend')}</span>
        </div>

        {isLoading ? (
          <div className="px-5 py-10 text-center text-sm text-gray-500">{t('dashboard.loading')}</div>
        ) : companies.length === 0 ? (
          <div className="px-5 py-10 text-center text-sm text-gray-500">{t('analytics.noData')}</div>
        ) : (
          <div className="divide-y divide-cosmic-700/30">
            {companies.map((c) => (
              <CompanyRow key={c.id} company={c} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
