/**
 * Formats a Date object to RFC3339 string expected by backend.
 * Example: 2024-01-01T00:00:00Z
 */
export const toRFC3339 = (date: Date): string => {
  return date.toISOString();
};

/**
 * Returns the start of the day for a given date.
 */
export const startOfDay = (date: Date): Date => {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  return d;
};

/**
 * Returns the end of the day for a given date.
 */
export const endOfDay = (date: Date): Date => {
  const d = new Date(date);
  d.setHours(23, 59, 59, 999);
  return d;
};

/**
 * Subtracts days from a date.
 */
export const subDays = (date: Date, amount: number): Date => {
  const d = new Date(date);
  d.setDate(d.getDate() - amount);
  return d;
};

/**
 * Checks if the first date is after the second date.
 */
export const isAfter = (date: Date, compare: Date): boolean => {
  return date.getTime() > compare.getTime();
};

export const formatDateInput = (date: Date): string => {
  return date.toISOString().split('T')[0];
};

export const startOfMonth = (date: Date): string => {
  const d = new Date(date);
  d.setDate(1);
  d.setHours(0, 0, 0, 0);
  return d.toISOString();
};

export const endOfMonth = (date: Date): string => {
  const d = new Date(date);
  d.setMonth(d.getMonth() + 1);
  d.setDate(0); // Last day of previous month (which is current month relative to +1)
  d.setHours(23, 59, 59, 999);
  return d.toISOString();
};
