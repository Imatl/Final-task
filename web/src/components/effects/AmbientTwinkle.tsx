import { useMemo } from 'react';

interface Star {
  id: number;
  x: number;
  y: number;
  size: number;
  color: string;
  duration: number;
  delay: number;
}

export function AmbientTwinkle({ starCount = 60 }: { starCount?: number }) {
  if (typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return null;
  }

  const colors = ['#a855f7', '#22d3ee', '#e879f9'];

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
  }, [starCount]);

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
