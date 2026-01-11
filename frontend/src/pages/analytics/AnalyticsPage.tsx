import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import { HeatmapChart } from '@/widgets/charts/HeatmapChart';
import { QualityRadar } from '@/widgets/charts/QualityRadar';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { toRFC3339 } from '@/shared/lib/date';
import { analyticsApi } from '@/entities/analytics/api';
import type {
  StreamChartData,
  HeatmapPoint,
  TrafficQuality,
} from '@/shared/api/types';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const { startDate, endDate } = useAnalyticsFilter();

  // Local State for Charts
  const [streamData, setStreamData] = useState<StreamChartData[]>([]);
  const [streamKeys, setStreamKeys] = useState<string[]>([]);

  const [heatmapData, setHeatmapData] = useState<HeatmapPoint[]>([]);
  const [qualityData, setQualityData] = useState<TrafficQuality | null>(null);

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

        // Parallel fetching
        const [streamRes, heatmapRes, qualityRes] = await Promise.all([
          analyticsApi.getStream(id, params),
          analyticsApi.getHeatmap(id, params),
          analyticsApi.getQuality(id, params),
        ]);

        // Process Stream Data
        const allKeysSet = new Set<string>();
        streamRes.forEach((item) => {
          if (item.values) {
            Object.keys(item.values).forEach((k) => allKeysSet.add(k));
          }
        });
        const collectedKeys = Array.from(allKeysSet);

        const processedStream: StreamChartData[] = streamRes.map((d) => {
          const point: StreamChartData = {
            time: new Date(d.time),
          };
          collectedKeys.forEach((key) => {
            point[key] = d.values?.[key] ?? 0;
          });
          return point;
        });
        processedStream.sort((a, b) => a.time.getTime() - b.time.getTime());

        setStreamData(processedStream);
        setStreamKeys(collectedKeys);

        // Process Heatmap & Quality
        setHeatmapData(heatmapRes);
        setQualityData(qualityRes);
      } catch (e) {
        console.error('Failed to load analytics data', e);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [id, startDate, endDate]);

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Header with Picker */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-slate-50">
            Link Analytics
          </h1>
          <p className="text-slate-500 text-sm">
            Deep dive into traffic sources and quality.
          </p>
        </div>
        <DateRangePicker />
      </div>

      {/* Main Stream Graph */}
      <Card>
        <CardHeader>
          <CardTitle>Traffic Dynamics</CardTitle>
          <p className="text-sm text-slate-400">Sessions over time</p>
        </CardHeader>
        <CardContent className="h-[350px]">
          {loading ? (
            <div className="h-full flex items-center justify-center">
              Loading...
            </div>
          ) : streamData.length > 0 ? (
            <StreamGraph data={streamData} keys={streamKeys} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-400">
              No traffic data for this period
            </div>
          )}
        </CardContent>
      </Card>

      {/* Secondary Metrics Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Heatmap */}
        <Card>
          <CardHeader>
            <CardTitle>Activity Heatmap</CardTitle>
            <p className="text-sm text-slate-400">
              Best time for engagement (UTC)
            </p>
          </CardHeader>
          <CardContent className="h-[300px]">
            {loading ? (
              <div className="h-full flex items-center justify-center">
                Loading...
              </div>
            ) : heatmapData.length > 0 ? (
              <HeatmapChart data={heatmapData} />
            ) : (
              <div className="h-full flex items-center justify-center text-slate-400">
                No activity data
              </div>
            )}
          </CardContent>
        </Card>

        {/* Quality Radar */}
        <Card>
          <CardHeader>
            <CardTitle>Traffic Quality Score</CardTitle>
            <p className="text-sm text-slate-400">
              AI-based fraud detection metrics
            </p>
          </CardHeader>
          <CardContent className="h-[300px]">
            {loading ? (
              <div className="h-full flex items-center justify-center">
                Loading...
              </div>
            ) : qualityData ? (
              <QualityRadar data={qualityData} />
            ) : (
              <div className="h-full flex items-center justify-center text-slate-400">
                No quality data
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
