import { useEffect, useState, useMemo, useCallback } from 'react';

interface Particle {
  id: number;
  x: number;
  y: number;
  layer: number;
}

const layerConfigs = [
  { speed: 0.2, size: 1, opacity: 0.3, color: '#a855f7' },
  { speed: 0.8, size: 2, opacity: 0.4, color: '#c084fc' },
  { speed: 1.0, size: 3, opacity: 0.7, color: '#e879f9' },
];

export function FloatingParticleField() {
  const [scrollY, setScrollY] = useState(0);

  const prefersReducedMotion = typeof window !== 'undefined'
    ? window.matchMedia('(prefers-reduced-motion: reduce)').matches
    : false;

  const particles = useMemo<Particle[]>(() => {
    const result: Particle[] = [];
    for (let layer = 0; layer < 3; layer++) {
      for (let i = 0; i < 10; i++) {
        result.push({
          id: layer * 10 + i,
          x: Math.random() * 100,
          y: Math.random() * 100,
          layer,
        });
      }
    }
    return result;
  }, []);

  const handleScroll = useCallback(() => {
    if (prefersReducedMotion) return;
    setScrollY(window.scrollY);
  }, [prefersReducedMotion]);

  useEffect(() => {
    if (prefersReducedMotion) return;
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [handleScroll, prefersReducedMotion]);

  if (prefersReducedMotion) return null;

  return (
    <div className="absolute inset-0 overflow-hidden">
      {particles.map((particle) => {
        const config = layerConfigs[particle.layer];
        const translateY = scrollY * 0.5 * config.speed;
        const animDuration = 10 + Math.random() * 20;
        const animDelay = Math.random() * -20;

        return (
          <div
            key={particle.id}
            className="absolute rounded-full will-change-transform animate-float"
            style={{
              left: `${particle.x}%`,
              top: `${particle.y}%`,
              width: `${config.size}px`,
              height: `${config.size}px`,
              backgroundColor: config.color,
              opacity: config.opacity,
              transform: `translateY(${translateY}px)`,
              animationDuration: `${animDuration}s`,
              animationDelay: `${animDelay}s`,
            }}
          />
        );
      })}
    </div>
  );
}
