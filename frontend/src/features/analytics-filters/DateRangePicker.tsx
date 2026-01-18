import React from 'react';
import { CalendarIcon } from 'lucide-react';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import type { DatePreset } from '@/entities/analytics/model/filters';
import { formatDateInput } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';
import { Button } from '@/shared/ui/button';

export const DateRangePicker = () => {
  const { preset, startDate, endDate, setPreset, setCustomRange } =
    useAnalyticsFilter();

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
    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 bg-card p-1 rounded-lg border border-border">
      <div className="flex items-center p-2 text-muted-foreground">
        <CalendarIcon size={16} />
      </div>

      <div className="flex items-center gap-1 bg-secondary rounded-md p-1">
        {(['7d', '30d', '90d'] as DatePreset[]).map((p) => (
          <Button
            key={p}
            variant="ghost"
            size="sm"
            onClick={() => handlePresetChange(p)}
            className={cn(
              'h-7 text-xs px-3 cursor-pointer',
              preset === p &&
                'bg-background shadow-sm text-primary'
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
            'h-7 text-xs px-3',
            preset === 'custom' &&
              'bg-background shadow-sm text-primary cursor-pointer'
          )}
        >
          Custom
        </Button>
      </div>

      {preset === 'custom' && (
        <div className="flex items-center gap-2 px-2 animate-in fade-in slide-in-from-left-2 duration-200">
          <input
            type="date"
            className="bg-transparent text-xs font-medium focus:outline-none text-foreground"
            value={formatDateInput(startDate)}
            onChange={handleStartChange}
            max={formatDateInput(new Date())}
          />
          <span className="text-muted-foreground">-</span>
          <input
            type="date"
            className="bg-transparent text-xs font-medium focus:outline-none text-foreground"
            value={formatDateInput(endDate)}
            onChange={handleEndChange}
            max={formatDateInput(new Date())}
          />
        </div>
      )}
    </div>
  );
};
