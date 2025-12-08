import { Outlet, Navigate, useLocation } from 'react-router-dom';
import { BottomMenu } from '@/widgets/bottom-menu/BottomMenu';

export const Layout = () => {
  const location = useLocation();
  if (location.pathname === '/') return <Navigate to="/campaigns" replace />;

  return (
    <div className="min-h-screen bg-background text-white pb-24 font-sans selection:bg-primary/30">
      <div className="max-w-md mx-auto md:max-w-lg min-h-screen relative">
        <header className="p-4 flex items-center justify-between sticky top-0 z-40 bg-background/80 backdrop-blur-sm">
          <div className="font-bold text-xl tracking-tighter bg-linear-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            Linkly.
          </div>
          <div className="w-8 h-8 rounded-full bg-slate-800 border border-slate-700" />
        </header>

        <main className="px-4">
          <Outlet />
        </main>
      </div>
      <BottomMenu />
    </div>
  );
};
