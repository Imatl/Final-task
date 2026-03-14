import { Outlet } from 'react-router-dom';
import { cn } from '@/lib/cn';
import { useUIStore } from '@/store/ui';
import { Sidebar } from './Sidebar';
import { AmbientTwinkle, FloatingParticleField } from '../effects';

export function MainLayout() {
  const { sidebarCollapsed } = useUIStore();

  return (
    <div className="min-h-screen bg-cosmic-950 relative overflow-hidden">
      <div className="fixed inset-0 pointer-events-none z-0">
        <FloatingParticleField />
        <AmbientTwinkle />
      </div>

      <div className="fixed inset-0 pointer-events-none z-0">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-neon-violet/5 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-neon-cyan/5 rounded-full blur-3xl" />
        <div className="absolute top-1/2 left-1/2 w-64 h-64 bg-neon-pink/5 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10">
        <Sidebar />
        <main
          className={cn(
            'min-h-screen transition-all duration-300',
            sidebarCollapsed ? 'pl-16' : 'pl-64'
          )}
        >
          <Outlet />
        </main>
      </div>
    </div>
  );
}
