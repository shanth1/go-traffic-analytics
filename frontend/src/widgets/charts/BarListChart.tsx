import { useMemo } from 'react';
import type { CategoryStat } from '@/shared/api/types';
import { cn } from '@/shared/lib/utils';

interface BarListChartProps {
  data: CategoryStat[];
  color?: string; // Expects a Tailwind class like 'bg-primary' or 'bg-chart-1'
  className?: string;
}

export const BarListChart = ({
  data,
  color = 'bg-primary',
  className,
}: BarListChartProps) => {
  const sortedData = useMemo(() => {
    return [...data].sort((a, b) => b.value - a.value);
  }, [data]);

  const maxValue = sortedData[0]?.value || 1;

  if (data.length === 0) {
    return <div className="text-sm text-muted-foreground py-4">No data available</div>;
  }

  return (
    <div className={cn('space-y-3', className)}>
      {sortedData.map((item, idx) => {
        const widthPercent = (item.value / maxValue) * 100;
        return (
          <div key={item.name + idx} className="group">
            <div className="flex justify-between text-xs mb-1">
              <span
                className="font-medium text-foreground truncate pr-2"
                title={item.name}
              >
                {item.name}
              </span>
              <span className="text-muted-foreground tabular-nums">
                {new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                  item.value
                )}
              </span>
            </div>
            <div className="h-2 w-full bg-secondary rounded-full overflow-hidden">
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
