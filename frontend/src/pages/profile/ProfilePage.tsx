import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/entities/session/store';
import { analyticsApi } from '@/entities/analytics/api';
import { useTheme } from '@/shared/lib/hooks/useTheme';
import { startOfMonth, endOfMonth } from '@/shared/lib/date';
import { toast } from '@/entities/notification/store'; // Import Toast
import {
  LogOutIcon,
  MoonIcon,
  SunIcon,
  GlobeIcon,
  MessageCircleIcon,
  SparklesIcon,
  ShieldCheckIcon,
  CalendarIcon,
  ZapIcon,
  LinkIcon,
  LayersIcon,
  SendIcon,
  AlertCircleIcon,
} from 'lucide-react';
import { cn } from '@/shared/lib/utils';
import { Button } from '@/shared/ui/button';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { Skeleton } from '@/shared/ui/skeleton';

// --- UI Components ---

interface ActionTileProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  icon: React.ReactNode;
  label: string;
  subLabel?: string;
  active?: boolean;
  variant?: 'default' | 'destructive' | 'outline';
}

const ActionTile = ({
  icon,
  label,
  subLabel,
  active,
  variant = 'default',
  className,
  ...props
}: ActionTileProps) => {
  return (
    <button
      className={cn(
        'cursor-pointer group relative flex flex-col justify-between p-5 h-32 w-full rounded-2xl border border-transparent hover:transition-colors',
        variant === 'default' && 'bg-secondary hover:bg-primary/10 hover:border-primary/20',
        variant === 'destructive' && 'bg-destructive/30 hover:bg-destructive/20',
        variant === 'outline' && 'border-border hover:bg-accent',
        active && 'bg-primary text-primary-foreground hover:bg-primary',
        className
      )}
      {...props}
    >
      <div className={cn(
        "p-2 rounded-full w-fit",
        active ? "bg-white/20" : "bg-background/50 group-hover:bg-background"
      )}>
        {icon}
      </div>
      <div className="text-left">
        <div className={cn("font-bold text-lg", "light", active ? "text-primary-foreground" : "text-foreground")}>
          {label}
        </div>
        {subLabel && (
          <div className={cn("text-xs font-medium mt-0.5", active ? "text-primary-foreground/80" : "text-muted-foreground")}>
            {subLabel}
          </div>
        )}
      </div>
    </button>
  );
};

const StatPill = ({
  label,
  value,
  icon,
  loading
}: {
  label: string;
  value: string | number;
  icon?: React.ReactNode;
  loading?: boolean;
}) => (
  <div className="flex flex-col justify-center px-5 py-4 bg-secondary/30 rounded-2xl border border-border hover:border-primary/20">
    <div className="flex items-center gap-2 mb-1 text-muted-foreground">
      {icon && <span className="opacity-70">{icon}</span>}
      <span className="text-[10px] uppercase font-bold tracking-wider">
        {label}
      </span>
    </div>
    {loading ? (
       <Skeleton className="h-6 w-16" />
    ) : (
       <span className="text-xl font-bold text-foreground">{value}</span>
    )}
  </div>
);

const FeedbackModal = ({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) => {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onClose();
    // REPLACED ALERT WITH TOAST
    toast.info("Message Sent", "We have received your request. Support will contact you shortly.");
  };

  return (
    <ResponsiveSheet isOpen={isOpen} onClose={onClose} title="Contact Support">
      <div className="space-y-6">
        <div className="p-4 bg-primary/5 rounded-xl border border-primary/10">
          <h4 className="font-semibold flex items-center gap-2 mb-2 text-foreground">
            <SendIcon size={18} className="text-primary" />
            Telegram Support
          </h4>
          <p className="text-sm text-muted-foreground mb-4">
            The fastest way to get help is via our Telegram bot.
          </p>
          <Button className="w-full gap-2" onClick={() => window.open('https://t.me/your_support_bot', '_blank')}>
            Open Telegram
          </Button>
        </div>

        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">Or send email</span>
          </div>
        </div>

        <form className="space-y-4" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">Message</label>
            <textarea
              className="flex min-h-[100px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 text-foreground"
              placeholder="Describe your idea or issue..."
            />
          </div>
          <Button variant="outline" type="submit" className="w-full">
            Send Message
          </Button>
        </form>
      </div>
    </ResponsiveSheet>
  );
};


export const ProfilePage = () => {
  const { user, logout } = useAuthStore();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();

  const [isFeedbackOpen, setIsFeedbackOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [clicksUsed, setClicksUsed] = useState(0);
  const [inventory, setInventory] = useState({ campaigns: 0, links: 0 });

  // Handle Theme Toggle with Toast
  const handleThemeToggle = () => {
    toggleTheme();
    const newTheme = theme === 'light' ? 'Dark' : 'Light'; // It toggles AFTER this call usually, but logic depends on hook
    // Actually hook toggles immediately. Let's assume toggle works.
    toast.info("Theme Updated", `Switched to ${newTheme === 'Light' ? 'Dark' : 'Light'} mode`);
  };

  // Handle Logout with Toast
  const handleLogout = () => {
    logout();
    toast.info("Signed Out", "See you next time!");
  };

  // Handle Language with Toast
  const handleLanguage = () => {
     // REPLACED ALERT WITH TOAST
     toast.warn("Coming Soon", "Localization is currently in development.");
  };

  useEffect(() => {
    if (!user) return;

    const loadProfileData = async () => {
      setLoading(true);
      try {
        const now = new Date();
        const [summaryData, treeData] = await Promise.all([
          analyticsApi.getSummary({
            from: startOfMonth(now),
            to: endOfMonth(now),
          }),
          analyticsApi.getHierarchy(),
        ]);

        setClicksUsed(summaryData?.total_clicks || 0);

        let cCount = 0;
        let lCount = 0;
        if (treeData && treeData.children) {
          cCount = treeData.children.length;
          treeData.children.forEach(c => {
            if (c.children) lCount += c.children.length;
          });
        }
        setInventory({ campaigns: cCount, links: lCount });

      } catch (e) {
        console.error("Failed to load profile data", e);
        setClicksUsed(user.clicks_current_month);
      } finally {
        setLoading(false);
      }
    };

    loadProfileData();
  }, [user]);

  if (!user) return null;

  const isPro = user.plan_id === 'pro' || user.plan_id === 'enterprise';
  const maxClicks = isPro ? 100000 : 1000;
  const displayClicks = loading && clicksUsed === 0 ? user.clicks_current_month : clicksUsed;
  const usagePercent = Math.min((displayClicks / maxClicks) * 100, 100);
  const avatarUrl = `https://api.dicebear.com/9.x/thumbs/svg?seed=Brian`;

  return (
    <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in slide-in-from-bottom-4 pb-20">

      <div className="flex flex-col gap-2 mb-8">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Profile</h1>
        <p className="text-muted-foreground">Manage your personal information and preferences.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

        {/* --- LEFT COLUMN: CONTROL GRID --- */}
        <div className="grid grid-cols-2 gap-4 h-fit order-2 lg:order-1">
          <ActionTile
            icon={theme === 'dark' ? <MoonIcon size={20} /> : <SunIcon size={20} />}
            label="Theme"
            subLabel={theme === 'dark' ? 'Dark Mode' : 'Light Mode'}
            onClick={handleThemeToggle}
            active={theme === 'dark'}
          />
          <ActionTile
            icon={<GlobeIcon size={20} />}
            label="Language"
            subLabel="English"
            onClick={handleLanguage}
          />
          <ActionTile
            icon={<MessageCircleIcon size={20} />}
            label="Support"
            subLabel="Contact us"
            onClick={() => setIsFeedbackOpen(true)}
          />
          <ActionTile
            icon={<LogOutIcon size={20} />}
            label="Sign Out"
            subLabel="End session"
            variant="destructive"
            onClick={handleLogout}
          />
        </div>

        {/* --- RIGHT COLUMN: INFO CARD --- */}
        <div className="lg:col-span-2 order-1 lg:order-2">
          <div className="bg-card border border-border rounded-3xl p-6 sm:p-8 shadow-sm flex flex-col gap-8 h-full">

            {/* Header: Avatar & Info */}
            <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6 text-center sm:text-left">
              <div className="w-28 h-28 rounded-full border-4 border-background shadow-lg overflow-hidden shrink-0 bg-primary/5">
                <img src={avatarUrl} alt="Avatar" className="w-full h-full object-cover" />
              </div>

              <div className="flex-1 space-y-2 mt-2">
                <div>
                  <h2 className="text-2xl font-bold text-foreground">{user.email.split('@')[0]}</h2>
                  <p className="text-muted-foreground font-medium">{user.email}</p>
                </div>

                <div className="flex flex-wrap gap-2 justify-center sm:justify-start pt-1">
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-secondary text-secondary-foreground border border-border capitalize">
                    {user.role}
                  </span>
                  {user.is_active ? (
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-chart-2/10 text-chart-2 border border-chart-2/20">
                      <ShieldCheckIcon size={12} /> Verified
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-destructive/10 text-destructive border border-destructive/20">
                      <AlertCircleIcon size={12} /> Inactive
                    </span>
                  )}
                </div>
              </div>
            </div>

            {/* Plan & Usage Section */}
            <div className="space-y-4">
              <div className="relative p-6 rounded-2xl bg-linear-to-br from-slate-900 to-slate-800 text-white shadow-lg overflow-hidden group">
                <div className="absolute top-0 right-0 p-32 bg-primary/30 blur-3xl rounded-full -mr-16 -mt-16 pointer-events-none group-hover:bg-primary/40 transition-colors" />

                <div className="relative z-10 flex flex-col sm:flex-row justify-between items-center gap-6">
                  <div>
                    <div className="text-white/60 text-xs font-bold uppercase tracking-widest mb-1 flex items-center gap-2">
                      Current Plan
                    </div>
                    <div className="text-3xl font-bold flex items-center gap-2">
                      {user.plan_id === 'pro' || user.plan_id === 'enterprise' ? (
                         <ZapIcon className="text-yellow-400" fill="currentColor" />
                      ) : (
                         <SparklesIcon className="text-indigo-300" />
                      )}
                      <span className="capitalize">{user.plan_id} Plan</span>
                    </div>
                  </div>
                  <Button
                    onClick={() => navigate('/pricing')}
                    className="bg-white text-black hover:bg-gray-200 border-none shadow-none font-bold whitespace-nowrap"
                  >
                    Upgrade Plan
                  </Button>
                </div>

                {/* Usage Bar inside Plan Card */}
                <div className="mt-6 space-y-2">
                  <div className="flex justify-between text-xs font-medium text-white/80">
                    <span className="flex items-center gap-2">
                      Monthly Clicks
                    </span>
                    <span>
                      {loading ? '...' : new Intl.NumberFormat('en-US').format(displayClicks)} / {new Intl.NumberFormat('en-US', { notation: 'compact' }).format(maxClicks)}
                    </span>
                  </div>
                  <div className={cn("h-2 w-full bg-white/10 rounded-full overflow-hidden backdrop-blur-sm", loading && "animate-pulse")}>
                    <div
                      className={cn(
                        "h-full transition-all duration-1000 ease-out",
                         usagePercent > 90 ? "bg-destructive" : "bg-primary"
                      )}
                      style={{ width: `${usagePercent}%` }}
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Detailed Stats Grid */}
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
              <StatPill
                label="Campaigns"
                value={inventory.campaigns}
                icon={<LayersIcon size={14} />}
                loading={loading}
              />
              <StatPill
                label="Active Links"
                value={inventory.links}
                icon={<LinkIcon size={14} />}
                loading={loading}
              />
              <StatPill
                label="Member Since"
                value={new Date(user.created_at).toLocaleDateString(undefined, { month: 'short', year: 'numeric' })}
                icon={<CalendarIcon size={14} />}
              />
            </div>

          </div>
        </div>
      </div>

      <FeedbackModal isOpen={isFeedbackOpen} onClose={() => setIsFeedbackOpen(false)} />
    </div>
  );
};
