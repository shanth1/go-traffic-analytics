import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { HierarchyTree } from '@/widgets/charts/HierarchyTree';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DonutChart } from '@/widgets/charts/DonutChart';
import type {
  HierarchyNode,
  AnalyticsSummary,
  GeoPoint,
} from '@/shared/api/types';
import { analyticsApi } from '@/entities/analytics/api';

export const ProfilePage = () => {
  const [treeData, setTreeData] = useState<HierarchyNode | null>(null);
  const [geoData, setGeoData] = useState<GeoPoint[] | null>(null);
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [browserStats, setBrowserStats] = useState<Record<
    string,
    number
  > | null>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const tree = await analyticsApi.getHierarchy();
        setTreeData(tree);

        const geo = await analyticsApi.getGeoStats();
        setGeoData(geo);

        const sum = await analyticsApi.getSummary();
        setSummary(sum);

        if (sum && sum.top_browsers) {
          const bStats: Record<string, number> = {};
          sum.top_browsers.forEach((b) => {
            bStats[b.name] = b.value;
          });
          setBrowserStats(bStats);
        }
      } catch (e) {
        console.error('Failed to load profile data', e);
      }
    };

    loadData();
  }, []);

  return (
    <div className="space-y-6 animate-in fade-in duration-700">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-slate-500">
              Total Clicks
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {summary?.total_clicks || 0}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2 min-h-[400px]">
          <CardHeader>
            <CardTitle>Campaign Structure</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {treeData ? (
              <HierarchyTree data={treeData} />
            ) : (
              <div className="h-full flex items-center justify-center">
                Loading campaign structure...
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="min-h-[400px]">
          <CardHeader>
            <CardTitle>Top Browsers</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {browserStats ? (
              <DonutChart data={browserStats} />
            ) : (
              <div className="h-full flex items-center justify-center">
                Loading browser stats...
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Traffic Geography</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px] w-full overflow-hidden">
          {geoData ? (
            <GeoMap data={geoData} />
          ) : (
            <div className="p-4 h-full flex items-center justify-center">
              Loading map...
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
