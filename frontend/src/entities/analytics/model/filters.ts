import { create } from 'zustand';
import { startOfDay, endOfDay, subDays, isAfter } from '@/shared/lib/date';

export type DatePreset = '7d' | '30d' | '90d' | 'custom';

interface AnalyticsFilterState {
  preset: DatePreset;
  startDate: Date;
  endDate: Date;

  // Actions
  setPreset: (preset: DatePreset) => void;
  setCustomRange: (start: Date, end: Date) => void;
}

const getPresetDates = (
  preset: DatePreset,
  currentStart?: Date,
  currentEnd?: Date
): { start: Date; end: Date } => {
  const end = endOfDay(new Date());
  let start = startOfDay(new Date());

  switch (preset) {
    case '7d':
      start = startOfDay(subDays(new Date(), 7));
      break;
    case '30d':
      start = startOfDay(subDays(new Date(), 30));
      break;
    case '90d':
      start = startOfDay(subDays(new Date(), 90));
      break;
    case 'custom':
      // Preserve existing custom dates if they exist and are valid
      if (currentStart && currentEnd && isAfter(currentEnd, currentStart)) {
        return {
          start: startOfDay(currentStart),
          end: endOfDay(currentEnd),
        };
      }
      start = startOfDay(subDays(new Date(), 30));
      break;
  }
  return { start, end };
};

export const useAnalyticsFilter = create<AnalyticsFilterState>((set, get) => {
  const { start, end } = getPresetDates('30d');

  return {
    preset: '30d',
    startDate: start,
    endDate: end,

    setPreset: (preset) => {
      if (preset === 'custom') {
        set({ preset });
        return;
      }
      const { startDate, endDate } = get();
      const { start, end } = getPresetDates(preset, startDate, endDate);
      set({ preset, startDate: start, endDate: end });
    },

    setCustomRange: (start, end) => {
      if (isAfter(end, start)) {
        set({
          preset: 'custom',
          startDate: startOfDay(start),
          endDate: endOfDay(end),
        });
      }
    },
  };
});
