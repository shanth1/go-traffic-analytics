import React, { useEffect, useState, useMemo } from 'react';
import {
  ActivityIcon,
  MousePointerClick,
  ShieldCheckIcon,
  LayersIcon,
  SmartphoneIcon,
  GlobeIcon,
  ServerIcon,
} from 'lucide-react';

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { GeoMap } from '@/widgets/charts/GeoMap';
import { BarListChart } from '@/widgets/charts/BarListChart';
import { DonutChart } from '@/widgets/charts/DonutChart';
import { HierarchyTree } from '@/widgets/charts/HierarchyTree';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';

import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { analyticsApi } from '@/entities/analytics/api';
import { toRFC3339 } from '@/shared/lib/date';
import { cn } from '@/shared/lib/utils';

import type {
  HierarchyNode,
  AnalyticsSummary,
  GeoPoint,
  StreamChartData,
  HeatmapPoint,
  TrafficQuality,
  CategoryStat,
} from '@/shared/api/types';

export const DashboardPage = () => {
  const { startDate, endDate } = useAnalyticsFilter();
  const [loading, setLoading] = useState(true);

  // --- Data States ---
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null);
  const [treeData, setTreeData] = useState<HierarchyNode | null>(null);

  // Charts Data
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);
  const [geoData, setGeoData] = useState<GeoPoint[]>([]);
  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [qualityData, setQualityData] = useState<TrafficQuality | null>(null);

  // Stats
  const [statsDevice, setStatsDevice] = useState<CategoryStat[]>([]);
  const [statsOS, setStatsOS] = useState<CategoryStat[]>([]);
  const [statsBrowser, setStatsBrowser] = useState<CategoryStat[]>([]);

  // --- Derived Metrics ---
  const deviceDonutData = useMemo(() => {
    const map: Record<string, number> = {};
    statsDevice.forEach((d) => (map[d.name] = d.value));
    return map;
  }, [statsDevice]);

  const topCountries = useMemo(() => {
    return geoData
      .sort((a, b) => b.value - a.value)
      .slice(0, 6)
      .map((g) => ({ name: g.country, value: g.value, share: 0 }));
  }, [geoData]);

  const inventoryStats = useMemo(() => {
    let campaigns = 0;
    let links = 0;
    if (treeData && treeData.children) {
      campaigns = treeData.children.length;
      treeData.children.forEach((c) => {
        if (c.children) links += c.children.length;
      });
    }
    return { campaigns, links };
  }, [treeData]);

  // --- Fetching ---
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const params = {
          from: toRFC3339(startDate),
          to: toRFC3339(endDate),
        };

        const [sum, tree, stream, geo, heatmap, quality, dev, os, browser] =
          await Promise.all([
            analyticsApi.getSummary(params),
            analyticsApi.getHierarchy(),
            analyticsApi.getStream(undefined, {
              ...params,
              group_by: 'device',
            }),
            analyticsApi.getGeoStats(params),
            analyticsApi.getHeatmap(undefined, params),
            analyticsApi.getQuality(undefined, params),
            analyticsApi.getStats('device', params),
            analyticsApi.getStats('os', params),
            analyticsApi.getStats('browser', params),
          ]);

        setSummary(sum);
        setTreeData(tree);

        // --- Process Stream (Sorting Logic) ---
        const keyTotals: Record<string, number> = {};

        // Calculate totals for each key to determine stack order
        stream.forEach((item) => {
          if (item.values) {
            Object.entries(item.values).forEach(([key, val]) => {
              keyTotals[key] = (keyTotals[key] || 0) + val;
            });
          }
        });

        // Sort descending: largest values first -> rendered at the bottom of the stack
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
        // ---------------------------------------

        setGeoData(geo);
        setHeatmapData(heatmap);
        setQualityData(quality);
        setStatsDevice(dev);
        setStatsOS(os);
        setStatsBrowser(browser);
      } catch (e) {
        console.error('Failed to load dashboard', e);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [startDate, endDate]);

  if (loading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center gap-4">
        <div className="w-10 h-10 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
        <p className="text-slate-500 font-medium">
          Assembling your command center...
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 pb-20">
      {/* --- HEADER --- */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-slate-200 dark:border-slate-800 pb-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
            Overview
          </h1>
          <p className="text-slate-500 mt-1">
            Global metrics across all campaigns and links.
          </p>
        </div>
        <DateRangePicker />
      </div>

      {/* --- 1. KPI ROW --- */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <KpiCard
          title="Total Clicks"
          value={summary?.total_clicks || 0}
          icon={<MousePointerClick className="text-indigo-500" />}
          trend="Volumetric"
        />
        <KpiCard
          title="Human Traffic"
          value={`${qualityData?.human_score || 0}%`}
          icon={<ShieldCheckIcon className="text-emerald-500" />}
          trend="Quality Score"
          isText
        />
        <KpiCard
          title="Active Campaigns"
          value={inventoryStats.campaigns}
          icon={<LayersIcon className="text-amber-500" />}
          trend="Inventory"
        />
        <KpiCard
          title="Active Links"
          value={inventoryStats.links}
          icon={<ActivityIcon className="text-pink-500" />}
          trend="Inventory"
        />
      </div>

      {/* --- 2. MAIN TREND --- */}
      <Card className="overflow-hidden">
        <CardHeader>
          <CardTitle>Traffic Trends</CardTitle>
          <CardDescription>
            Click volume distribution by device type over time
          </CardDescription>
        </CardHeader>
        <CardContent className="h-[350px]">
          {streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <NoData />
          )}
        </CardContent>
      </Card>

      {/* --- 3. GEOGRAPHY --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Global Heatmap</CardTitle>
            <CardDescription>Traffic intensity by region</CardDescription>
          </CardHeader>
          <CardContent className="h-[400px] w-full p-0 overflow-hidden">
            {geoData.length > 0 ? <GeoMap data={geoData} /> : <NoData />}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Top Countries</CardTitle>
          </CardHeader>
          <CardContent>
            <BarListChart data={topCountries} color="bg-indigo-500" />
          </CardContent>
        </Card>
      </div>

      {/* --- 4. TECHNOLOGY STACK --- */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <SmartphoneIcon size={18} className="text-slate-400" /> Device
              Share
            </CardTitle>
          </CardHeader>
          <CardContent className="h-[250px] flex items-center justify-center">
            {Object.keys(deviceDonutData).length > 0 ? (
              <DonutChart data={deviceDonutData} />
            ) : (
              <NoData />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ServerIcon size={18} className="text-slate-400" /> Operating
              Systems
            </CardTitle>
          </CardHeader>
          <CardContent>
            <BarListChart data={statsOS} color="bg-emerald-500" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <GlobeIcon size={18} className="text-slate-400" /> Browsers
            </CardTitle>
          </CardHeader>
          <CardContent>
            <BarListChart data={statsBrowser} color="bg-amber-500" />
          </CardContent>
        </Card>
      </div>

      {/* --- 5. BEHAVIOR & STRUCTURE --- */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Peak Hours (UTC)</CardTitle>
            <CardDescription>Engagement heatmap</CardDescription>
          </CardHeader>
          <CardContent className="h-[300px]">
            {heatmapData.length > 0 ? (
              <HeatmapChart data={heatmapData} />
            ) : (
              <NoData />
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Traffic Quality Radar</CardTitle>
            <CardDescription>Bot vs Human Analysis</CardDescription>
          </CardHeader>
          <CardContent className="h-[300px]">
            {qualityData ? <QualityRadar data={qualityData} /> : <NoData />}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Account Hierarchy</CardTitle>
            <CardDescription>Campaigns & Links</CardDescription>
          </CardHeader>
          <CardContent className="h-[300px]">
            {treeData ? <HierarchyTree data={treeData} /> : <NoData />}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

// --- Sub-components ---

interface KpiCardProps {
  title: string;
  value: number | string;
  icon: React.ReactNode;
  trend: string;
  isText?: boolean;
}

const KpiCard = ({
  title,
  value,
  icon,
  trend,
  isText = false,
}: KpiCardProps) => (
  <Card className="hover:border-indigo-200 dark:hover:border-indigo-900 transition-colors">
    <CardContent className="p-6">
      <div className="flex items-center justify-between mb-4">
        <div className="p-2 bg-slate-100 dark:bg-slate-900 rounded-lg">
          {icon}
        </div>
        <span className="text-xs font-medium text-slate-400 bg-slate-50 dark:bg-slate-900 px-2 py-1 rounded">
          {trend}
        </span>
      </div>
      <div>
        <p className="text-sm font-medium text-slate-500">{title}</p>
        <h3
          className={cn(
            'font-bold text-slate-900 dark:text-slate-50',
            isText ? 'text-xl' : 'text-3xl'
          )}
        >
          {typeof value === 'number'
            ? new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
                value
              )
            : value}
        </h3>
      </div>
    </CardContent>
  </Card>
);

const NoData = () => (
  <div className="h-full w-full flex items-center justify-center text-slate-400 text-sm italic bg-slate-50/50 dark:bg-slate-900/50 rounded-lg border border-dashed border-slate-200 dark:border-slate-800">
    No data available
  </div>
);
