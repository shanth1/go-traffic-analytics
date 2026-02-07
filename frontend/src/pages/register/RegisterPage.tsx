import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  SendIcon,
  SparklesIcon,
  ArrowLeft,
  CheckCircle2Icon,
  MessageCircleIcon,
} from 'lucide-react';

import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { TG_SUPPORT_URL, APP_TITLE } from '@/shared/config';
import { feedbackApi } from '@/entities/feedback/api';
import { toast } from '@/entities/notification/store';

export const RegisterPage = () => {
  const { t } = useTranslation();

  // Form State
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');

  // UI State
  const [isSending, setIsSending] = useState(false);
  const [isSent, setIsSent] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !email.trim()) return;

    setIsSending(true);
    try {
      await feedbackApi.send({
        name,
        email,
        message: message || 'Requesting access via registration page',
      });

      setIsSent(true);
      toast.info(t('auth.success_title'), t('auth.success_desc'));
    } catch (error) {
      console.error('Registration request failed', error);
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/20 p-4 animate-in fade-in zoom-in-95 duration-500">
      <div className="w-full max-w-md space-y-6">
        {/* Brand Header */}
        <div className="flex flex-col items-center gap-2 text-center">
          <div className="p-3 bg-primary/10 rounded-xl text-primary shadow-lg shadow-primary/20">
            <SparklesIcon size={32} />
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            {APP_TITLE}
          </h1>
        </div>

        <Card className="border-border shadow-xl overflow-hidden">
          {isSent ? (
            /* --- STATE: SUCCESS --- */
            <div className="p-8 flex flex-col items-center text-center space-y-4">
              <div className="w-16 h-16 bg-green-100 text-green-600 rounded-full flex items-center justify-center mb-2">
                <CheckCircle2Icon size={32} />
              </div>
              <h2 className="text-xl font-bold text-foreground">
                {t('auth.success_title')}
              </h2>
              <p className="text-muted-foreground text-sm leading-relaxed">
                {t('auth.success_desc')}
              </p>
              <div className="pt-4 w-full">
                <Link to="/login" className="w-full">
                  <Button variant="outline" className="w-full gap-2">
                    <ArrowLeft size={16} />
                    {t('auth.back_to_login')}
                  </Button>
                </Link>
              </div>
            </div>
          ) : (
            /* --- STATE: FORM --- */
            <>
              <CardHeader className="space-y-1 text-center pb-2">
                <CardTitle className="text-xl font-bold text-foreground">
                  {t('auth.register_title')}
                </CardTitle>
                <p className="text-sm text-muted-foreground px-2">
                  {t('auth.register_subtitle')}
                </p>
              </CardHeader>

              <CardContent className="space-y-6 pt-2">
                {/* 1. Telegram Option (Highlighted) */}
                <div className="p-4 bg-primary/5 rounded-xl border border-primary/10 space-y-3">
                  <div className="flex items-center gap-2 font-semibold text-foreground text-sm">
                    <SendIcon size={16} className="text-primary" />
                    {t('auth.telegram_title')}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {t('auth.telegram_desc')}
                  </p>
                  <Button
                    className="w-full h-10 gap-2 text-sm shadow-sm"
                    onClick={() => window.open(TG_SUPPORT_URL, '_blank')}
                  >
                    {t('auth.open_telegram')}
                  </Button>
                </div>

                {/* Divider */}
                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <span className="w-full border-t border-border" />
                  </div>
                  <div className="relative flex justify-center text-xs uppercase font-medium">
                    <span className="bg-card px-2 text-muted-foreground">
                      {t('auth.or_form')}
                    </span>
                  </div>
                </div>

                {/* 2. Manual Form */}
                <form onSubmit={handleSubmit} className="space-y-3">
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-foreground ml-1">
                      {t('auth.form_name')}{' '}
                      <span className="text-destructive">*</span>
                    </label>
                    <Input
                      placeholder="John Doe"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-foreground ml-1">
                      {t('auth.form_email')}{' '}
                      <span className="text-destructive">*</span>
                    </label>
                    <Input
                      type="email"
                      placeholder="john@company.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      required
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-foreground ml-1">
                      {t('auth.form_message')}
                    </label>
                    <textarea
                      className="flex min-h-20 w-full rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50 resize-none"
                      placeholder="..."
                      value={message}
                      onChange={(e) => setMessage(e.target.value)}
                    />
                  </div>

                  <Button
                    type="submit"
                    className="w-full mt-2"
                    variant="secondary"
                    disabled={isSending}
                  >
                    {isSending ? (
                      t('auth.form_sending')
                    ) : (
                      <>
                        <MessageCircleIcon size={16} className="mr-2" />
                        {t('auth.form_submit')}
                      </>
                    )}
                  </Button>
                </form>

                {/* Footer Link */}
                <div className="pt-2 text-center">
                  <Link
                    to="/login"
                    className="text-xs text-muted-foreground hover:text-primary transition-colors flex items-center justify-center gap-1"
                  >
                    <ArrowLeft size={12} />
                    {t('auth.back_to_login')}
                  </Link>
                </div>
              </CardContent>
            </>
          )}
        </Card>
      </div>
    </div>
  );
};
