import { useEffect, useState } from 'react';

type Theme = 'dark' | 'light';

export function useTheme() {
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window !== 'undefined') {
      return document.documentElement.classList.contains('dark') ? 'dark' : 'light';
    }
    return 'light';
  });

  useEffect(() => {
    const root = window.document.documentElement;
    const metaThemeColor = document.getElementById('theme-color-meta');

    if (theme === 'dark') {
      root.classList.add('dark');
      // Set status bar to dark color (slate-950)
      metaThemeColor?.setAttribute('content', '#020617');
    } else {
      root.classList.remove('dark');
      // Set status bar to white
      metaThemeColor?.setAttribute('content', '#ffffff');
    }

    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  return { theme, toggleTheme };
}
