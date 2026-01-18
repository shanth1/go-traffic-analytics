export const REDIRECT_HOST =
  import.meta.env.VITE_REDIRECT_HOST || window.location.origin;

export const APP_TITLE = import.meta.env.VITE_APP_TITLE || 'Analytics';

/**
 * Helper to construct full short URL
 */
export const getShortLink = (slug: string) => {
  // Remove trailing slash from host if present
  const host = REDIRECT_HOST.replace(/\/$/, '');
  return `${host}/${slug}`;
};
