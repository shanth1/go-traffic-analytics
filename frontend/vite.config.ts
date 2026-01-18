import path from 'path';
import react from '@vitejs/plugin-react';
import { defineConfig, loadEnv } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  const appTitle = env.VITE_APP_TITLE || 'Analytics';

  return {
    plugins: [
      react(),
      tailwindcss(),
      {
        name: 'html-title-replace',
        transformIndexHtml(html) {
          return html.replace(/%VITE_APP_TITLE%/g, appTitle);
        },
      },
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  };
});
