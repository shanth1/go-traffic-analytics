import { useState } from 'react';
import { CheckIcon, ZapIcon } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/shared/ui/card';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { useAuthStore } from '@/entities/session/store';

const PLANS = [
  {
    id: 'free',
    name: 'Starter',
    price: '$0',
    description: 'Perfect for side projects and hobbies.',
    features: ['1,000 clicks/mo', 'Unlimited Links', 'Basic Analytics', '7-day Data Retention'],
    buttonText: 'Current Plan',
    highlight: false,
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '$29',
    period: '/mo',
    description: 'For creators and growing businesses.',
    features: ['100,000 clicks/mo', 'Custom Slugs', 'Geo & Device Analytics', '90-day Data Retention', 'Priority Support'],
    buttonText: 'Upgrade to Pro',
    highlight: true,
  },
  {
    id: 'enterprise',
    name: 'Enterprise',
    price: 'Custom',
    description: 'For large scale organizations.',
    features: ['Unlimited clicks', 'SSO & SLA', 'Dedicated Manager', 'Custom Contracts', 'Raw Data Export'],
    buttonText: 'Contact Sales',
    highlight: false,
  },
];

export const PricingPage = () => {
  const { user } = useAuthStore();
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);

  const handleAction = (planId: string) => {
    if (planId === user?.plan_id) return;
    setSelectedPlan(planId);
  };

  return (
    <div className="max-w-5xl mx-auto space-y-8 animate-in fade-in duration-500 py-10">
      <div className="text-center space-y-4 max-w-2xl mx-auto">
        <h1 className="text-4xl font-bold tracking-tight text-foreground">
          Simple, transparent pricing
        </h1>
        <p className="text-lg text-muted-foreground">
          Choose the plan that's right for you. All plans include our core link shortening features.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 pt-8">
        {PLANS.map((plan) => {
          const isCurrent = user?.plan_id === plan.id;

          return (
            <Card
              key={plan.id}
              className={`flex flex-col relative transition-all duration-200 ${
                plan.highlight
                  ? 'border-primary shadow-lg scale-105 z-10 bg-card'
                  : 'bg-secondary/20 border-border hover:border-primary/30'
              }`}
            >
              {plan.highlight && (
                <div className="absolute -top-4 left-0 right-0 flex justify-center">
                  <span className="bg-primary text-primary-foreground text-xs font-bold px-3 py-1 rounded-full flex items-center gap-1 shadow-sm">
                    <ZapIcon size={12} fill="currentColor" />
                    MOST POPULAR
                  </span>
                </div>
              )}

              <CardHeader>
                <CardTitle className="text-xl">{plan.name}</CardTitle>
                <CardDescription>{plan.description}</CardDescription>
              </CardHeader>
              <CardContent className="flex-1 space-y-6">
                <div className="flex items-baseline gap-1">
                  <span className="text-4xl font-bold text-foreground">{plan.price}</span>
                  {plan.period && <span className="text-muted-foreground font-medium">{plan.period}</span>}
                </div>

                <ul className="space-y-3">
                  {plan.features.map((feature) => (
                    <li key={feature} className="flex items-start gap-3 text-sm text-muted-foreground">
                      <CheckIcon className="h-5 w-5 text-primary shrink-0" />
                      {feature}
                    </li>
                  ))}
                </ul>
              </CardContent>
              <CardFooter>
                <Button
                  className="w-full"
                  variant={isCurrent ? "outline" : (plan.highlight ? "default" : "secondary")}
                  disabled={isCurrent}
                  onClick={() => handleAction(plan.id)}
                >
                  {isCurrent ? "Current Plan" : plan.buttonText}
                </Button>
              </CardFooter>
            </Card>
          );
        })}
      </div>

      {/* Contact Modal for Upgrade */}
      <ResponsiveSheet
        isOpen={!!selectedPlan}
        onClose={() => setSelectedPlan(null)}
        title={`Upgrade to ${PLANS.find(p => p.id === selectedPlan)?.name}`}
      >
        <div className="space-y-6 text-center py-4">
          <div className="p-4 bg-primary/10 rounded-full w-16 h-16 mx-auto flex items-center justify-center text-primary">
            <ZapIcon size={32} />
          </div>
          <div className="space-y-2">
            <h3 className="text-lg font-semibold">Payment Gateway Integration</h3>
            <p className="text-muted-foreground text-sm">
              We are currently integrating Stripe. To upgrade your plan immediately, please contact our support team, and we will handle it manually within 1 hour.
            </p>
          </div>

          <Button className="w-full" onClick={() => window.open('https://t.me/your_support_bot', '_blank')}>
            Contact Support to Upgrade
          </Button>

          <Button variant="ghost" onClick={() => setSelectedPlan(null)} className="w-full">
            Cancel
          </Button>
        </div>
      </ResponsiveSheet>
    </div>
  );
};
