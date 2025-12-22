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

        // 1. Собираем все уникальные ключи
        const allKeysSet = new Set<string>();
        rawData.forEach((item) => {
          if (item.values) {
            Object.keys(item.values).forEach((k) => allKeysSet.add(k));
          }
        });
        const collectedKeys = Array.from(allKeysSet);

        // 2. Трансформируем данные
        const processedData: StreamChartData[] = rawData.map((d) => {
          // Инициализируем объект. TypeScript знает, что time - это Date.
          const point: StreamChartData = {
            time: new Date(d.time),
          };

          // Заполняем динамические ключи
          collectedKeys.forEach((key) => {
            // Используем оператор ?? (nullish coalescing), чтобы не потерять 0, если он придет явно
            // d.values?.[key] вернет number или undefined
            const value = d.values?.[key] ?? 0;

            // Присваиваем значение.
            // TS не ругается, так как value (number) входит в тип (number | Date)
            point[key] = value;
          });

          return point;
        });

        // 3. Сортируем
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
      {/* ... Ваш заголовок ... */}

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
          ) : chartData.length > 0 ? (
            <StreamGraph data={chartData} keys={keys} />
          ) : (
            <div className="h-full flex items-center justify-center text-slate-400">
              Нет данных для отображения графика
            </div>
          )}
        </CardContent>
      </Card>

      {/* ... Остальные карточки ... */}
    </div>
  );
};
