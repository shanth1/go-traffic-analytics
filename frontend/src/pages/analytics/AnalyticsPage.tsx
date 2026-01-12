import { useEffect, useState, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import {
  MousePointerClick,
  GlobeIcon,
  SmartphoneIcon,
  MonitorIcon,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { SankeyChart } from '@/widgets/charts/SankeyChart';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
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
  GeoPoint,
  AnalyticsSummary,
} from '@/shared/api/types';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const { startDate, endDate } = useAnalyticsFilter();

  // --- State ---
  const [loading, setLoading] = useState(true);

  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [qualityData, setQualityData] = useState<TrafficQuality | null>(null);
  const [flowData, setFlowData] = useState<SankeyData>({
    nodes: [],
    links: [],
  });
  const [geoData, setGeoData] = useState<GeoPoint[]>([]);

  // Categorical Stats
  const [statsOS, setStatsOS] = useState<CategoryStat[]>([]);
  const [statsBrowser, setStatsBrowser] = useState<CategoryStat[]>([]);
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  // --- Derived State for Top Lists ---
  const topCountries = useMemo<CategoryStat[]>(() => {
    // Convert GeoPoints to CategoryStat for the BarList component
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 8) // Top 8
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

  const topCountryName = topCountries[0]?.name || '-';

  // --- Data Fetching ---
  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      setLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
          link_id: id, // Ensure we filter by this link
        };

        const [
          sumRes,
          streamRes,
          heatmapRes,
          qualityRes,
          flowRes,
          geoRes,
          osRes,
          browserRes,
          deviceRes,
        ] = await Promise.all([
          analyticsApi.getSummary(params),
          analyticsApi.getStream(id, params),
          analyticsApi.getHeatmap(id, params),
          analyticsApi.getQuality(id, params),
          analyticsApi.getFlow(id, params),
          analyticsApi.getGeoStats(params),
          analyticsApi.getStats('os', params),
          analyticsApi.getStats('browser', params),
          analyticsApi.getStats('device', params),
        ]);

        // Process Summary
        setSummary(sumRes);

        // Process Stream
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

        setHeatmapData(heatmapRes);
        setQualityData(qualityRes);
        setFlowData(flowRes);
        setGeoData(geoRes);
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
      <div className="min-h-[60vh] flex flex-col items-center justify-center gap-4 text-slate-500 animate-pulse">
        <div className="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
        <p className="font-medium">Aggregating analytics data...</p>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500 pb-10">
      {/* 1. Page Header & Filters */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-slate-200 dark:border-slate-800 pb-6">
        <div>
          <div className="flex items-center gap-2 text-slate-500 mb-1">
            <span className="text-xs font-bold bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400 px-2 py-0.5 rounded uppercase tracking-wider">
              Link Analytics
            </span>
          </div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
            /{id}
          </h1>
          <p className="text-slate-500 text-sm mt-1">
            Detailed performance report.
          </p>
        </div>
        <DateRangePicker />
      </div>

      {/* 2. KPI Cards Row */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <KpiCard
          title="Total Clicks"
          value={summary?.total_clicks || 0}
          icon={<MousePointerClick size={18} />}
        />
        <KpiCard
          title="Top Country"
          value={topCountryName}
          icon={<GlobeIcon size={18} />}
          isText
        />
        <KpiCard
          title="Top OS"
          value={statsOS[0]?.name || '-'}
          icon={<MonitorIcon size={18} />}
          isText
        />
        <KpiCard
          title="Top Device"
          value={statsDevice[0]?.name || '-'}
          icon={<SmartphoneIcon size={18} />}
          isText
        />
      </div>

      {/* 3. Main Stream Graph (Volume) */}
      <Card>
        <CardHeader>
          <CardTitle>Traffic Volume & Dynamics</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px]">
          {streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* 4. Geography Section (Map + List) */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Global Reach</CardTitle>
            <p className="text-sm text-slate-400">
              Click distribution by country
            </p>
          </CardHeader>
          <CardContent className="h-[400px] w-full overflow-hidden p-0">
            {geoData.length > 0 ? <GeoMap data={geoData} /> : <NoData />}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top Countries</CardTitle>
          </CardHeader>
          <CardContent>
            <BarListChart data={topCountries} color="bg-blue-500" />
          </CardContent>
        </Card>
      </div>

      {/* 5. Tech Stack (3 cols) */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatsCard
          title="Device Type"
          data={statsDevice}
          color="bg-indigo-500"
        />
        <StatsCard
          title="Operating System"
          data={statsOS}
          color="bg-emerald-500"
        />
        <StatsCard title="Browser" data={statsBrowser} color="bg-amber-500" />
      </div>

      {/* 6. Behavior Section (Flow + Quality) */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>User Journey Flow</CardTitle>
            <p className="text-sm text-slate-400">Source → Device → Country</p>
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
            <CardTitle>Traffic Quality Score</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {qualityData ? <QualityRadar data={qualityData} /> : <NoData />}
          </CardContent>
        </Card>
      </div>

      {/* 7. Timing (Heatmap) */}
      <Card>
        <CardHeader>
          <CardTitle>Engagement Heatmap</CardTitle>
          <p className="text-sm text-slate-400">Best performing times (UTC)</p>
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

// --- Sub-components ---

const KpiCard = ({
  title,
  value,
  icon,
  isText = false,
}: {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  isText?: boolean;
}) => (
  <Card>
    <CardContent className="p-6 flex flex-col justify-between h-full">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-slate-500">{title}</span>
        <div className="text-slate-400">{icon}</div>
      </div>
      <div
        className={`font-bold text-slate-900 dark:text-slate-100 ${isText ? 'text-lg truncate' : 'text-3xl'}`}
        title={String(value)}
      >
        {isText
          ? value
          : new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
              Number(value)
            )}
      </div>
    </CardContent>
  </Card>
);

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
  <div className="h-full w-full flex items-center justify-center text-slate-400 text-sm italic bg-slate-50/50 dark:bg-slate-900/50 rounded-lg border border-dashed border-slate-200 dark:border-slate-800">
    No data available for this period
  </div>
);
