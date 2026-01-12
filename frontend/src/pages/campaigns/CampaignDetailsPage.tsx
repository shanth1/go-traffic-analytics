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
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Pagination } from '@/shared/ui/pagination'; // Import Pagination
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { api } from '@/shared/api/base';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';

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

export const CampaignDetailsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  // URL-Synced Tab State
  const activeTab = searchParams.get('tab') || 'overview';
  const setActiveTab = (tab: string) => {
    setSearchParams({ tab });
  };

  const { startDate, endDate } = useAnalyticsFilter();

  // --- Global Loading State ---
  const [loading, setLoading] = useState(true);

  // --- Campaign Data ---
  const [campaign, setCampaign] = useState<Campaign | null>(null);

  // --- Links Data (Paginated) ---
  const [links, setLinks] = useState<Link[]>([]);
  const [linksMeta, setLinksMeta] = useState<PaginationMeta>({
    limit: 10,
    offset: 0,
    total: 0,
  });
  const [linksLoading, setLinksLoading] = useState(false);

  // --- Analytics Data ---
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [geoData, setGeoData] = useState<GeoPoint[]>([]);
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  // --- Helpers ---
  const topCountries = useMemo(() => {
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 6)
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

  // --- Fetch Analytics & Metadata (Depends on Date) ---
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

        // 1. Get Campaign Info
        // Note: Ideally backend has GET /campaigns/:id. Using list search as fallback.
        const metaReq = api.get<CampaignsListResponse>('/campaigns', {
          params: { limit: 100 },
        });

        // 2. Get Analytics
        const analyticsReq = Promise.all([
          analyticsApi.getSummary(params),
          analyticsApi.getStream(undefined, params),
          analyticsApi.getHeatmap(undefined, params),
          analyticsApi.getGeoStats(params),
          analyticsApi.getStats('device', params),
        ]);

        const [metaRes, analyticsRes] = await Promise.all([
          metaReq,
          analyticsReq,
        ]);
        const [sum, stream, heatmap, geo, dev] = analyticsRes;

        const foundCampaign = metaRes.data.data.find((c) => c.id === id);
        if (foundCampaign) setCampaign(foundCampaign);

        setSummary(sum);
        setGeoData(geo);
        setHeatmapData(heatmap);
        setStatsDevice(dev);

        // Process Stream
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

  // --- Fetch Links (Depends on Offset) ---
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

  // Initial Links Load or Page Change
  useEffect(() => {
    fetchLinks(linksMeta.offset);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, linksMeta.offset]);

  const handlePageChange = (newOffset: number) => {
    setLinksMeta((prev) => ({ ...prev, offset: newOffset }));
    // Scroll to top of list if needed, or just let it update
  };

  if (loading && !campaign) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center gap-4 text-slate-500 animate-pulse">
        <div className="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
        <p className="font-medium">Loading Campaign Intelligence...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 pb-20">
      {/* --- HEADER --- */}
      <div className="flex flex-col gap-6 border-b border-slate-200 dark:border-slate-800 pb-6">
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            size="sm"
            className="gap-2 text-slate-500 pl-0 hover:text-indigo-600"
            onClick={() => navigate('/campaigns')}
          >
            <ArrowLeft size={16} /> Back to Campaigns
          </Button>
        </div>

        <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
          <div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 rounded-lg">
                <LayersIcon size={24} />
              </div>
              <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
                {campaign?.name || 'Campaign Details'}
              </h1>
            </div>

            <div className="flex items-center gap-4 mt-3 ml-1 text-sm text-slate-500">
              <span className="flex items-center gap-1.5">
                <CalendarIcon size={14} /> Created{' '}
                {campaign
                  ? new Date(campaign.created_at).toLocaleDateString()
                  : '-'}
              </span>
              <span className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400">
                <LinkIcon size={12} /> {linksMeta.total} Links
              </span>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <DateRangePicker />
          </div>
        </div>

        {/* TABS */}
        <div className="flex gap-1 bg-slate-100 dark:bg-slate-900 p-1.5 rounded-xl w-fit border border-slate-200 dark:border-slate-800">
          <button
            onClick={() => setActiveTab('overview')}
            className={cn(
              'flex items-center gap-2 px-5 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
              activeTab === 'overview'
                ? 'bg-white dark:bg-slate-800 text-indigo-600 shadow-sm'
                : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
            )}
          >
            <BarChart2Icon size={16} />
            Analytics
          </button>
          <button
            onClick={() => setActiveTab('links')}
            className={cn(
              'flex items-center gap-2 px-5 py-2 text-sm font-semibold rounded-lg transition-all duration-200',
              activeTab === 'links'
                ? 'bg-white dark:bg-slate-800 text-indigo-600 shadow-sm'
                : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
            )}
          >
            <LinkIcon size={16} />
            Links List
          </button>
        </div>
      </div>

      {/* --- CONTENT: OVERVIEW TAB --- */}
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

      {/* --- CONTENT: LINKS TAB --- */}
      {activeTab === 'links' && (
        <div className="space-y-4 animate-in slide-in-from-bottom-4 duration-300">
          {linksLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-16 bg-slate-100 dark:bg-slate-900 rounded-lg animate-pulse"
                />
              ))}
            </div>
          ) : links.length === 0 ? (
            <div className="text-center py-20 bg-slate-50 dark:bg-slate-900/50 rounded-lg border border-dashed border-slate-200 dark:border-slate-800">
              <p className="text-slate-500 mb-4">
                No links in this campaign yet.
              </p>
              <Button onClick={() => navigate('/links')}>Create Link</Button>
            </div>
          ) : (
            <>
              <div className="grid gap-3">
                {links.map((link) => (
                  <div
                    key={link.id}
                    className="flex flex-col sm:flex-row sm:items-center justify-between p-4 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-lg hover:border-indigo-300 transition-colors gap-4"
                  >
                    <div className="flex items-start sm:items-center gap-4 overflow-hidden">
                      <div
                        className={cn(
                          'w-10 h-10 rounded-full flex items-center justify-center shrink-0 bg-slate-100 dark:bg-slate-900',
                          link.is_active ? 'text-green-600' : 'text-red-500'
                        )}
                      >
                        <LinkIcon size={18} />
                      </div>
                      <div className="min-w-0">
                        <div className="flex items-center gap-2">
                          <h4 className="font-bold text-slate-900 dark:text-slate-100 truncate">
                            /{link.slug}
                          </h4>
                          {!link.is_active && (
                            <span className="text-[10px] bg-red-100 text-red-600 px-1.5 rounded uppercase font-bold">
                              Inactive
                            </span>
                          )}
                        </div>
                        <p
                          className="text-xs text-slate-500 truncate max-w-[200px] sm:max-w-[300px]"
                          title={link.target_url}
                        >
                          {link.target_url}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center justify-end gap-2 border-t sm:border-t-0 pt-3 sm:pt-0 border-slate-100 dark:border-slate-800">
                      <span className="text-xs text-slate-400 mr-2 hidden md:inline-block">
                        {new Date(link.created_at).toLocaleDateString()}
                      </span>

                      <RouterLink to={`/links/${link.id}`}>
                        <Button variant="secondary" size="sm" className="gap-2">
                          <BarChart2Icon size={14} /> Analytics
                        </Button>
                      </RouterLink>

                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => {
                          window.open(
                            `http://${window.location.host}/${link.slug}`,
                            '_blank'
                          );
                        }}
                      >
                        <ExternalLinkIcon size={16} />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              {/* PAGINATION */}
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
        <span className="text-sm font-medium text-slate-500">{title}</span>
        <div className="text-slate-400">{icon}</div>
      </div>
      <div
        className={`font-bold text-slate-900 dark:text-slate-100 ${isText ? 'text-lg truncate' : 'text-3xl'}`}
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
  <div className="h-full w-full flex items-center justify-center text-slate-400 text-sm italic">
    No data available
  </div>
);
