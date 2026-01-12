import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { HierarchyTree } from '@/widgets/charts/HierarchyTree';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DonutChart } from '@/widgets/charts/DonutChart';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { toRFC3339 } from '@/shared/lib/date';
import { analyticsApi } from '@/entities/analytics/api';
import type {
  HierarchyNode,
  AnalyticsSummary,
  GeoPoint,
} from '@/shared/api/types';

export const DashboardPage = () => {
  const { startDate, endDate } = useAnalyticsFilter();

  const [treeData, setTreeData] = useState<HierarchyNode | null>(null);
  const [geoData, setGeoData] = useState<GeoPoint[] | null>(null);
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [browserStats, setBrowserStats] = useState<Record<
    string,
    number
  > | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
        };

        const [tree, geo, sum] = await Promise.all([
          analyticsApi.getHierarchy(),
          analyticsApi.getGeoStats(params),
          analyticsApi.getSummary(params),
        ]);

        setTreeData(tree);
        setGeoData(geo);
        setSummary(sum);

        if (sum && sum.top_browsers) {
          const bStats: Record<string, number> = {};
          sum.top_browsers.forEach((b) => {
            bStats[b.name] = b.value;
          });
          setBrowserStats(bStats);
        }
      } catch (e) {
        console.error('Failed to load dashboard data', e);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [startDate, endDate]);

  return (
    <div className="space-y-6 animate-in fade-in duration-700">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
            Overview
          </h1>
          <p className="text-slate-500">General account performance.</p>
        </div>
        <DateRangePicker />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-slate-500">
              Total Clicks
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div
              className={`text-3xl font-bold ${isLoading ? 'opacity-50' : ''}`}
            >
              {new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                summary?.total_clicks || 0
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2 min-h-[400px]">
          <CardHeader>
            <CardTitle>Account Structure</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {treeData ? (
              <HierarchyTree data={treeData} />
            ) : (
              <div className="h-full flex items-center justify-center text-slate-400">
                Loading structure...
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
              <div className="h-full flex items-center justify-center text-slate-400">
                Loading stats...
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Geography</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px] w-full overflow-hidden">
          {geoData ? (
            <GeoMap data={geoData} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-400">
              Loading map...
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
