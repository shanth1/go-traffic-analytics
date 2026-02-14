import { useEffect, useState, useMemo } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import {
  ArrowLeft,
  LinkIcon,
  BarChart2Icon,
  CalendarIcon,
  MousePointerClick,
  GlobeIcon,
  SmartphoneIcon,
  LayersIcon,
  MapPinIcon,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Pagination } from '@/shared/ui/pagination';
import { Skeleton } from '@/shared/ui/skeleton';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { ConfirmationModal } from '@/shared/ui/confirmation-modal';
import { LinkCard } from '@/entities/link/ui/LinkCard';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { api } from '@/shared/api/base';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';
import { useTranslation } from 'react-i18next';

import type {
  Campaign,
  Link,
  StreamChartData,
  HeatmapPoint,
  GeoStats,
  AnalyticsSummary,
  CategoryStat,
  CampaignsListResponse,
  LinksListResponse,
  PaginationMeta,
} from '@/shared/api/types';
import { CreateLinkFeature } from '@/features/create-link/CreateLinkFeature';
import { useLinkStore } from '@/entities/link/model/store';
import { toast } from '@/entities/notification/store';

export const CampaignDetailsPage = () => {
  const { t } = useTranslation();
  const { id } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const deleteLinkStore = useLinkStore((s) => s.deleteLink);

  const activeTab = searchParams.get('tab') || 'overview';
  const setActiveTab = (tab: string) => {
    setSearchParams({ tab });
  };

  const { startDate, endDate } = useAnalyticsFilter();

  const [analyticsLoading, setAnalyticsLoading] = useState(true);
  const [metaLoading, setMetaLoading] = useState(true);

  const [campaign, setCampaign] = useState<Campaign | null>(null);

  const [links, setLinks] = useState<Link[]>([]);
  const [linksMeta, setLinksMeta] = useState<PaginationMeta>({
    limit: 10,
    offset: 0,
    total: 0,
  });
  const [linksLoading, setLinksLoading] = useState(false);

  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [geoStats, setGeoStats] = useState<GeoStats>({
    countries: [],
    cities: [],
  });
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);

  const topCountries = useMemo(() => {
    return geoStats.countries
      .sort((a, b) => b.value - a.value)
      .slice(0, 6)
      .map((g) => ({ name: g.country, value: g.value }));
  }, [geoStats.countries]);

  const topCities = useMemo(() => {
    return geoStats.cities
      .sort((a, b) => b.value - a.value)
      .slice(0, 6)
      .map((c) => ({ name: c.city, value: c.value }));
  }, [geoStats.cities]);

  const topCityName = topCities[0]?.name || '-';

  useEffect(() => {
    if (!id) return;
    const loadMeta = async () => {
      setMetaLoading(true);
      try {
        const { data } = await api.get<CampaignsListResponse>('/campaigns', {
          params: { limit: 100 },
        });
        const found = data.data.find((c) => c.id === id);
        if (found) setCampaign(found);
      } catch (e) {
        console.error(e);
      } finally {
        setMetaLoading(false);
      }
    };
    loadMeta();
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const loadAnalytics = async () => {
      setAnalyticsLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
          campaign_id: id,
        };

        const [sum, stream, heatmap, geo, dev] = await Promise.all([
          analyticsApi.getSummary(params),
          analyticsApi.getStream(undefined, params),
          analyticsApi.getHeatmap(undefined, params),
          analyticsApi.getGeoStats(params),
          analyticsApi.getStats('device', params),
        ]);

        setSummary(sum);
        setGeoStats(geo);
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
        setAnalyticsLoading(false);
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

  const handleDeleteConfirm = async () => {
    if (!deleteId) return;
    setIsDeleting(true);
    try {
      await deleteLinkStore(deleteId);
      toast.info(
        'Link Deleted',
        'The short link has been permanently removed.'
      );
      setDeleteId(null);
      fetchLinks(linksMeta.offset);
    } catch (e) {
      console.error('Failed to delete link', e);
    } finally {
      setIsDeleting(false);
    }
  };

  const handleCreateSuccess = () => {
    setLinksMeta((prev) => ({ ...prev, offset: 0 }));
    fetchLinks(0);
  };

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
            <ArrowLeft size={16} /> {t('campaigns.details.back_to_list')}
          </Button>
        </div>

        <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
          <div>
            {metaLoading ? (
              <div className="space-y-2">
                <Skeleton className="h-10 w-64" />
                <Skeleton className="h-4 w-48" />
              </div>
            ) : (
              <>
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
                    <CalendarIcon size={14} /> {t('common.created')}{' '}
                    {campaign
                      ? new Date(campaign.created_at).toLocaleDateString()
                      : '-'}
                  </span>
                  <span className="flex items-center gap-1.5 px-2 py-0.5 rounded bg-secondary text-secondary-foreground font-medium">
                    <LinkIcon size={12} />{' '}
                    {t('campaigns.details.links_count', {
                      count: linksMeta.total,
                    })}
                  </span>
                </div>
              </>
            )}
          </div>

          <div className="flex flex-col sm:flex-row items-end gap-3 w-full sm:w-auto">
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
            {t('campaigns.details.tabs.analytics')}
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
            {t('campaigns.details.tabs.links_list')}
          </button>
        </div>
      </div>

      {activeTab === 'overview' && (
        <div className="space-y-6 animate-in slide-in-from-bottom-4 duration-300">
          {/* KPI ROW */}
          <div className="grid grid-cols-2 lg:grid-cols-5 gap-4">
            <KpiCard
              title={t('campaigns.details.kpi.total_clicks')}
              value={summary?.total_clicks || 0}
              icon={<MousePointerClick size={18} />}
              loading={analyticsLoading}
            />
            <KpiCard
              title={t('campaigns.details.kpi.active_links')}
              value={linksMeta.total}
              icon={<LinkIcon size={18} />}
              isText
              loading={metaLoading}
            />
            <KpiCard
              title={t('campaigns.details.kpi.top_country')}
              value={topCountries[0]?.name || '-'}
              icon={<GlobeIcon size={18} />}
              isText
              loading={analyticsLoading}
            />
            <KpiCard
              title={t('dashboard.kpi.top_city')}
              value={topCityName}
              icon={<MapPinIcon size={18} />}
              isText
              loading={analyticsLoading}
            />
            <KpiCard
              title={t('campaigns.details.kpi.top_device')}
              value={statsDevice[0]?.name || '-'}
              icon={<SmartphoneIcon size={18} />}
              isText
              loading={analyticsLoading}
            />
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Performance Trend</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {analyticsLoading ? (
                <Skeleton className="w-full h-full" />
              ) : streamData.length > 0 ? (
                <StreamGraph data={streamData} keys={streamKeys} />
              ) : (
                <NoData />
              )}
            </CardContent>
          </Card>

          {/* GEO SECTION: DASHBOARD STYLE (3 Separate Cards) */}
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-6">
            {/* 1. MAP (Takes 2 cols) */}
            <Card className="lg:col-span-2">
              <CardHeader>
                <CardTitle>{t('analytics.charts.global_reach')}</CardTitle>
              </CardHeader>
              <CardContent className="h-[400px] w-full p-0 overflow-hidden">
                {analyticsLoading ? (
                  <div className="p-6 h-full">
                    <Skeleton className="w-full h-full" />
                  </div>
                ) : geoStats.countries.length > 0 ? (
                  <GeoMap data={geoStats.countries} />
                ) : (
                  <NoData />
                )}
              </CardContent>
            </Card>

            {/* 2. TOP COUNTRIES */}
            <Card>
              <CardHeader>
                <CardTitle>{t('dashboard.charts.top_countries')}</CardTitle>
              </CardHeader>
              <CardContent>
                {analyticsLoading ? (
                  <div className="space-y-3">
                    {[1, 2, 3, 4, 5].map((i) => (
                      <Skeleton key={i} className="h-6 w-full" />
                    ))}
                  </div>
                ) : (
                  <BarListChart data={topCountries} color="bg-primary" />
                )}
              </CardContent>
            </Card>

            {/* 3. TOP CITIES */}
            <Card>
              <CardHeader>
                <CardTitle>{t('dashboard.charts.top_cities')}</CardTitle>
              </CardHeader>
              <CardContent>
                {analyticsLoading ? (
                  <div className="space-y-3">
                    {[1, 2, 3, 4, 5].map((i) => (
                      <Skeleton key={i} className="h-6 w-full" />
                    ))}
                  </div>
                ) : (
                  <BarListChart data={topCities} color="bg-chart-5" />
                )}
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>{t('analytics.charts.engagement_heatmap')}</CardTitle>
            </CardHeader>
            <CardContent className="h-[300px]">
              {analyticsLoading ? (
                <Skeleton className="w-full h-full" />
              ) : heatmapData.length > 0 ? (
                <HeatmapChart data={heatmapData} />
              ) : (
                <NoData />
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === 'links' && (
        <div className="space-y-4 animate-in slide-in-from-bottom-4 duration-300">
          <div className="flex justify-between items-center bg-card p-4 rounded-xl border border-border">
            <div>
              <h3 className="text-lg font-bold">
                {t('campaigns.details.manage_links_title')}
              </h3>
              <p className="text-sm text-muted-foreground">
                {t('campaigns.details.manage_links_desc')}
              </p>
            </div>
            <CreateLinkFeature
              selectedCampaignId={id}
              onSuccess={handleCreateSuccess}
            />
          </div>

          {linksLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-24 bg-muted/50 rounded-xl animate-pulse"
                />
              ))}
            </div>
          ) : links.length === 0 ? (
            <div className="text-center py-20 bg-muted/20 rounded-xl border border-dashed border-border">
              <p className="text-muted-foreground mb-4">
                {t('campaigns.details.empty_links')}
              </p>
            </div>
          ) : (
            <>
              <div className="grid gap-3">
                {links.map((link) => (
                  <LinkCard key={link.id} link={link} onDelete={setDeleteId} />
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

      <ConfirmationModal
        isOpen={!!deleteId}
        onClose={() => setDeleteId(null)}
        onConfirm={handleDeleteConfirm}
        title="Delete Link?"
        description="Are you sure? This short link will stop working immediately, and users will see a 404 error. All analytics data for this link will be permanently removed."
        confirmLabel="Yes, Delete Link"
        isLoading={isDeleting}
      />
    </div>
  );
};

interface KpiCardProps {
  title: string;
  value: string | number;
  icon: React.ReactElement;
  isText?: boolean;
  loading?: boolean;
}

const KpiCard = ({
  title,
  value,
  icon,
  isText = false,
  loading = false,
}: KpiCardProps) => (
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
        >
          {isText
            ? value
            : new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                value as number
              )}
        </div>
      )}
    </CardContent>
  </Card>
);

const NoData = () => (
  <div className="h-full w-full flex items-center justify-center text-muted-foreground text-sm italic">
    No data available
  </div>
);
