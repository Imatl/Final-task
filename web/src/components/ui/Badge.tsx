import { cn } from '@/lib/cn';

const VARIANTS = {
  default: 'bg-cosmic-700 text-gray-300',
  success: 'bg-green-500/20 text-green-400',
  warning: 'bg-amber-500/20 text-amber-400',
  error: 'bg-red-500/20 text-red-400',
  info: 'bg-cyan-500/20 text-cyan-400',
  neon: 'bg-neon-violet/20 text-neon-purple',
};

interface BadgeProps {
  children: React.ReactNode;
  variant?: keyof typeof VARIANTS;
  dot?: boolean;
  className?: string;
}

export function Badge({ children, variant = 'default', dot = false, className }: BadgeProps) {
  return (
    <span className={cn('inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium', VARIANTS[variant], className)}>
      {dot && <span className={cn('w-1.5 h-1.5 rounded-full', variant === 'success' ? 'bg-green-400' : variant === 'error' ? 'bg-red-400' : variant === 'warning' ? 'bg-amber-400' : 'bg-neon-violet')} />}
      {children}
    </span>
  );
}
