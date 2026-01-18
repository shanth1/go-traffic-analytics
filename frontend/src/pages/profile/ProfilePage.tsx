import { useEffect, useState } from 'react';
import { useAuthStore } from '@/entities/session/store';
import { analyticsApi } from '@/entities/analytics/api';
import { startOfMonth, endOfMonth } from '@/shared/lib/date';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import {
  UserIcon,
  CreditCardIcon,
  SettingsIcon,
  ShieldIcon,
  LogOutIcon,
  GlobeIcon,
  RefreshCwIcon,
} from 'lucide-react';

export const ProfilePage = () => {
  const { user, logout } = useAuthStore();
  const [realtimeClicks, setRealtimeClicks] = useState<number>(0);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchUsage = async () => {
      setIsLoading(true);
      try {
        const now = new Date();
        const data = await analyticsApi.getSummary({
          from: startOfMonth(now),
          to: endOfMonth(now),
        });
        setRealtimeClicks(data?.total_clicks || 0);
      } catch (e) {
        console.error('Failed to sync usage', e);
        if (user) setRealtimeClicks(user.clicks_current_month);
      } finally {
        setIsLoading(false);
      }
    };
    fetchUsage();
  }, [user]);

  if (!user) return null;

  const isPro = user.plan_id === 'pro' || user.plan_id === 'enterprise';
  const maxClicks = isPro ? 100000 : 1000;

  const displayClicks =
    isLoading && realtimeClicks === 0
      ? user.clicks_current_month
      : realtimeClicks;
  const clicksUsage = (displayClicks / maxClicks) * 100;

  return (
    <div className="space-y-8 animate-in fade-in duration-500 max-w-4xl mx-auto">
      {/* Header Profile Info */}
      <div className="flex items-center gap-4 mb-8">
        <div className="h-20 w-20 rounded-full bg-primary/10 flex items-center justify-center text-primary">
          <UserIcon size={40} />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-foreground">{user.email.split('@')[0]}</h1>
          <p className="text-muted-foreground">{user.email}</p>
          <div className="flex items-center gap-2 mt-2">
            <span className="px-2 py-0.5 rounded-full bg-secondary text-xs font-medium border border-border capitalize text-foreground">
              {user.role}
            </span>
            <span
              className={`px-2 py-0.5 rounded-full text-xs font-medium border capitalize ${
                user.plan_id === 'free'
                  ? 'bg-secondary border-border text-muted-foreground'
                  : 'bg-chart-3/10 border-chart-3/20 text-chart-3'
              }`}
            >
              {user.plan_id} Plan
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Usage Stats (Fresh Data) */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CreditCardIcon size={20} className="text-muted-foreground" />
              Subscription & Usage
            </CardTitle>
            <CardDescription>
              Usage for{' '}
              {new Date().toLocaleString('default', {
                month: 'long',
                year: 'numeric',
              })}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="font-medium flex items-center gap-2 text-foreground">
                  Total Clicks
                  {isLoading && (
                    <RefreshCwIcon
                      size={12}
                      className="animate-spin text-muted-foreground"
                    />
                  )}
                </span>
                <span className="text-muted-foreground tabular-nums">
                  {new Intl.NumberFormat('en-US').format(displayClicks)} /{' '}
                  {new Intl.NumberFormat('en-US', {
                    notation: 'compact',
                  }).format(maxClicks)}
                </span>
              </div>
              <div className="h-2 w-full bg-secondary rounded-full overflow-hidden">
                <div
                  className={`h-full transition-all duration-1000 ease-out ${
                    clicksUsage > 90 ? 'bg-destructive' : 'bg-primary'
                  }`}
                  style={{ width: `${Math.min(clicksUsage, 100)}%` }}
                />
              </div>
              <p className="text-xs text-muted-foreground text-right">
                Refreshes automatically based on real-time analytics
              </p>
            </div>

            <div className="pt-4 border-t border-border">
              <Button variant="outline" className="w-full">
                Upgrade Plan
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Settings Stubs */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <SettingsIcon size={20} className="text-muted-foreground" />
              Preferences
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-sm font-medium text-foreground">
                <GlobeIcon size={16} /> Language
              </div>
              <select className="h-8 rounded-md border border-input bg-transparent px-2 text-xs focus:ring-ring text-foreground">
                <option>English</option>
                <option disabled>Russian (Coming Soon)</option>
              </select>
            </div>

            <div className="flex items-center justify-between pt-2">
              <div className="flex items-center gap-2 text-sm font-medium text-foreground">
                <ShieldIcon size={16} /> Password
              </div>
              <Button variant="ghost" size="sm" className="h-8">
                Change
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-destructive/20">
        <CardHeader>
          <CardTitle className="text-destructive text-lg">
            Session Control
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Button variant="destructive" className="gap-2" onClick={logout}>
            <LogOutIcon size={16} /> Sign Out
          </Button>
        </CardContent>
      </Card>
    </div>
  );
};
