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
    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 bg-white dark:bg-slate-950 p-1 rounded-lg border border-slate-200 dark:border-slate-800">
      <div className="flex items-center p-2 text-slate-500">
        <CalendarIcon size={16} />
      </div>

      <div className="flex items-center gap-1 bg-slate-100 dark:bg-slate-900 rounded-md p-1">
        {(['7d', '30d', '90d'] as DatePreset[]).map((p) => (
          <Button
            key={p}
            variant="ghost"
            size="sm"
            onClick={() => handlePresetChange(p)}
            className={cn(
              'h-7 text-xs px-3',
              preset === p &&
                'bg-white dark:bg-slate-800 shadow-sm text-indigo-600 dark:text-indigo-400'
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
              'bg-white dark:bg-slate-800 shadow-sm text-indigo-600 dark:text-indigo-400'
          )}
        >
          Custom
        </Button>
      </div>

      {preset === 'custom' && (
        <div className="flex items-center gap-2 px-2 animate-in fade-in slide-in-from-left-2 duration-200">
          <input
            type="date"
            className="bg-transparent text-xs font-medium focus:outline-none dark:text-slate-200"
            value={formatDateInput(startDate)}
            onChange={handleStartChange}
            max={formatDateInput(new Date())}
          />
          <span className="text-slate-400">-</span>
          <input
            type="date"
            className="bg-transparent text-xs font-medium focus:outline-none dark:text-slate-200"
            value={formatDateInput(endDate)}
            onChange={handleEndChange}
            max={formatDateInput(new Date())}
          />
        </div>
      )}
    </div>
  );
};
