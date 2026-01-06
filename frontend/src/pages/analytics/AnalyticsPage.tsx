import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import type { StreamChartData } from '@/shared/api/types';
import { analyticsApi } from '@/entities/analytics/api';

export const AnalyticsPage = () => {
  const { id } = useParams();
  const [chartData, setChartData] = useState<StreamChartData[]>([]);
  const [keys, setKeys] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      setLoading(true);
      try {
        const rawData = await analyticsApi.getStream(id);

        const allKeysSet = new Set<string>();
        rawData.forEach((item) => {
          if (item.values) {
            Object.keys(item.values).forEach((k) => allKeysSet.add(k));
          }
        });
        const collectedKeys = Array.from(allKeysSet);

        const processedData: StreamChartData[] = rawData.map((d) => {
          const point: StreamChartData = {
            time: new Date(d.time),
          };

          collectedKeys.forEach((key) => {
            const value = d.values?.[key] ?? 0;

            point[key] = value;
          });

          return point;
        });

        processedData.sort((a, b) => a.time.getTime() - b.time.getTime());

        setKeys(collectedKeys);
        setChartData(processedData);
      } catch (e) {
        console.error('Failed to load stream data', e);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [id]);

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Traffic Dynamics (Stream)</CardTitle>
          <p className="text-sm text-slate-400">
            Device distribution over time
          </p>
        </CardHeader>
        <CardContent className="h-[400px]">
          {loading ? (
            <div className="h-full flex items-center justify-center">
              Loading...
            </div>
          ) : chartData.length > 0 ? (
            <StreamGraph data={chartData} keys={keys} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-400">
              No data to display the chart
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
