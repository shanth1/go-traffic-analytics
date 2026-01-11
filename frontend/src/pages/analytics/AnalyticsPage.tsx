import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { SankeyChart } from '@/widgets/charts/SankeyChart';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { toRFC3339 } from '@/shared/lib/date';
import { analyticsApi } from '@/entities/analytics/api';
import type {
  StreamChartData,
  HeatmapPoint,
  TrafficQuality,
  SankeyData,
  CategoryStat,
} from '@/shared/api/types';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const { startDate, endDate } = useAnalyticsFilter();

  // State
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [qualityData, setQualityData] = useState<TrafficQuality | null>(null);
  const [flowData, setFlowData] = useState<SankeyData>({
    nodes: [],
    links: [],
  });

  // Categorical Stats
  const [statsOS, setStatsOS] = useState<CategoryStat[]>([]);
  const [statsBrowser, setStatsBrowser] = useState<CategoryStat[]>([]);
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      setLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
        };

        // Parallel heavy fetching
        // Note: In real production, we might want to lazy load visuals below the fold
        const [
          streamRes,
          heatmapRes,
          qualityRes,
          flowRes,
          osRes,
          browserRes,
          deviceRes,
        ] = await Promise.all([
          analyticsApi.getStream(id, params),
          analyticsApi.getHeatmap(id, params),
          analyticsApi.getQuality(id, params),
          analyticsApi.getFlow(id, params),
          analyticsApi.getStats('os', { link_id: id, ...params }),
          analyticsApi.getStats('browser', { link_id: id, ...params }),
          analyticsApi.getStats('device', { link_id: id, ...params }),
        ]);

        // --- Process Stream ---
        const allKeysSet = new Set<string>();
        streamRes.forEach((item) => {
          if (item.values) {
            Object.keys(item.values).forEach((k) => allKeysSet.add(k));
          }
        });
        const collectedKeys = Array.from(allKeysSet);
        const processedStream = streamRes
          .map((d) => {
            const point: StreamChartData = { time: new Date(d.time) };
            collectedKeys.forEach((key) => {
              point[key] = d.values?.[key] ?? 0;
            });
            return point;
          })
          .sort((a, b) => a.time.getTime() - b.time.getTime());

        setStreamData(processedStream);
        setStreamKeys(collectedKeys);

        // --- Set Others ---
        setHeatmapData(heatmapRes);
        setQualityData(qualityRes);
        setFlowData(flowRes);
        setStatsOS(osRes);
        setStatsBrowser(browserRes);
        setStatsDevice(deviceRes);
      } catch (e) {
        console.error('Failed to load analytics data', e);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [id, startDate, endDate]);

  if (loading) {
    return (
      <div className="p-10 flex justify-center text-slate-500">
        Loading analytics...
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 pb-10">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-slate-50">
            Analytics Deep Dive
          </h1>
        </div>
        <DateRangePicker />
      </div>

      {/* 1. Main Traffic Stream */}
      <Card>
        <CardHeader>
          <CardTitle>Traffic Volume</CardTitle>
        </CardHeader>
        <CardContent className="h-[350px]">
          {streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* 2. Categorical Breakdown (3 cols) */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatsCard
          title="Top Devices"
          data={statsDevice}
          color="bg-indigo-500"
        />
        <StatsCard title="Top OS" data={statsOS} color="bg-emerald-500" />
        <StatsCard
          title="Top Browsers"
          data={statsBrowser}
          color="bg-amber-500"
        />
      </div>

      {/* 3. Flow & Quality Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>User Flow Journey</CardTitle>
            <p className="text-sm text-slate-400">Referer → Device → Country</p>
          </CardHeader>
          <CardContent className="h-[350px]">
            {flowData.nodes.length > 0 ? (
              <SankeyChart data={flowData} />
            ) : (
              <NoData />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Traffic Quality</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {qualityData ? <QualityRadar data={qualityData} /> : <NoData />}
          </CardContent>
        </Card>
      </div>

      {/* 4. Heatmap */}
      <Card>
        <CardHeader>
          <CardTitle>Engagement Heatmap</CardTitle>
        </CardHeader>
        <CardContent className="h-[300px]">
          {heatmapData.length > 0 ? (
            <HeatmapChart data={heatmapData} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>
    </div>
  );
};

// --- Sub-components for cleaner file ---

const StatsCard = ({
  title,
  data,
  color,
}: {
  title: string;
  data: CategoryStat[];
  color: string;
}) => (
  <Card>
    <CardHeader className="pb-2">
      <CardTitle className="text-lg">{title}</CardTitle>
    </CardHeader>
    <CardContent>
      <BarListChart data={data} color={color} />
    </CardContent>
  </Card>
);

const NoData = () => (
  <div className="h-full w-full flex items-center justify-center text-slate-400 text-sm italic">
    No data available for this period
  </div>
);
