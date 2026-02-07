import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '@/entities/session/store';
import { Input } from '@/shared/ui/input';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { LayoutDashboardIcon } from 'lucide-react';
import { APP_TITLE } from '@/shared/config';
import { toast } from '@/entities/notification/store';
import { useTranslation } from 'react-i18next';

export const LoginPage = () => {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const login = useAuthStore((state) => state.login);
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login({ email, password });
      toast.info(t('auth.welcome_toast_title'), t('auth.welcome_toast_desc'));
      navigate('/');
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/20 p-4">
      <div className="w-full max-w-md space-y-6">
        <div className="flex flex-col items-center gap-2 text-center">
          <div className="p-3 bg-primary/10 rounded-xl text-primary">
            <LayoutDashboardIcon size={32} />
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            {APP_TITLE}
          </h1>
        </div>

        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-xl font-bold text-center text-foreground">
              {t('auth.title')}
            </CardTitle>
            <p className="text-sm text-muted-foreground text-center">
              {t('auth.subtitle')}
            </p>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleLogin} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground">
                  {t('auth.email_label')}
                </label>
                <Input
                  type="email"
                  placeholder={t('auth.email_placeholder')}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground">
                  {t('auth.password_label')}
                </label>
                <Input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <Button className="w-full" type="submit" disabled={loading}>
                {loading ? t('auth.signing_in') : t('auth.sign_in_btn')}
              </Button>
            </form>

            <div className="mt-6 text-center text-sm">
              <span className="text-muted-foreground">
                {t('auth.no_account')}{' '}
              </span>
              <Link
                to="/register"
                className="font-semibold text-primary hover:underline underline-offset-4 transition-colors"
              >
                {t('auth.sign_up')}
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
