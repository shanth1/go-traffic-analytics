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
import { CampaignDetailsPage } from '@/pages/campaigns/CampaignDetailsPage';
import { LinksPage } from '@/pages/links/LinksPage';
import { AnalyticsPage } from '@/pages/analytics/AnalyticsPage';
import { ProfilePage } from '@/pages/profile/ProfilePage';
import { PricingPage } from '@/pages/pricing/PricingPage';
import { useAuthStore } from '@/entities/session/store';
import { Toaster } from '@/shared/ui/toaster';

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
      <Toaster /> {/* Global Notifications Layer */}
      <Routes>
        <Route element={<PublicRoute />}>
          <Route path="/login" element={<LoginPage />} />
        </Route>

        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/pricing" element={<PricingPage />} />
          <Route path="/campaigns" element={<CampaignsPage />} />
          <Route path="/campaigns/:id" element={<CampaignDetailsPage />} />
          <Route path="/links" element={<LinksPage />} />
          <Route path="/links/:id" element={<AnalyticsPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};
