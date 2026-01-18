/**
 * Exports JS constants that map to CSS variables for libraries like Visx
 * that cannot consume Tailwind classes directly.
 */

export const CHART_COLORS = [
  'var(--chart-1)',
  'var(--chart-2)',
  'var(--chart-3)',
  'var(--chart-4)',
  'var(--chart-5)',
];

export const THEME_COLORS = {
  background: 'var(--background)',
  foreground: 'var(--foreground)',
  primary: 'var(--primary)',
  secondary: 'var(--secondary)',
  muted: 'var(--muted)',
  border: 'var(--border)',
  grid: 'var(--muted-foreground)', // often used for chart grids
};
