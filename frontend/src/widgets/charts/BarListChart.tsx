import { useMemo } from 'react';
import type { CategoryStat } from '@/shared/api/types';
import { cn } from '@/shared/lib/utils';

interface BarListChartProps {
  data: CategoryStat[];
  color?: string;
  className?: string;
}

export const BarListChart = ({
  data,
  color = 'bg-indigo-600',
  className,
}: BarListChartProps) => {
  // Sort descending just in case API didn't
  const sortedData = useMemo(() => {
    return [...data].sort((a, b) => b.value - a.value);
  }, [data]);

  const maxValue = sortedData[0]?.value || 1;

  if (data.length === 0) {
    return <div className="text-sm text-slate-400 py-4">No data available</div>;
  }

  return (
    <div className={cn('space-y-3', className)}>
      {sortedData.map((item, idx) => {
        // Calculate width relative to the largest item in the list
        const widthPercent = (item.value / maxValue) * 100;

        return (
          <div key={item.name + idx} className="group">
            <div className="flex justify-between text-xs mb-1">
              <span
                className="font-medium text-slate-700 dark:text-slate-200 truncate pr-2"
                title={item.name}
              >
                {item.name}
              </span>
              <span className="text-slate-500 tabular-nums">
                {new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                  item.value
                )}
              </span>
            </div>
            <div className="h-2 w-full bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
              <div
                className={cn(
                  'h-full rounded-full transition-all duration-500',
                  color
                )}
                style={{ width: `${widthPercent}%` }}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
};
