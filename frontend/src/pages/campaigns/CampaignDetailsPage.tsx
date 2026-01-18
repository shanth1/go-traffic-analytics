import { useEffect, useState, useMemo } from 'react';
import {
  useParams,
  useNavigate,
  Link as RouterLink,
  useSearchParams,
} from 'react-router-dom';
import {
  ArrowLeft,
  LinkIcon,
  BarChart2Icon,
  CalendarIcon,
  MousePointerClick,
  GlobeIcon,
  SmartphoneIcon,
  ExternalLinkIcon,
  LayersIcon,
  QrCodeIcon,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Pagination } from '@/shared/ui/pagination';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { QrCodeModal } from '@/widgets/qr/QrCodeModal';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { api } from '@/shared/api/base';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';
import { getShortLink } from '@/shared/config';

import type {
  Campaign,
  Link,
  StreamChartData,
  HeatmapPoint,
  GeoPoint,
  AnalyticsSummary,
  CategoryStat,
  CampaignsListResponse,
  LinksListResponse,
  PaginationMeta,
} from '@/shared/api/types';
import { CreateLinkFeature } from '@/features/create-link/CreateLinkFeature';

export const CampaignDetailsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const activeTab = searchParams.get('tab') || 'overview';
  const setActiveTab = (tab: string) => {
    setSearchParams({ tab });
  };

  const { startDate, endDate } = useAnalyticsFilter();
  const [loading, setLoading] = useState(true);
  const [campaign, setCampaign] = useState<Campaign | null>(null);

  const [links, setLinks] = useState<Link[]>([]);
  const [linksMeta, setLinksMeta] = useState<PaginationMeta>({
    limit: 10,
    offset: 0,
    total: 0,
  });
  const [linksLoading, setLinksLoading] = useState(false);
  const [qrSlug, setQrSlug] = useState<string | null>(null);

  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [geoData, setGeoData] = useState<GeoPoint[]>([]);
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  const topCountries = useMemo(() => {
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 6)
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

  useEffect(() => {
    if (!id) return;
    const loadAnalytics = async () => {
      setLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
          campaign_id: id,
        };
        const metaReq = api.get<CampaignsListResponse>('/campaigns', {
          params: { limit: 100 },
        });
        const analyticsReq = Promise.all([
          analyticsApi.getSummary(params),
          analyticsApi.getStream(undefined, params),
          analyticsApi.getHeatmap(undefined, params),
          analyticsApi.getGeoStats(params),
          analyticsApi.getStats('device', params),
        ]);
        const [metaRes, analyticsRes] = await Promise.all([metaReq, analyticsReq]);
        const [sum, stream, heatmap, geo, dev] = analyticsRes;

        const foundCampaign = metaRes.data.data.find((c) => c.id === id);
        if (foundCampaign) setCampaign(foundCampaign);

        setSummary(sum);
        setGeoData(geo);
        setHeatmapData(heatmap);
        setStatsDevice(dev);

        const keyTotals: Record<string, number> = {};
        stream.forEach((item) => {
          if (item.values)
            Object.entries(item.values).forEach(
              ([k, v]) => (keyTotals[k] = (keyTotals[k] || 0) + v)
            );
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
      } catch (e) {
        console.error('Failed to load campaign analytics', e);
      } finally {
        setLoading(false);
      }
    };
    loadAnalytics();
  }, [id, startDate, endDate]);

  const fetchLinks = async (offset: number) => {
    if (!id) return;
    setLinksLoading(true);
    try {
      const { data } = await api.get<LinksListResponse>('/links', {
        params: {
          campaign_id: id,
          limit: linksMeta.limit,
          offset: offset,
        },
      });
      setLinks(data.data);
      setLinksMeta(data.meta);
    } catch (e) {
      console.error('Failed to fetch links', e);
    } finally {
      setLinksLoading(false);
    }
  };

  useEffect(() => {
    fetchLinks(linksMeta.offset);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, linksMeta.offset]);

  const handlePageChange = (newOffset: number) => {
    setLinksMeta((prev) => ({ ...prev, offset: newOffset }));
  };

  if (loading && !campaign) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center gap-4 text-muted-foreground animate-pulse">
        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
        <p className="font-medium">Loading Campaign Intelligence...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 pb-20">
      {/* HEADER */}
      <div className="flex flex-col gap-6 border-b border-border pb-6">
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            size="sm"
            className="gap-2 text-muted-foreground pl-0 hover:text-primary"
            onClick={() => navigate('/campaigns')}
          >
            <ArrowLeft size={16} /> Back to Campaigns
          </Button>
        </div>

        <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
          <div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 text-primary rounded-lg">
                <LayersIcon size={24} />
              </div>
              <h1 className="text-3xl font-bold text-foreground">
                {campaign?.name || 'Campaign Details'}
              </h1>
            </div>

            <div className="flex items-center gap-4 mt-3 ml-1 text-sm text-muted-foreground">
              <span className="flex items-center gap-1.5">
                <CalendarIcon size={14} /> Created{' '}
                {campaign
                  ? new Date(campaign.created_at).toLocaleDateString()
                  : '-'}
              </span>
              <span className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-secondary text-secondary-foreground">
                <LinkIcon size={12} /> {linksMeta.total} Links
              </span>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <DateRangePicker />
          </div>
        </div>

        {/* TABS */}
        <div className="flex gap-1 bg-muted p-1.5 rounded-xl w-fit border border-border">
          <button
            onClick={() => setActiveTab('overview')}
            className={cn(
              'flex cursor-pointer items-center gap-2 px-5 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
              activeTab === 'overview'
                ? 'bg-background text-primary shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <BarChart2Icon size={16} />
            Analytics
          </button>
          <button
            onClick={() => setActiveTab('links')}
            className={cn(
              'flex cursor-pointer items-center gap-2 px-5 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
              activeTab === 'links'
                ? 'bg-background text-primary shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <LinkIcon size={16} />
            Links List
          </button>
        </div>
      </div>

      {/* CONTENT: OVERVIEW TAB */}
      {activeTab === 'overview' && (
        <div className="space-y-6 animate-in slide-in-from-bottom-4 duration-300">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <KpiCard
              title="Total Clicks"
              value={summary?.total_clicks || 0}
              icon={<MousePointerClick size={18} />}
            />
            <KpiCard
              title="Active Links"
              value={linksMeta.total}
              icon={<LinkIcon size={18} />}
              isText
            />
            <KpiCard
              title="Top Country"
              value={topCountries[0]?.name || '-'}
              icon={<GlobeIcon size={18} />}
              isText
            />
            <KpiCard
              title="Top Device"
              value={statsDevice[0]?.name || '-'}
              icon={<SmartphoneIcon size={18} />}
              isText
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Performance Trend</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {streamData.length > 0 ? (
                <StreamGraph data={streamData} keys={streamKeys} />
              ) : (
                <NoData />
              )}
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Global Reach</CardTitle>
              </CardHeader>
              <CardContent className="h-[300px] w-full overflow-hidden p-0">
                {geoData.length > 0 ? <GeoMap data={geoData} /> : <NoData />}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Activity Heatmap</CardTitle>
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
        </div>
      )}

      {/* CONTENT: LINKS TAB */}
      {activeTab === 'links' && (
        <div className="space-y-4 animate-in slide-in-from-bottom-4 duration-300">
          {linksLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-16 bg-muted rounded-lg animate-pulse"
                />
              ))}
            </div>
          ) : links.length === 0 ? (
            <div className="text-center py-20 bg-muted/20 rounded-lg border border-dashed border-border">
              <p className="text-muted-foreground mb-4">
                No links in this campaign yet.
              </p>
              <CreateLinkFeature selectedCampaignId={id} />
            </div>
          ) : (
            <>
              <div className="grid gap-3">
                {links.map((link) => (
                  <div
                    key={link.id}
                    className="flex flex-col sm:flex-row sm:items-center justify-between p-4 bg-card border border-border rounded-lg hover:border-primary/50 transition-colors gap-4"
                  >
                    <div className="flex items-start sm:items-center gap-4 overflow-hidden">
                      <div
                        className={cn(
                          'w-10 h-10 rounded-full flex items-center justify-center shrink-0 bg-secondary',
                          link.is_active ? 'text-chart-2' : 'text-destructive'
                        )}
                      >
                        <LinkIcon size={18} />
                      </div>
                      <div className="min-w-0">
                        <div className="flex items-center gap-2">
                          <h4 className="font-bold text-foreground truncate">
                            /{link.slug}
                          </h4>
                          {!link.is_active && (
                            <span className="text-[10px] bg-destructive/10 text-destructive px-1.5 rounded uppercase font-bold">
                              Inactive
                            </span>
                          )}
                        </div>
                        <p
                          className="text-xs text-muted-foreground truncate max-w-[200px] sm:max-w-[300px]"
                          title={link.target_url}
                        >
                          {link.target_url}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center justify-end gap-2 border-t sm:border-t-0 pt-3 sm:pt-0 border-border">
                      <span className="text-xs text-muted-foreground mr-2 hidden md:inline-block">
                        {new Date(link.created_at).toLocaleDateString()}
                      </span>

                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setQrSlug(link.slug)}
                        title="QR Code"
                      >
                        <QrCodeIcon size={16} />
                      </Button>

                      <RouterLink to={`/links/${link.id}`}>
                        <Button variant="secondary" size="sm" className="gap-2">
                          <BarChart2Icon size={14} /> Analytics
                        </Button>
                      </RouterLink>

                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => {
                          window.open(getShortLink(link.slug), '_blank');
                        }}
                      >
                        <ExternalLinkIcon size={16} />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="flex justify-center pt-4">
                <Pagination
                  total={linksMeta.total}
                  limit={linksMeta.limit}
                  offset={linksMeta.offset}
                  onChange={handlePageChange}
                />
              </div>
            </>
          )}
        </div>
      )}

      {qrSlug && (
        <QrCodeModal
          isOpen={!!qrSlug}
          onClose={() => setQrSlug(null)}
          slug={qrSlug}
        />
      )}
    </div>
  );
};

// --- Utils ---
interface KpiCardProps {
  title: string;
  value: string | number;
  icon: React.ReactElement;
  isText?: boolean;
}

const KpiCard = ({ title, value, icon, isText = false }: KpiCardProps) => (
  <Card>
    <CardContent className="p-6 flex flex-col justify-between h-full">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-muted-foreground">{title}</span>
        <div className="text-muted-foreground">{icon}</div>
      </div>
      <div
        className={`font-bold text-foreground ${isText ? 'text-lg truncate' : 'text-3xl'}`}
      >
        {isText
          ? value
          : new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
              value as number
            )}
      </div>
    </CardContent>
  </Card>
);

const NoData = () => (
  <div className="h-full w-full flex items-center justify-center text-muted-foreground text-sm italic">
    No data available
  </div>
);
