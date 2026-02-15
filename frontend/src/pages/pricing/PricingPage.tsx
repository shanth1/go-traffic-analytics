import { useState } from 'react';
import {
  CheckIcon,
  ZapIcon,
  XIcon,
  ShieldCheckIcon,
  WalletIcon,
  EyeIcon,
  InfinityIcon,
} from 'lucide-react';
import { Button } from '@/shared/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from '@/shared/ui/card';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { useAuthStore } from '@/entities/session/store';
import { useIsDesktop } from '@/shared/lib/hooks';
import { TG_SUPPORT_URL } from '@/shared/config';
import { useTranslation } from 'react-i18next';
import { cn } from '@/shared/lib/utils';

// --- STRICT TYPES ---
type TranslationParams = Record<string, string | number>;

interface FeatureDef {
  key: string;
  included: boolean;
  params?: TranslationParams;
  icon?: React.ElementType;
  highlight?: boolean;
  isNegative?: boolean;
  badge?: string;
}

interface PlanConfig {
  id: string;
  price: string;
  period: string;
  subPrice?: string;
  subPriceParams?: TranslationParams;
  highlight: boolean;
  badge: string | null;
  viewsEstimateKey: string;
  viewsEstimateParams?: TranslationParams;
  features: FeatureDef[];
}

export const PricingPage = () => {
  const { t } = useTranslation();
  const { user } = useAuthStore();
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);
  const isDesktop = useIsDesktop();

  // --- CONFIGURATION ---
  const plansConfig: PlanConfig[] = [
    {
      id: 'starter', // Lite
      price: '290 ₽',
      period: 'pricing.per_month',
      highlight: false,
      badge: null,
      viewsEstimateKey: 'pricing.features.views_estimate',
      viewsEstimateParams: { count: '160k' },
      features: [
        {
          key: 'pricing.features.clicks_limit',
          params: { count: '5,000' },
          included: true,
        },
        {
          key: 'pricing.features.links_limit',
          params: { count: '500' },
          included: true,
        },
        {
          key: 'pricing.features.retention',
          params: { days: '90' },
          included: true,
        },

        {
          key: 'pricing.features.no_export',
          included: false,
          isNegative: true,
        },
        { key: 'pricing.features.no_ml', included: false, isNegative: true },
        {
          key: 'pricing.features.clicks_burn',
          included: false,
          isNegative: true,
        },
      ],
    },
    {
      id: 'pro', // Pro
      price: '990 ₽',
      period: 'pricing.per_month',
      highlight: true,
      badge: 'pricing.most_popular',
      viewsEstimateKey: 'pricing.features.views_estimate',
      viewsEstimateParams: { count: '1.6M' },
      features: [
        {
          key: 'pricing.features.clicks_limit',
          params: { count: '50,000' },
          included: true,
        },
        { key: 'pricing.features.links_unlimited', included: true },
        {
          key: 'pricing.features.retention',
          params: { days: '365' },
          included: true,
        },
        { key: 'pricing.features.export', included: true },

        // SOON FEATURES (PRO)
        {
          key: 'pricing.features.ml',
          included: true,
          badge: 'pricing.badges.soon',
        },
        {
          key: 'pricing.features.utm',
          included: true,
          badge: 'pricing.badges.soon',
        },
        // {
        //   key: 'pricing.features.custom_slug',
        //   included: true,
        //   badge: 'pricing.badges.soon',
        // },
        {
          key: 'pricing.features.geo',
          included: true,
          badge: 'pricing.badges.soon',
        },
        // { key: 'pricing.features.support', included: true },
        {
          key: 'pricing.features.clicks_burn',
          included: false,
          isNegative: true,
        },
      ],
    },
    {
      id: 'enterprise', // Scale
      price: '0.05 ₽',
      period: 'pricing.per_click',
      subPrice: 'pricing.min_deposit',
      subPriceParams: { amount: '5,000' },
      highlight: false,
      badge: 'pricing.best_value',
      viewsEstimateKey: 'pricing.features.views_unlimited',
      features: [
        {
          key: 'pricing.features.clicks_no_burn',
          included: true,
          icon: InfinityIcon,
          highlight: true,
        },
        {
          key: 'pricing.features.all_pro_features',
          included: true,
          icon: CheckIcon,
        }, // Includes Pro

        { key: 'pricing.features.clicks_unlimited', included: true },
        { key: 'pricing.features.retention_lifetime', included: true },

        // ENTERPRISE EXCLUSIVES
        {
          key: 'pricing.features.custom_domain',
          included: true,
        }, // Manual setup - Ready
        { key: 'pricing.features.sla', included: true },

        // ENTERPRISE SOON
        {
          key: 'pricing.features.custom_error',
          included: true,
          badge: 'pricing.badges.soon',
        },
      ],
    },
  ];

  const handleAction = (planId: string) => {
    if (planId === user?.plan_id) return;
    setSelectedPlan(planId);
  };

  return (
    <div className="max-w-6xl mx-auto space-y-10 animate-in fade-in duration-500 py-10 px-4">
      <div className="text-center space-y-4 max-w-2xl mx-auto">
        <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight text-foreground">
          {t('pricing.title')}
        </h1>
        <p className="text-lg text-muted-foreground">{t('pricing.subtitle')}</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 pt-6 items-start">
        {plansConfig.map((plan) => {
          const isCurrent = user?.plan_id === plan.id;
          const isEnterprise = plan.id === 'enterprise';

          const sortedFeatures = [...plan.features].sort((a, b) => {
            if (a.isNegative === b.isNegative) return 0;
            return a.isNegative ? 1 : -1;
          });

          return (
            <Card
              key={plan.id}
              className={cn(
                'flex flex-col relative transition-all duration-300 h-full border-border overflow-visible',
                plan.highlight
                  ? 'border-primary shadow-2xl shadow-primary/10 scale-100 lg:scale-105 z-10 bg-card'
                  : 'bg-card/50 hover:border-primary/30 hover:bg-card',
                isEnterprise && 'bg-linear-to-b from-card to-secondary/30'
              )}
            >
              {plan.badge && (
                <div className="absolute -top-3.5 left-0 right-0 flex justify-center z-20">
                  <span
                    className={cn(
                      'text-xs font-bold px-4 py-1.5 rounded-full flex items-center gap-1.5 shadow-lg tracking-wider uppercase',
                      plan.highlight
                        ? 'bg-primary text-primary-foreground shadow-primary/20'
                        : 'bg-emerald-600 text-white shadow-emerald-500/20'
                    )}
                  >
                    {plan.highlight ? (
                      <ZapIcon size={12} fill="currentColor" />
                    ) : (
                      <ShieldCheckIcon size={12} />
                    )}
                    {t(plan.badge)}
                  </span>
                </div>
              )}

              <CardHeader className={cn(plan.highlight ? 'pt-8' : 'pt-6')}>
                <div className="flex justify-between items-start mb-2">
                  <CardTitle className="text-2xl font-bold">
                    {t(`pricing.plans.${plan.id}.name`)}
                  </CardTitle>
                  {isEnterprise && (
                    <ShieldCheckIcon className="text-muted-foreground opacity-50" />
                  )}
                </div>
                <CardDescription className="text-sm h-10 line-clamp-2">
                  {t(`pricing.plans.${plan.id}.desc`)}
                </CardDescription>
              </CardHeader>

              <CardContent className="flex-1 space-y-6">
                <div>
                  <div className="flex items-baseline gap-1">
                    <span className="text-5xl font-extrabold text-foreground tracking-tight">
                      {plan.price}
                    </span>
                    <span className="text-muted-foreground font-medium text-lg">
                      {t(plan.period)}
                    </span>
                  </div>
                  {plan.subPrice && (
                    <div className="mt-2 inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-xs font-semibold border border-emerald-500/20">
                      <WalletIcon size={12} />
                      {t(plan.subPrice, { ...plan.subPriceParams } as Record<
                        string,
                        string | number
                      >)}
                    </div>
                  )}
                </div>

                <div className="p-3 bg-secondary/50 rounded-lg border border-border/50 text-xs text-muted-foreground">
                  <div className="flex items-center gap-1.5">
                    <EyeIcon
                      size={14}
                      className={
                        isEnterprise ? 'text-emerald-500' : 'text-primary'
                      }
                    />
                    <span className="font-medium text-foreground">
                      {t(
                        plan.viewsEstimateKey,
                        plan.viewsEstimateParams as Record<
                          string,
                          string | number
                        >
                      )}
                    </span>
                  </div>
                </div>

                <div className="h-px bg-border/50 w-full" />

                <ul className="space-y-3">
                  {sortedFeatures.map((feature, idx) => {
                    const IconComponent = feature.icon
                      ? feature.icon
                      : feature.included
                        ? CheckIcon
                        : XIcon;

                    return (
                      <li
                        key={idx}
                        className={cn(
                          'flex items-start gap-3 text-sm min-h-6',
                          feature.highlight &&
                            'font-semibold text-emerald-600 dark:text-emerald-400',
                          feature.isNegative
                            ? 'text-muted-foreground opacity-60'
                            : 'text-foreground'
                        )}
                      >
                        <IconComponent
                          className={cn(
                            'h-5 w-5 shrink-0 mt-0.5',
                            feature.highlight
                              ? 'text-emerald-500'
                              : feature.included
                                ? 'text-primary'
                                : 'text-muted-foreground'
                          )}
                        />
                        <div className="flex flex-wrap items-center gap-1.5">
                          <span className="leading-tight">
                            {t(feature.key, { ...feature.params } as Record<
                              string,
                              string | number
                            >)}
                          </span>

                          {feature.badge && (
                            <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-bold uppercase tracking-wide bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border border-indigo-500/20">
                              {t(feature.badge)}
                            </span>
                          )}
                        </div>
                      </li>
                    );
                  })}
                </ul>
              </CardContent>

              <CardFooter className="pb-8 pt-4">
                <Button
                  className={cn(
                    'w-full h-12 text-base font-semibold transition-all duration-300',
                    isCurrent
                      ? 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
                      : plan.highlight
                        ? 'bg-primary text-primary-foreground shadow-lg shadow-primary/25 hover:shadow-primary/40 hover:-translate-y-0.5'
                        : isEnterprise
                          ? 'bg-emerald-600 text-white hover:bg-emerald-700 shadow-lg shadow-emerald-500/20'
                          : 'bg-primary/10 text-primary hover:bg-primary/20'
                  )}
                  disabled={isCurrent}
                  onClick={() => handleAction(plan.id)}
                >
                  {isCurrent
                    ? t('pricing.current_plan')
                    : t(`pricing.plans.${plan.id}.btn`)}
                </Button>
              </CardFooter>
            </Card>
          );
        })}
      </div>

      <ResponsiveSheet
        isOpen={!!selectedPlan}
        onClose={() => setSelectedPlan(null)}
        title={
          selectedPlan
            ? t('pricing.upgrade_modal.title', {
                plan: t(`pricing.plans.${selectedPlan}.name`),
              })
            : ''
        }
      >
        <div className="space-y-8 text-center py-6 px-2">
          <div className="p-5 bg-linear-to-br from-primary/10 to-purple-500/10 rounded-full w-20 h-20 mx-auto flex items-center justify-center text-primary shadow-inner">
            <WalletIcon size={36} />
          </div>

          <div className="space-y-3">
            <h3 className="text-xl font-bold text-foreground">
              {t('pricing.upgrade_modal.integration_title')}
            </h3>
            <p className="text-muted-foreground leading-relaxed">
              {t('pricing.upgrade_modal.integration_desc')}
            </p>
          </div>

          <div className="space-y-3 pt-4">
            <Button
              size="lg"
              className="w-full gap-2 h-12 text-base shadow-xl shadow-primary/20"
              onClick={() => window.open(TG_SUPPORT_URL, '_blank')}
            >
              {t('pricing.upgrade_modal.contact_btn')}
            </Button>

            {isDesktop && (
              <Button
                variant="ghost"
                onClick={() => setSelectedPlan(null)}
                className="w-full"
              >
                {t('common.cancel')}
              </Button>
            )}
          </div>
        </div>
      </ResponsiveSheet>
    </div>
  );
};
