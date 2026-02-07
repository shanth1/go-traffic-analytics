import React from 'react';
import { CalendarIcon } from 'lucide-react';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import type { DatePreset } from '@/entities/analytics/model/filters';
import { formatDateInput } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';
import { Button } from '@/shared/ui/button';
import { motion, AnimatePresence } from 'framer-motion';
import { useTheme } from '@/shared/lib/hooks/useTheme';

export const DateRangePicker = () => {
  const { preset, startDate, endDate, setPreset, setCustomRange } =
    useAnalyticsFilter();

  const { theme } = useTheme();

  const handlePresetChange = (p: DatePreset) => {
    setPreset(p);
  };

  const handleStartChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.value) {
      setCustomRange(new Date(e.target.value), endDate);
    }
  };

  const handleEndChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.value) {
      setCustomRange(startDate, new Date(e.target.value));
    }
  };

  return (
    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 bg-card p-1.5 rounded-xl border border-border shadow-sm w-full sm:w-auto overflow-hidden">
      {/* Icon & Presets Group */}
      <div className="flex items-center gap-1 w-full sm:w-auto overflow-x-auto no-scrollbar">
        <div className="flex items-center px-2 text-muted-foreground shrink-0">
          <CalendarIcon size={16} />
        </div>

        <div className="flex items-center gap-1 bg-secondary/50 rounded-lg p-1 shrink-0">
          {(['7d', '30d', '90d'] as DatePreset[]).map((p) => (
            <Button
              key={p}
              variant="ghost"
              size="sm"
              onClick={() => handlePresetChange(p)}
              className={cn(
                'h-7 text-xs px-3 rounded-md transition-all duration-200',
                preset === p
                  ? 'bg-background text-primary shadow-sm font-semibold hover:bg-background'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              {p.toUpperCase()}
            </Button>
          ))}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handlePresetChange('custom')}
            className={cn(
              'h-7 text-xs px-3 rounded-md transition-all duration-200',
              preset === 'custom'
                ? 'bg-background text-primary shadow-sm font-semibold hover:bg-background'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            Custom
          </Button>
        </div>
      </div>

      {/* Custom Inputs with Fixed Width Container */}
      <AnimatePresence>
        {preset === 'custom' && (
          <motion.div
            initial={{ width: 0, opacity: 0 }}
            animate={{ width: 'auto', opacity: 1 }}
            exit={{ width: 0, opacity: 0 }}
            transition={{ type: 'spring', stiffness: 400, damping: 30 }}
            className="overflow-hidden sm:border-l sm:border-border"
          >
            <div className="w-[260px] flex items-center gap-2 pl-3 pr-1">
              <input
                type="date"
                // Using inline style for colorScheme ensures the browser renders the correct native icon color
                style={{ colorScheme: theme }}
                className={cn(
                  'bg-transparent text-xs font-medium focus:outline-none cursor-pointer uppercase',
                  'text-muted-foreground hover:text-foreground transition-colors',
                  // Reset any previous filters
                  '[&::-webkit-calendar-picker-indicator]:filter-none',
                  // Ensure proper spacing and cursor
                  '[&::-webkit-calendar-picker-indicator]:cursor-pointer [&::-webkit-calendar-picker-indicator]:ml-1'
                )}
                value={formatDateInput(startDate)}
                onChange={handleStartChange}
                max={formatDateInput(new Date())}
              />
              <span className="text-muted-foreground/50 text-[10px] uppercase font-bold">
                To
              </span>
              <input
                type="date"
                style={{ colorScheme: theme }}
                className={cn(
                  'bg-transparent text-xs font-medium focus:outline-none cursor-pointer uppercase',
                  'text-muted-foreground hover:text-foreground transition-colors',
                  '[&::-webkit-calendar-picker-indicator]:filter-none',
                  '[&::-webkit-calendar-picker-indicator]:cursor-pointer [&::-webkit-calendar-picker-indicator]:ml-1'
                )}
                value={formatDateInput(endDate)}
                onChange={handleEndChange}
                max={formatDateInput(new Date())}
              />
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
