import { useEffect, useState, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  MousePointerClick,
  GlobeIcon,
  MonitorIcon,
  ArrowLeft,
  ExternalLinkIcon,
  LayersIcon,
  CopyIcon,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { SankeyChart } from '@/widgets/charts/SankeyChart';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { linkApi } from '@/entities/link/api';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';

import type {
  StreamChartData,
  HeatmapPoint,
  TrafficQuality,
  SankeyData,
  CategoryStat,
  GeoPoint,
  AnalyticsSummary,
  Link,
} from '@/shared/api/types';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { startDate, endDate } = useAnalyticsFilter();

  // --- Data State ---
  const [loading, setLoading] = useState(true);
  const [linkMeta, setLinkMeta] = useState<Link | null>(null);

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

  const [statsOS, setStatsOS] = useState<CategoryStat[]>([]);
  const [statsBrowser, setStatsBrowser] = useState<CategoryStat[]>([]);
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  // --- Derived Metrics ---

  const topCountries = useMemo<CategoryStat[]>(() => {
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 8)
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

  // Extract Top Referrers from Sankey Nodes (Layer 0)
  const topReferrers = useMemo<CategoryStat[]>(() => {
    if (!flowData.nodes.length) return [];

    return flowData.nodes
      .filter((n) => n.layer === 0)
      .map((n) => {
        const val = flowData.links
          .filter((l) => l.source === n.id)
          .reduce((acc, curr) => acc + curr.value, 0);
        return { name: n.id, value: val || 0, share: 0 };
      })
      .sort((a, b) => b.value - a.value)
      .slice(0, 8);
  }, [flowData]);

  // --- Fetching ---
  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      setLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
          link_id: id,
        };

        const metaPromise = linkApi.getLinkById(id);

        const analyticsPromise = Promise.all([
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

        const [meta, analyticsRes] = await Promise.all([
          metaPromise,
          analyticsPromise,
        ]);
        const [sum, stream, heatmap, quality, flow, geo, os, browser, device] =
          analyticsRes;

        if (meta) setLinkMeta(meta);

        setSummary(sum);

        const allKeysSet = new Set<string>();
        stream.forEach((item) => {
          if (item.values)
            Object.keys(item.values).forEach((k) => allKeysSet.add(k));
        });
        const collectedKeys = Array.from(allKeysSet);
        setStreamKeys(collectedKeys);
        setStreamData(
          stream
            .map((d) => {
              const p: StreamChartData = { time: new Date(d.time) };
              collectedKeys.forEach((k) => (p[k] = d.values?.[k] ?? 0));
              return p;
            })
            .sort((a, b) => a.time.getTime() - b.time.getTime())
        );

        setHeatmapData(heatmap);
        setQualityData(quality);
        setFlowData(flow);
        setGeoData(geo);
        setStatsOS(os);
        setStatsBrowser(browser);
        setStatsDevice(device);
      } catch (e) {
        console.error('Failed to load data', e);
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

  const topCountryName = topCountries[0]?.name || '-';

  return (
    <div className="space-y-8 animate-in fade-in duration-500 pb-20">
      {/* --- 1. RICH HEADER --- */}
      <div className="flex flex-col gap-4 border-b border-slate-200 dark:border-slate-800 pb-6">
        {/* Breadcrumbs & Actions */}
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            size="sm"
            className="gap-2 text-slate-500 pl-0 hover:text-indigo-600"
            onClick={() => navigate('/links')}
          >
            <ArrowLeft size={16} /> Back to Links
          </Button>

          {linkMeta && (
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                className="gap-2"
                onClick={() => {
                  navigator.clipboard.writeText(
                    window.location.host + '/' + linkMeta.slug
                  );
                }}
              >
                <CopyIcon size={14} /> Copy Short Link
              </Button>
              <Button
                variant="secondary"
                size="sm"
                className="gap-2"
                onClick={() => {
                  if (linkMeta.campaign_id)
                    navigate(`/campaigns/${linkMeta.campaign_id}`);
                  else navigate('/campaigns');
                }}
              >
                <LayersIcon size={14} /> Campaign
              </Button>
            </div>
          )}
        </div>

        {/* Title & Target */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
                /{linkMeta?.slug || id}
              </h1>
              {linkMeta && (
                <span
                  className={cn(
                    'px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider',
                    linkMeta.is_active
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-red-100 text-red-700'
                  )}
                >
                  {linkMeta.is_active ? 'Active' : 'Inactive'}
                </span>
              )}
            </div>

            {linkMeta && (
              <div className="flex items-center gap-2 text-slate-500 hover:text-indigo-600 transition-colors">
                <ExternalLinkIcon size={14} />
                <a
                  href={linkMeta.target_url}
                  target="_blank"
                  rel="noreferrer"
                  className="text-sm truncate max-w-[300px] md:max-w-[500px] underline underline-offset-4"
                >
                  {linkMeta.target_url}
                </a>
              </div>
            )}
          </div>

          <DateRangePicker />
        </div>
      </div>

      {/* --- 2. KPI --- */}
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
          title="Top Source"
          value={topReferrers[0]?.name || 'Direct'}
          icon={<ExternalLinkIcon size={18} />}
          isText
        />
      </div>

      {/* --- 3. VOLUME --- */}
      <Card>
        <CardHeader>
          <CardTitle>Traffic Dynamics</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px]">
          {streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* --- 4. REFERRERS & GEO --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* New Referrers Widget */}
        <Card>
          <CardHeader>
            <CardTitle>Top Referrers</CardTitle>
            <p className="text-sm text-slate-400">Where traffic comes from</p>
          </CardHeader>
          <CardContent>
            <BarListChart data={topReferrers} color="bg-violet-500" />
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Global Reach</CardTitle>
          </CardHeader>
          <CardContent className="h-[400px] w-full overflow-hidden p-0">
            {geoData.length > 0 ? <GeoMap data={geoData} /> : <NoData />}
          </CardContent>
        </Card>
      </div>

      {/* --- 5. TECH STACK --- */}
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

      {/* --- 6. FLOW & QUALITY --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Traffic Flow</CardTitle>
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
            <CardTitle>Quality Score</CardTitle>
          </CardHeader>
          <CardContent className="h-[350px]">
            {qualityData ? <QualityRadar data={qualityData} /> : <NoData />}
          </CardContent>
        </Card>
      </div>

      {/* --- 7. HEATMAP --- */}
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

// --- Reusable Components ---

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
  <div className="h-full w-full flex items-center justify-center text-slate-400 text-sm italic bg-slate-50/50 dark:bg-slate-900/50 rounded-lg">
    No data available
  </div>
);
