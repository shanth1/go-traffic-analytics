import { Outlet, useLocation, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  LayoutDashboardIcon,
  LayersIcon,
  LinkIcon,
  UserIcon,
  LogOutIcon,
  CreditCardIcon,
  FileDownIcon,
  SparklesIcon,
} from 'lucide-react';
import { cn } from '@/shared/lib/utils';
import { useAuthStore } from '@/entities/session/store';
import { useTranslation } from 'react-i18next';
import { APP_TITLE } from '@/shared/config';

const MAIN_NAV = [
  { label: 'Dashboard', path: '/', icon: LayoutDashboardIcon },
  { label: 'Campaigns', path: '/campaigns', icon: LayersIcon },
  { label: 'Links', path: '/links', icon: LinkIcon },
];

const SECONDARY_NAV = [
  { label: 'Export Data', path: '/export', icon: FileDownIcon },
  { label: 'Plans & Billing', path: '/pricing', icon: CreditCardIcon },
];

export const MainLayout = () => {
  const { user, logout } = useAuthStore();
  const location = useLocation();
  const { t } = useTranslation();

  const isPathActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  if (!user) return null;

  const avatarUrl = `https://api.dicebear.com/9.x/thumbs/svg?seed=Brian`;

  return (
    <div className="flex h-screen w-full bg-muted/20 text-foreground overflow-hidden font-sans">
      {/*
        ========================================
        DESKTOP SIDEBAR (Floating Glass Island)
        ========================================
      */}
      <aside className="hidden md:flex flex-col w-[260px] h-[calc(100vh-24px)] m-3 mr-0 bg-card/80 backdrop-blur-xl border border-border shadow-2xl shadow-primary/5 rounded-4xl overflow-hidden z-50 relative">
        {/* 1. PROFILE BUTTON (Interactive) */}
        <div className="p-3">
          <Link to="/profile">
            <button className="cursor-pointer flex items-center gap-3 w-full p-2 rounded-2xl hover:bg-secondary/80 transition-colors border border-transparent hover:border-border group text-left">
              <div className="relative shrink-0">
                <div className="w-10 h-10 rounded-full bg-background border border-border overflow-hidden">
                  <img
                    src={avatarUrl}
                    alt="avatar"
                    className="w-full h-full object-cover"
                  />
                </div>
                {/* Online/Active Dot */}
                <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-card rounded-full" />
              </div>

              <div className="flex-1 min-w-0">
                <div className="font-semibold text-sm truncate text-foreground group-hover:text-primary transition-colors">
                  {user.email.split('@')[0]}
                </div>
                <div className="text-xs text-muted-foreground truncate">
                  View Profile
                </div>
              </div>
            </button>
          </Link>
        </div>

        <div className="h-px bg-border mx-4 opacity-50" />

        {/* 2. MAIN NAVIGATION */}
        <div className="flex-1 overflow-y-auto py-4 px-3 space-y-6 scrollbar-none">
          {/* Section: Main */}
          <div className="space-y-1">
            <div className="px-3 text-[10px] uppercase font-bold text-muted-foreground tracking-widest mb-2 opacity-70">
              {t('common.overview')}
            </div>
            {MAIN_NAV.map((item) => {
              const isActive = isPathActive(item.path);
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className="relative flex items-center gap-3 px-3 py-3 rounded-xl transition-all duration-300 group"
                >
                  {/* Liquid Active Background */}
                  {isActive && (
                    <motion.div
                      layoutId="desktop-nav-active"
                      className="absolute inset-0 bg-primary shadow-lg shadow-primary/20 rounded-xl"
                      initial={false}
                      transition={{
                        type: 'spring',
                        stiffness: 300,
                        damping: 30,
                      }}
                    />
                  )}

                  {/* Hover Background */}
                  {!isActive && (
                    <div className="absolute inset-0 bg-secondary/50 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity" />
                  )}

                  <item.icon
                    size={20}
                    className={cn(
                      'relative z-10 transition-colors duration-200',
                      isActive
                        ? 'text-primary-foreground'
                        : 'text-muted-foreground group-hover:text-foreground'
                    )}
                  />
                  <span
                    className={cn(
                      'relative z-10 font-medium text-sm transition-colors duration-200',
                      isActive
                        ? 'text-primary-foreground'
                        : 'text-muted-foreground group-hover:text-foreground'
                    )}
                  >
                    {item.label}
                  </span>
                </Link>
              );
            })}
          </div>

          {/* Section: Auxiliary */}
          <div className="space-y-1">
            <div className="px-3 text-[10px] uppercase font-bold text-muted-foreground tracking-widest mb-2 opacity-70">
              Other
            </div>
            {SECONDARY_NAV.map((item) => {
              const isActive = isPathActive(item.path);
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className="relative flex items-center gap-3 px-3 py-3 rounded-xl transition-all duration-300 group"
                >
                  {isActive && (
                    <motion.div
                      layoutId="desktop-nav-active"
                      className="absolute inset-0 bg-primary shadow-lg shadow-primary/20 rounded-xl"
                      initial={false}
                      transition={{
                        type: 'spring',
                        stiffness: 300,
                        damping: 30,
                      }}
                    />
                  )}

                  {!isActive && (
                    <div className="absolute inset-0 bg-secondary/50 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity" />
                  )}

                  <item.icon
                    size={20}
                    className={cn(
                      'relative z-10 transition-colors duration-200',
                      isActive
                        ? 'text-primary-foreground'
                        : 'text-muted-foreground group-hover:text-foreground'
                    )}
                  />
                  <span
                    className={cn(
                      'relative z-10 font-medium text-sm transition-colors duration-200',
                      isActive
                        ? 'text-primary-foreground'
                        : 'text-muted-foreground group-hover:text-foreground'
                    )}
                  >
                    {item.label}
                  </span>
                </Link>
              );
            })}
          </div>
        </div>

        {/* 3. BOTTOM ACTIONS */}
        <div className="p-3 border-t border-border mt-auto">
          <button
            onClick={logout}
            className="flex cursor-pointer w-full items-center gap-3 px-3 py-3 text-sm font-medium text-destructive hover:bg-destructive/10 rounded-xl transition-all duration-200 group"
          >
            <LogOutIcon
              size={20}
              className="opacity-70 group-hover:opacity-100"
            />
            <span>Log Out</span>
          </button>
        </div>
      </aside>

      {/*
        ========================================
        MAIN CONTENT AREA
        ========================================
      */}
      <main className="flex-1 relative h-full overflow-hidden">
        {/* Decorative Background Blobs (Using Semantic Colors) */}
        <div className="absolute top-[-20%] right-[-10%] w-[600px] h-[600px] bg-primary/5 rounded-full blur-[120px] pointer-events-none" />
        <div className="absolute bottom-[-20%] left-[-10%] w-[500px] h-[500px] bg-chart-4/5 rounded-full blur-[120px] pointer-events-none" />

        <div className="h-full overflow-y-auto overflow-x-hidden p-4 md:p-6 pb-32 md:pb-6 scroll-smooth">
          {/* Mobile Header */}
          <div className="md:hidden flex items-center justify-between mb-6 px-2">
            <h1 className="text-xl font-extrabold tracking-tight text-foreground flex items-center gap-2">
              <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-primary-foreground shadow-lg shadow-primary/20">
                <SparklesIcon size={18} />
              </div>
              {APP_TITLE}
            </h1>
            <Link to="/profile">
              <div className="w-9 h-9 rounded-full bg-secondary border border-border overflow-hidden">
                <img
                  src={avatarUrl}
                  alt="user"
                  className="w-full h-full object-cover"
                />
              </div>
            </Link>
          </div>

          <div className="max-w-6xl mx-auto">
            <Outlet />
          </div>
        </div>
      </main>

      {/*
        ========================================
        MOBILE FLOATING NAVBAR
        ========================================
      */}
      <nav className="md:hidden fixed bottom-6 left-6 right-6 h-[72px] bg-card/90 backdrop-blur-2xl border border-border rounded-4xl shadow-xl z-50 flex items-center justify-between px-2">
        {/* 1. PROFILE (First Item) */}
        <Link
          to="/profile"
          className="relative flex-1 flex flex-col items-center justify-center h-full group"
        >
          {isPathActive('/profile') && (
            <motion.div
              layoutId="mobile-nav-active"
              className="absolute w-12 h-12 bg-primary rounded-full shadow-lg shadow-primary/30 -z-10"
              transition={{ type: 'spring', stiffness: 300, damping: 25 }}
            />
          )}
          <UserIcon
            size={24}
            className={cn(
              'transition-all duration-300',
              isPathActive('/profile')
                ? 'text-primary-foreground scale-110'
                : 'text-muted-foreground group-active:scale-90'
            )}
            strokeWidth={isPathActive('/profile') ? 2.5 : 2}
          />
        </Link>

        {/* 2. MAIN NAV ITEMS */}
        {MAIN_NAV.map((item) => {
          const isActive = isPathActive(item.path);
          return (
            <Link
              key={item.path}
              to={item.path}
              className="relative flex-1 flex flex-col items-center justify-center h-full group"
            >
              {isActive && (
                <motion.div
                  layoutId="mobile-nav-active"
                  className="absolute w-12 h-12 bg-primary rounded-full shadow-lg shadow-primary/30 -z-10"
                  transition={{ type: 'spring', stiffness: 300, damping: 25 }}
                />
              )}

              <item.icon
                size={24}
                className={cn(
                  'transition-all duration-300',
                  isActive
                    ? 'text-primary-foreground scale-110'
                    : 'text-muted-foreground group-active:scale-90'
                )}
                strokeWidth={isActive ? 2.5 : 2}
              />
            </Link>
          );
        })}
      </nav>
    </div>
  );
};
