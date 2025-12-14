import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  Outlet,
} from 'react-router-dom';

// Layouts & Widgets
import { MainLayout } from '@/widgets/layouts/MainLayout';

// Pages
import { LoginPage } from '@/pages/login/LoginPage';
import { CampaignsPage } from '@/pages/campaigns/CampaignsPage';
import { LinksPage } from '@/pages/links/LinksPage';
// Placeholders (Future Part 3)
import { ProfilePage } from '@/pages/profile/ProfilePage';
import { AnalyticsPage } from '@/pages/analytics/AnalyticsPage';

// Store
import { useAuthStore } from '@/entities/session/store';

// Компонент для защиты маршрутов
const ProtectedRoute = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <MainLayout />;
  // MainLayout внутри себя содержит <Outlet />, куда будут подставляться страницы
};

const PublicRoute = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
};

export const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<PublicRoute />}>
          <Route path="/login" element={<LoginPage />} />
        </Route>

        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<ProfilePage />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
          <Route path="/links" element={<LinksPage />} />

          <Route path="/links/:id" element={<AnalyticsPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};
