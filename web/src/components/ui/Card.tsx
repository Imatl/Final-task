import { cn } from '@/lib/cn';
import type { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  variant?: 'default' | 'glass' | 'metric';
  hover?: boolean;
}

export function Card({ children, className, variant = 'default', hover = false }: CardProps) {
  return (
    <div
      className={cn(
        'rounded-xl border transition-all duration-200',
        variant === 'default' && 'bg-cosmic-800/90 border-cosmic-700/50 shadow-[0_4px_20px_rgba(0,0,0,0.4),0_0_40px_rgba(168,85,247,0.05)]',
        variant === 'glass' && 'bg-cosmic-800/60 border-cosmic-700/30 backdrop-blur-sm',
        variant === 'metric' && 'bg-cosmic-800/90 border-cosmic-700/50 shadow-[0_4px_20px_rgba(0,0,0,0.4)] relative overflow-hidden before:absolute before:top-0 before:left-0 before:right-0 before:h-[2px] before:bg-gradient-to-r before:from-velvet-600 before:via-neon-violet before:to-neon-pink',
        hover && 'hover:border-neon-violet/30 hover:-translate-y-0.5 hover:shadow-[0_8px_30px_rgba(0,0,0,0.5),0_0_50px_rgba(168,85,247,0.1)] cursor-pointer',
        className
      )}
    >
      {children}
    </div>
  );
}
