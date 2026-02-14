// FILE: frontend/src/pages/ml/MlAnalyticsPage.tsx
import { useState } from 'react';
import {
  BrainCircuitIcon,
  SparklesIcon,
  TrendingUpIcon,
  LineChartIcon,
  ShieldAlertIcon,
  Loader2Icon,
} from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Card, CardContent } from '@/shared/ui/card';
import { useTranslation } from 'react-i18next';
import { toast } from '@/entities/notification/store';
import { feedbackApi } from '@/entities/feedback/api';
import { useAuthStore } from '@/entities/session/store';

export const MlAnalyticsPage = () => {
  const { t } = useTranslation();
  const { user } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const handleNotify = async () => {
    setLoading(true);
    try {
      await feedbackApi.send({
        name: user?.email ? user.email.split('@')[0] : 'Guest',
        email: user?.email || 'unknown@user.com',
        message:
          'User requested early access to ML/AI Analytics features via Waitlist button.',
      });

      toast.info(t('ml.notify_success'), t('ml.notify_success_desc'));
    } catch (error) {
      console.error('Failed to subscribe to ML waitlist', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-8 animate-in fade-in duration-500 min-h-[calc(100vh-100px)] flex flex-col items-center justify-center relative overflow-hidden">
      {/* Background Decor */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/5 rounded-full blur-[100px] -z-10" />

      <div className="max-w-4xl w-full space-y-12 text-center z-10 px-4">
        {/* Header */}
        <div className="space-y-6">
          <div className="w-24 h-24 bg-linear-to-br from-primary to-purple-600 rounded-3xl flex items-center justify-center mx-auto shadow-2xl shadow-primary/30 mb-8 animate-in zoom-in duration-500">
            <BrainCircuitIcon size={48} className="text-white" />
          </div>

          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight text-foreground">
            {t('ml.title')}
          </h1>

          <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            {t('ml.desc')}
          </p>

          <div className="flex justify-center pt-2">
            <span className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-primary/10 text-primary text-sm font-bold uppercase tracking-wider border border-primary/20">
              <SparklesIcon size={14} /> {t('ml.comming_soon')}
            </span>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left">
          {/* Feature 1: Anomaly */}
          <Card className="bg-card/40 backdrop-blur-md border-primary/10 hover:border-primary/30 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
            <CardContent className="p-8 space-y-5">
              <div className="p-3 bg-red-500/10 text-red-500 rounded-2xl w-fit">
                <ShieldAlertIcon size={28} />
              </div>
              <div>
                <h3 className="font-bold text-xl mb-2">{t('ml.feature_1')}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {t('ml.feature_1_desc')}
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Feature 2: Prediction */}
          <Card className="bg-card/40 backdrop-blur-md border-primary/10 hover:border-primary/30 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
            <CardContent className="p-8 space-y-5">
              <div className="p-3 bg-blue-500/10 text-blue-500 rounded-2xl w-fit">
                <LineChartIcon size={28} />
              </div>
              <div>
                <h3 className="font-bold text-xl mb-2">{t('ml.feature_2')}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {t('ml.feature_2_desc')}
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Feature 3: Smart Optimization */}
          <Card className="bg-card/40 backdrop-blur-md border-primary/10 hover:border-primary/30 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
            <CardContent className="p-8 space-y-5">
              <div className="p-3 bg-emerald-500/10 text-emerald-500 rounded-2xl w-fit">
                <TrendingUpIcon size={28} />
              </div>
              <div>
                <h3 className="font-bold text-xl mb-2">{t('ml.feature_3')}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {t('ml.feature_3_desc')}
                </p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* CTA */}
        <div className="pt-8">
          <Button
            size="lg"
            className="h-14 px-10 text-lg font-semibold gap-3 shadow-xl shadow-primary/20 hover:scale-105 transition-transform"
            onClick={handleNotify}
            disabled={loading}
          >
            {loading ? (
              <Loader2Icon className="animate-spin" size={20} />
            ) : (
              <SparklesIcon size={20} />
            )}
            {t('ml.notify_btn')}
          </Button>
          <p className="text-xs font-medium text-muted-foreground mt-6 uppercase tracking-widest opacity-60">
            Available for Pro & Enterprise
          </p>
        </div>
      </div>
    </div>
  );
};
