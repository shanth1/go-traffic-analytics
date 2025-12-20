import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import type { StackedPoint } from '@/shared/api/types';
import { analyticsApi } from '@/entities/analytics/api';

export const AnalyticsPage = () => {
  const { id } = useParams(); // link ID
  const [streamData, setStreamData] = useState<StackedPoint[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    const loadData = async () => {
      setLoading(true);
      try {
        const data = await analyticsApi.getStream(id);
        setStreamData(data);
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
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">Аналитика ссылки</h1>
          <p className="text-slate-500">Link ID: {id}</p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Динамика трафика (Stream)</CardTitle>
          <p className="text-sm text-slate-400">
            Распределение устройств во времени
          </p>
        </CardHeader>
        <CardContent className="h-[400px]">
          {loading ? (
            <div className="h-full flex items-center justify-center">
              Loading...
            </div>
          ) : streamData.length > 0 ? (
            <StreamGraph data={streamData} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-400">
              Нет данных для отображения графика
            </div>
          )}
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="h-[300px] flex items-center justify-center bg-slate-50 dark:bg-slate-900 border-dashed">
          <span className="text-slate-400">Heatmap (In Development)</span>
        </Card>
        <Card className="h-[300px] flex items-center justify-center bg-slate-50 dark:bg-slate-900 border-dashed">
          <span className="text-slate-400">Quality Radar (In Development)</span>
        </Card>
      </div>
    </div>
  );
};
