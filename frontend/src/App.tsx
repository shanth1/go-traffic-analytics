import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from '@/app/Layout';
import { CampaignsPage } from '@/pages/CampaignsPage';
import { AnalyticsPage } from '@/pages/AnalyticsPage';
import { ProfilePage } from '@/pages/ProfilePage';
import { CreatePage } from '@/pages/CreatePage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="campaigns" element={<CampaignsPage />} />
          <Route path="analytics" element={<AnalyticsPage />} />
          <Route path="profile" element={<ProfilePage />} />
          <Route path="create" element={<CreatePage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
