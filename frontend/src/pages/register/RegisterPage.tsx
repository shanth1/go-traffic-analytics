import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { SendIcon, SparklesIcon, ArrowLeft } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { TG_SUPPORT_URL, APP_TITLE } from '@/shared/config';

export const RegisterPage = () => {
  const { t } = useTranslation();

  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/20 p-4 animate-in fade-in zoom-in-95 duration-500">
      <div className="w-full max-w-md space-y-6">
        {/* Logo / Brand */}
        <div className="flex flex-col items-center gap-2 text-center">
          <div className="p-3 bg-primary/10 rounded-xl text-primary shadow-lg shadow-primary/20">
            <SparklesIcon size={32} />
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            {APP_TITLE}
          </h1>
        </div>

        <Card className="border-border shadow-xl">
          <CardHeader className="space-y-1 text-center pb-2">
            <CardTitle className="text-xl font-bold text-foreground">
              {t('auth.register_title')}
            </CardTitle>
            <p className="text-sm text-muted-foreground leading-relaxed px-4">
              {t('auth.register_subtitle')}
            </p>
          </CardHeader>

          <CardContent className="space-y-4 pt-4">
            {/* Main CTA Button */}
            <Button
              className="w-full h-12 gap-2 text-base font-medium shadow-md hover:shadow-lg transition-all"
              onClick={() => window.open(TG_SUPPORT_URL, '_blank')}
            >
              <SendIcon size={18} />
              {t('auth.contact_support')}
            </Button>

            {/* Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">Or</span>
              </div>
            </div>

            {/* Back to Login */}
            <Link to="/login" className="block">
              <Button variant="outline" className="w-full gap-2">
                <ArrowLeft size={16} />
                {t('auth.back_to_login')}
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
