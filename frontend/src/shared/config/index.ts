export const API_URL =
  import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const REDIRECT_HOST =
  import.meta.env.VITE_REDIRECT_HOST || window.location.origin;

export const APP_TITLE = import.meta.env.VITE_APP_TITLE || 'Analytics';

export const TG_SUPPORT_URL =
  import.meta.env.VITE_TG_SUPPORT_URL || 'https://telegram.me/BotFather';

/**
 * Helper to construct full short URL
 */
export const getShortLink = (slug: string) => {
  // Remove trailing slash from host if present
  const host = REDIRECT_HOST.replace(/\/$/, '');
  return `${host}/${slug}`;
};
