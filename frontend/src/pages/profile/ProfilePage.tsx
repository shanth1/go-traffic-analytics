import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { HierarchyTree } from '@/widgets/charts/HierarchyTree';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DonutChart } from '@/widgets/charts/DonutChart';
import type { HierarchyNode } from '@/shared/api/types';
import { useAuthStore } from '@/entities/session/store';

// Типы для стейта
type BrowserStats = Record<string, number>;
type GeoStats = Record<string, number>;

export const ProfilePage = () => {
  const user = useAuthStore((s) => s.user);

  // Указываем конкретные типы вместо any
  const [treeData, setTreeData] = useState<HierarchyNode | null>(null);
  const [geoData, setGeoData] = useState<GeoStats | null>(null);
  const [browserStats, setBrowserStats] = useState<BrowserStats | null>(null);

  useEffect(() => {
    // Используем setTimeout чтобы избежать синхронного обновления внутри эффекта
    // Это удовлетворяет линтер и имитирует сетевую задержку
    const timer = setTimeout(() => {
      // 1. Mock Tree Data
      setTreeData({
        name: user?.email || 'User',
        type: 'root',
        children: [
          {
            name: 'Summer Sale',
            type: 'campaign',
            children: [
              { name: '/sale', type: 'link' },
              { name: '/bonus', type: 'link' },
            ],
          },
          {
            name: 'Socials',
            type: 'campaign',
            children: [
              { name: '/insta', type: 'link' },
              { name: '/yt', type: 'link' },
            ],
          },
        ],
      });

      // 2. Mock Geo Data
      setGeoData({
        'United States of America': 500,
        Russia: 300,
        Germany: 150,
        Brazil: 80,
      });

      // 3. Mock Stats
      setBrowserStats({ Chrome: 65, Safari: 20, Firefox: 10, Edge: 5 });
    }, 100);

    return () => clearTimeout(timer);
  }, [user]);

  return (
    <div className="space-y-6 animate-in fade-in duration-700">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-slate-500">
              Клики за месяц
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {user?.clicks_current_month || 0}
            </div>
          </CardContent>
        </Card>
        {/* ... другие карточки (код тот же) ... */}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2 min-h-[400px]">
          <CardHeader>
            <CardTitle>Структура кампаний</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {treeData ? (
              <HierarchyTree data={treeData} />
            ) : (
              <div className="h-full flex items-center justify-center">
                Loading...
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="min-h-[400px]">
          <CardHeader>
            <CardTitle>Топ браузеров</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {browserStats ? (
              <DonutChart data={browserStats} />
            ) : (
              <div>Loading...</div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>География трафика</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px] w-full overflow-hidden">
          {geoData ? (
            <GeoMap data={geoData} />
          ) : (
            <div className="p-4">Loading Map...</div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
