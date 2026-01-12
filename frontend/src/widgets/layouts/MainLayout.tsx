import { Outlet, useLocation, Link } from 'react-router-dom';
import {
  LayoutDashboardIcon, // New icon for dashboard
  LayersIcon,
  LinkIcon,
  LogOutIcon,
  UserIcon,
} from 'lucide-react';
import { cn } from '@/shared/lib/utils';
import { useAuthStore } from '@/entities/session/store';

const NAV_ITEMS = [
  { label: 'Dashboard', path: '/', icon: LayoutDashboardIcon },
  { label: 'Campaigns', path: '/campaigns', icon: LayersIcon },
  { label: 'Links', path: '/links', icon: LinkIcon },
  { label: 'Profile', path: '/profile', icon: UserIcon },
];

export const MainLayout = () => {
  const logout = useAuthStore((state) => state.logout);
  const location = useLocation();

  const isObjActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  return (
    <div className="flex min-h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-50 overflow-hidden">
      {/* --- Desktop Sidebar --- */}
      <aside className="hidden md:flex flex-col w-64 border-r border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 fixed inset-y-0 left-0 z-40 shrink-0">
        <div className="p-6">
          <h1 className="text-2xl font-bold tracking-tight text-indigo-600 flex items-center gap-2">
            <LayoutDashboardIcon className="text-indigo-600" />
            Analytics
          </h1>
        </div>

        <nav className="flex-1 px-4 space-y-2 overflow-y-auto">
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
          <div className="mb-4 px-4 flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-slate-200 dark:bg-slate-800 flex items-center justify-center">
              <UserIcon size={16} />
            </div>
            <div className="flex flex-col">
              <span className="text-xs font-bold">My Account</span>
              <Link
                to="/profile"
                className="text-[10px] text-slate-500 hover:text-indigo-500"
              >
                View Profile
              </Link>
            </div>
          </div>
          <button
            onClick={logout}
            className="flex items-center gap-3 px-4 py-2 w-full text-sm font-medium text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
          >
            <LogOutIcon size={18} />
            Logout
          </button>
        </div>
      </aside>

      {/* --- Main Content Area --- */}

      <main className="flex-1 min-w-0 md:ml-64 pb-24 md:pb-8 relative overflow-x-hidden">
        <div className="container mx-auto p-4 md:p-8 max-w-7xl">
          <Outlet />
        </div>
      </main>

      {/* --- Mobile Bottom Nav --- */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white/90 backdrop-blur-md dark:bg-slate-950/90 border-t border-slate-200 dark:border-slate-800 z-50 pb-safe shadow-lg">
        <div className="flex justify-around items-center h-16">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                'flex flex-col items-center justify-center w-full h-full gap-1 active:scale-95 transition-transform',
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
        </div>
      </nav>
    </div>
  );
};
