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
  QrCodeIcon,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Skeleton } from '@/shared/ui/skeleton'; // New
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { SankeyChart } from '@/widgets/charts/SankeyChart';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { QrCodeModal } from '@/widgets/qr/QrCodeModal';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { linkApi } from '@/entities/link/api';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';
import { getShortLink } from '@/shared/config';

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
import { toast } from '@/entities/notification/store';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { startDate, endDate } = useAnalyticsFilter();

  const [loading, setLoading] = useState(true);
  const [metaLoading, setMetaLoading] = useState(true);
  const [linkMeta, setLinkMeta] = useState<Link | null>(null);
  const [isQrOpen, setIsQrOpen] = useState(false);

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

  const topCountries = useMemo<CategoryStat[]>(() => {
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 8)
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

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

  // Load Meta separately to show header fast
  useEffect(() => {
    if (!id) return;
    const loadMeta = async () => {
      setMetaLoading(true);
      try {
        const meta = await linkApi.getLinkById(id);
        if (meta) setLinkMeta(meta);
      } catch (e) {
        console.error(e);
      } finally {
        setMetaLoading(false);
      }
    };
    loadMeta();
  }, [id]);

  // Load Heavy Analytics
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

        const [sum, stream, heatmap, quality, flow, geo, os, browser, device] =
          await Promise.all([
            analyticsApi.getSummary(params),
            analyticsApi.getStream(id, params),
            analyticsApi.getHeatmap(id, params),
            analyticsApi.getQuality(id, params),
            analyticsApi.getFlow(id, params),
            analyticsApi.getCountryStats(params),
            analyticsApi.getStats('os', params),
            analyticsApi.getStats('browser', params),
            analyticsApi.getStats('device', params),
          ]);

        setSummary(sum);

        const keyTotals: Record<string, number> = {};
        stream.forEach((item) => {
          if (item.values) {
            Object.entries(item.values).forEach(([key, val]) => {
              keyTotals[key] = (keyTotals[key] || 0) + val;
            });
          }
        });

        const sortedKeys = Object.keys(keyTotals).sort(
          (a, b) => keyTotals[b] - keyTotals[a]
        );
        setStreamKeys(sortedKeys);

        setStreamData(
          stream
            .map((d) => {
              const p: StreamChartData = { time: new Date(d.time) };
              sortedKeys.forEach((k) => (p[k] = d.values?.[k] ?? 0));
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

  const topCountryName = topCountries[0]?.name || '-';

  return (
    <div className="space-y-8 animate-in fade-in duration-500 pb-20">
      {/* --- 1. RICH HEADER --- */}
      <div className="flex flex-col gap-4 border-b border-border pb-6">
        {/* Top Row: Back Button & Actions */}
        <div className="flex flex-wrap items-center justify-between gap-4">
          <Button
            variant="ghost"
            size="sm"
            className="gap-2 text-muted-foreground pl-0 hover:text-primary shrink-0"
            onClick={() => navigate('/links')}
          >
            <ArrowLeft size={16} /> Back to Links
          </Button>

          <div className="flex items-center gap-2 ml-auto">
            {metaLoading ? (
              <Skeleton className="h-9 w-32" />
            ) : (
              linkMeta && (
                <>
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-2 px-3"
                    onClick={() => setIsQrOpen(true)}
                    title="Show QR Code"
                  >
                    <QrCodeIcon size={16} />
                    <span className="hidden sm:inline">QR Code</span>
                  </Button>

                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-2 px-3"
                    onClick={() => {
                      navigator.clipboard.writeText(
                        getShortLink(linkMeta.slug)
                      );
                      toast.info(
                        'Copied to clipboard',
                        getShortLink(linkMeta.slug)
                      );
                    }}
                    title="Copy Link"
                  >
                    <CopyIcon size={16} />
                    <span className="hidden sm:inline">Copy</span>
                  </Button>

                  <Button
                    variant="secondary"
                    size="sm"
                    className="gap-2 px-3"
                    onClick={() => {
                      if (linkMeta.campaign_id)
                        navigate(`/campaigns/${linkMeta.campaign_id}`);
                      else navigate('/campaigns');
                    }}
                    title="View Campaign"
                  >
                    <LayersIcon size={16} />
                    <span className="hidden sm:inline">Campaign</span>
                  </Button>
                </>
              )
            )}
          </div>
        </div>

        {/* Title & Filters Row */}
        <div className="flex flex-col xl:flex-row justify-between items-start xl:items-end gap-4">
          <div className="w-full xl:w-auto">
            {metaLoading ? (
              <div className="space-y-2">
                <Skeleton className="h-10 w-full sm:w-64" />
                <Skeleton className="h-4 w-full sm:w-96" />
              </div>
            ) : (
              <>
                <div className="flex flex-wrap items-center gap-3 mb-2">
                  <h1 className="text-2xl sm:text-3xl font-bold text-foreground break-all">
                    /{linkMeta?.slug || id}
                  </h1>
                  {linkMeta && (
                    <span
                      className={cn(
                        'px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider shrink-0',
                        linkMeta.is_active
                          ? 'bg-chart-2/10 text-chart-2'
                          : 'bg-destructive/10 text-destructive'
                      )}
                    >
                      {linkMeta.is_active ? 'Active' : 'Inactive'}
                    </span>
                  )}
                </div>

                {linkMeta && (
                  <div className="flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors overflow-hidden">
                    <ExternalLinkIcon size={14} className="shrink-0" />
                    <a
                      href={linkMeta.target_url}
                      target="_blank"
                      rel="noreferrer"
                      className="text-sm truncate underline underline-offset-4 block w-full"
                    >
                      {linkMeta.target_url}
                    </a>
                  </div>
                )}
              </>
            )}
          </div>

          <div className="w-full xl:w-auto">
            <DateRangePicker />
          </div>
        </div>
      </div>

      {/* --- 2. KPI --- */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <KpiCard
          title="Total Clicks"
          value={summary?.total_clicks || 0}
          icon={<MousePointerClick size={18} />}
          loading={loading}
        />
        <KpiCard
          title="Top Country"
          value={topCountryName}
          icon={<GlobeIcon size={18} />}
          isText
          loading={loading}
        />
        <KpiCard
          title="Top OS"
          value={statsOS[0]?.name || '-'}
          icon={<MonitorIcon size={18} />}
          isText
          loading={loading}
        />
        <KpiCard
          title="Top Source"
          value={(() => {
            const v = topReferrers[0]?.name || 'Direct';
            const idx = v.indexOf(':');
            return idx !== -1 ? v.slice(idx + 1).trim() : v;
          })()}
          icon={<ExternalLinkIcon size={18} />}
          isText
          loading={loading}
        />
      </div>

      {/* --- 3. VOLUME --- */}
      <Card>
        <CardHeader>
          <CardTitle>Traffic Dynamics</CardTitle>
        </CardHeader>
        <CardContent className="h-[400px]">
          {loading ? (
            <Skeleton className="w-full h-full" />
          ) : streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* --- 4. REFERRERS & GEO --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Top Referrers</CardTitle>
            <p className="text-sm text-muted-foreground">
              Where traffic comes from
            </p>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3, 4].map((i) => (
                  <Skeleton key={i} className="h-6 w-full" />
                ))}
              </div>
            ) : (
              <BarListChart data={topReferrers} color="bg-chart-5" />
            )}
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Global Reach</CardTitle>
          </CardHeader>
          <CardContent className="h-[400px] w-full overflow-hidden p-0">
            {loading ? (
              <div className="p-6 h-full">
                <Skeleton className="w-full h-full" />
              </div>
            ) : geoData.length > 0 ? (
              <GeoMap data={geoData} />
            ) : (
              <NoData />
            )}
          </CardContent>
        </Card>
      </div>

      {/* --- 5. TECH STACK --- */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatsCard
          title="Device Type"
          data={statsDevice}
          color="bg-chart-1"
          loading={loading}
        />
        <StatsCard
          title="Operating System"
          data={statsOS}
          color="bg-chart-2"
          loading={loading}
        />
        <StatsCard
          title="Browser"
          data={statsBrowser}
          color="bg-chart-3"
          loading={loading}
        />
      </div>

      {/* --- 6. FLOW & QUALITY --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Traffic Flow</CardTitle>
            <p className="text-sm text-muted-foreground">
              Referer → Device → Country
            </p>
          </CardHeader>
          <CardContent className="h-[350px]">
            {loading ? (
              <Skeleton className="w-full h-full" />
            ) : flowData.nodes.length > 0 ? (
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
            {loading ? (
              <Skeleton className="w-full h-full rounded-full" />
            ) : qualityData ? (
              <QualityRadar data={qualityData} />
            ) : (
              <NoData />
            )}
          </CardContent>
        </Card>
      </div>

      {/* --- 7. HEATMAP --- */}
      <Card>
        <CardHeader>
          <CardTitle>Engagement Heatmap</CardTitle>
        </CardHeader>
        <CardContent className="h-[300px]">
          {loading ? (
            <Skeleton className="w-full h-full" />
          ) : heatmapData.length > 0 ? (
            <HeatmapChart data={heatmapData} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* QR Modal */}
      {linkMeta && (
        <QrCodeModal
          isOpen={isQrOpen}
          onClose={() => setIsQrOpen(false)}
          slug={linkMeta.slug}
        />
      )}
    </div>
  );
};

const KpiCard = ({
  title,
  value,
  icon,
  isText = false,
  loading = false,
}: {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  isText?: boolean;
  loading?: boolean;
}) => (
  <Card>
    <CardContent className="p-6 flex flex-col justify-between h-full">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-muted-foreground">
          {title}
        </span>
        <div className="text-muted-foreground">{icon}</div>
      </div>
      {loading ? (
        <Skeleton className="h-8 w-24" />
      ) : (
        <div
          className={`font-bold text-foreground ${isText ? 'text-lg truncate' : 'text-3xl'}`}
          title={String(value)}
        >
          {isText
            ? value
            : new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                Number(value)
              )}
        </div>
      )}
    </CardContent>
  </Card>
);

const StatsCard = ({
  title,
  data,
  color,
  loading,
}: {
  title: string;
  data: CategoryStat[];
  color: string;
  loading: boolean;
}) => (
  <Card>
    <CardHeader className="pb-2">
      <CardTitle className="text-lg">{title}</CardTitle>
    </CardHeader>
    <CardContent>
      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-6 w-full" />
          ))}
        </div>
      ) : (
        <BarListChart data={data} color={color} />
      )}
    </CardContent>
  </Card>
);

const NoData = () => (
  <div className="h-full w-full flex items-center justify-center text-muted-foreground text-sm italic bg-muted/20 rounded-lg">
    No data available
  </div>
);
