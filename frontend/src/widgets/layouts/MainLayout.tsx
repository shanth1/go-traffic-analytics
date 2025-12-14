import { Outlet, useLocation, Link } from 'react-router-dom';
import { LayersIcon, LinkIcon, LogOutIcon, UserIcon } from 'lucide-react';
import { cn } from '@/shared/lib/utils';
import { useAuthStore } from '@/entities/session/store';

const NAV_ITEMS = [
  { label: 'Профиль', path: '/', icon: UserIcon },
  { label: 'Кампании', path: '/campaigns', icon: LayersIcon },
  { label: 'Ссылки', path: '/links', icon: LinkIcon },
];

export const MainLayout = () => {
  const logout = useAuthStore((state) => state.logout);
  const location = useLocation();

  const isObjActive = (path: string) => location.pathname === path;

  return (
    <div className="flex min-h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-50">
      {/* --- Desktop Sidebar --- */}
      <aside className="hidden md:flex flex-col w-64 border-r border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 fixed inset-y-0 left-0 z-50">
        <div className="p-6">
          <h1 className="text-2xl font-bold tracking-tight text-indigo-600">
            Analytics
          </h1>
        </div>

        <nav className="flex-1 px-4 space-y-2">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium',
                isObjActive(item.path)
                  ? 'bg-indigo-50 text-indigo-600 dark:bg-indigo-900/20 dark:text-indigo-400'
                  : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'
              )}
            >
              <item.icon size={20} />
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="p-4 border-t border-slate-200 dark:border-slate-800">
          <button
            onClick={logout}
            className="flex items-center gap-3 px-4 py-3 w-full text-sm font-medium text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
          >
            <LogOutIcon size={20} />
            Выйти
          </button>
        </div>
      </aside>

      {/* --- Main Content Area --- */}
      <main className="flex-1 md:ml-64 pb-20 md:pb-0 relative">
        <div className="container mx-auto p-4 md:p-8 max-w-7xl">
          <Outlet />
        </div>
      </main>

      {/* --- Mobile Bottom Nav --- */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 z-50 pb-safe">
        <div className="flex justify-around items-center h-16">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                'flex flex-col items-center justify-center w-full h-full gap-1',
                isObjActive(item.path)
                  ? 'text-indigo-600 dark:text-indigo-400'
                  : 'text-slate-500 dark:text-slate-500'
              )}
            >
              <item.icon
                size={24}
                strokeWidth={isObjActive(item.path) ? 2.5 : 2}
              />
              <span className="text-[10px] font-medium">{item.label}</span>
            </Link>
          ))}
          <button
            onClick={logout}
            className="flex flex-col items-center justify-center w-full h-full gap-1 text-slate-500"
          >
            <LogOutIcon size={24} />
            <span className="text-[10px] font-medium">Выход</span>
          </button>
        </div>
      </nav>
    </div>
  );
};
