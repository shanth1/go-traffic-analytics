import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { StreamGraph } from '@/widgets/charts/StreamGraph';
import type { StackedPoint } from '@/shared/api/types';

// Mock generator
const generateStreamData = (points: number): StackedPoint[] => {
  return Array.from({ length: points }).map((_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - (points - i));
    return {
      time: date.toISOString(),
      values: {
        desktop: Math.floor(Math.random() * 100),
        mobile: Math.floor(Math.random() * 150),
        tablet: Math.floor(Math.random() * 50),
      },
    };
  });
};

export const AnalyticsPage = () => {
  const { id } = useParams();
  const [streamData, setStreamData] = useState<StackedPoint[]>([]);

  useEffect(() => {
    // Mock async request via setTimeout
    const timer = setTimeout(() => {
      setStreamData(generateStreamData(20));
    }, 100);

    return () => clearTimeout(timer);
  }, [id]);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">Аналитика ссылки</h1>
          <p className="text-slate-500">ID: {id}</p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Динамика трафика (Stream)</CardTitle>
          <p className="text-sm text-slate-400">
            Распределение устройств за последние 20 дней
          </p>
        </CardHeader>
        <CardContent className="h-[400px]">
          {/* HOC automatically passes parentWidth/parentHeight, we pass 'data' */}
          {streamData.length > 0 ? (
            <StreamGraph data={streamData} />
          ) : (
            <div>Loading...</div>
          )}
        </CardContent>
      </Card>

      {/* Остальная разметка без изменений */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="h-[300px] flex items-center justify-center bg-slate-50 dark:bg-slate-900 border-dashed">
          <span className="text-slate-400">Heatmap (Placeholder)</span>
        </Card>
        <Card className="h-[300px] flex items-center justify-center bg-slate-50 dark:bg-slate-900 border-dashed">
          <span className="text-slate-400">Radar Quality (Placeholder)</span>
        </Card>
      </div>
    </div>
  );
};
