import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Building2, Users, TrendingUp } from 'lucide-react';
import { companiesApi, type Company } from '@/api/client';

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
