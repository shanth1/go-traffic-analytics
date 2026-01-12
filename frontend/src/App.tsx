import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  Outlet,
} from 'react-router-dom';
import { MainLayout } from '@/widgets/layouts/MainLayout';
import { LoginPage } from '@/pages/login/LoginPage';
import { DashboardPage } from '@/pages/dashboard/DashboardPage';
import { CampaignsPage } from '@/pages/campaigns/CampaignsPage';
import { LinksPage } from '@/pages/links/LinksPage';
import { ProfilePage } from '@/pages/profile/ProfilePage';
import { AnalyticsPage } from '@/pages/analytics/AnalyticsPage';
import { useAuthStore } from '@/entities/session/store';

const ProtectedRoute = () => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <MainLayout />;
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
          <Route path="/" element={<DashboardPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
          <Route path="/links" element={<LinksPage />} />

          <Route path="/links/:id" element={<AnalyticsPage />} />
          {/* Campaign Analytics will go here later: /campaigns/:id */}
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};
