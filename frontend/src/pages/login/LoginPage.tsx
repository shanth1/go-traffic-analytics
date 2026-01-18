import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/entities/session/store';
import { Input } from '@/shared/ui/input';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { LayoutDashboardIcon } from 'lucide-react';
import { AxiosError } from 'axios';
import { APP_TITLE } from '@/shared/config';

export const LoginPage = () => {
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
      navigate('/');
    } catch (error) {
      console.error(error);
      if (error instanceof AxiosError && error.response) {
        alert(
          `Login error: ${error.response.data?.error || 'Invalid credentials'}`
        );
      } else {
        alert('Network or server error');
      }
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
          <h1 className="text-2xl font-bold tracking-tight text-foreground">{APP_TITLE}</h1>
        </div>

        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-xl font-bold text-center text-foreground">
              Welcome back
            </CardTitle>
            <p className="text-sm text-muted-foreground text-center">
              Enter your credentials to access your dashboard
            </p>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleLogin} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground">Email</label>
                <Input
                  type="email"
                  placeholder="admin@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground">Password</label>
                <Input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <Button className="w-full" type="submit" disabled={loading}>
                {loading ? 'Signing in...' : 'Sign In'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
