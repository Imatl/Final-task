import { useMemo } from 'react';
import { useThemeStore } from '@/store/theme';

interface Star {
  id: number;
  x: number;
  y: number;
  size: number;
  color: string;
  duration: number;
  delay: number;
}

const DARK_COLORS  = ['#a855f7', '#22d3ee', '#e879f9'];
const LIGHT_COLORS = ['#A67C42', '#C4847A', '#7A8F52', '#BF9A5A', '#C09878', '#8FAA70', '#D4A882'];

export function AmbientTwinkle({ starCount = 60 }: { starCount?: number }) {
  if (typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return null;
  }

  const { theme } = useThemeStore();
  const colors = theme === 'light' ? LIGHT_COLORS : DARK_COLORS;

  const stars = useMemo<Star[]>(() => {
    return Array.from({ length: starCount }, (_, i) => ({
      id: i,
      x: Math.random() * 100,
      y: Math.random() * 100,
      size: 1 + Math.random() * 2,
      color: colors[Math.floor(Math.random() * colors.length)],
      duration: 2000 + Math.random() * 4000,
      delay: Math.random() * 6000,
    }));
  }, [starCount, theme]);

  return (
    <div className="absolute inset-0 overflow-hidden">
      {stars.map((star) => (
        <div
          key={star.id}
          className="absolute rounded-full animate-twinkle"
          style={{
            left: `${star.x}%`,
            top: `${star.y}%`,
            width: `${star.size}px`,
            height: `${star.size}px`,
            backgroundColor: star.color,
            animationDuration: `${star.duration}ms`,
            animationDelay: `${star.delay}ms`,
          } as React.CSSProperties}
        />
      ))}
    </div>
  );
}
