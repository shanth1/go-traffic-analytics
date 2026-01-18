import { useMemo } from 'react';
import { Group } from '@visx/group';
import { scaleLinear } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import { AxisBottom, AxisLeft } from '@visx/axis';
import type { HeatmapPoint } from '@/shared/api/types';
import { THEME_COLORS } from '@/shared/config/theme';

interface HeatmapProps {
  parentWidth?: number;
  parentHeight?: number;
  data: HeatmapPoint[];
}

const DAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const HOURS = Array.from({ length: 24 }, (_, i) => i);

const HeatmapChartBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: HeatmapProps) => {
  const width = parentWidth;
  const height = parentHeight;
  const margin = { top: 10, left: 40, right: 10, bottom: 30 };
  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;

  const binData = useMemo(() => {
    const map = new Map<string, number>();
    data.forEach((d) => {
      map.set(`${d.day}-${d.hour}`, d.count);
    });
    return map;
  }, [data]);

  const maxCount = useMemo(() => {
    return data.length > 0 ? Math.max(...data.map((d) => d.count)) : 0;
  }, [data]);

  // FIX: Используем scaleLinear для Opacity (числа), а не для Цветов (строк)
  const opacityScale = useMemo(() => scaleLinear<number>({
    domain: [0, maxCount],
    range: [0.05, 1], // От почти прозрачного до полного цвета
    clamp: true,
  }), [maxCount]);

  const xScale = scaleLinear({
    domain: [0, 24],
    range: [0, xMax],
  });

  const yScale = scaleLinear({
    domain: [0, 7],
    range: [0, yMax],
  });

  if (width < 10) return null;

  const binWidth = xMax / 24;
  const binHeight = yMax / 7;

  return (
    <svg width={width} height={height}>
      <Group top={margin.top} left={margin.left}>
        {DAYS.map((day, dIndex) =>
          HOURS.map((hour, hIndex) => {
            const count = binData.get(`${dIndex}-${hIndex}`) || 0;
            return (
              <rect
                key={`cell-${dIndex}-${hIndex}`}
                x={xScale(hIndex)}
                y={yScale(dIndex)}
                width={binWidth - 2}
                height={binHeight - 2}
                // FIX: Цвет берется из CSS переменной
                fill={THEME_COLORS.primary}
                // FIX: Насыщенность зависит от данных
                fillOpacity={count > 0 ? opacityScale(count) : 0.05}
                rx={2}
                style={{ transition: 'fill-opacity 0.3s ease' }}
              >
                <title>{`${day} ${hour}:00 - ${count} clicks`}</title>
              </rect>
            );
          })
        )}
        <AxisLeft
          scale={yScale}
          top={binHeight / 2}
          tickValues={[0, 1, 2, 3, 4, 5, 6]}
          tickFormat={(v) => DAYS[Number(v)]}
          stroke="transparent"
          tickStroke="transparent"
          tickLabelProps={() => ({
            fill: THEME_COLORS.foreground,
            fontSize: 11,
            textAnchor: 'end',
            dy: 4,
            dx: -5,
            opacity: 0.6
          })}
        />

        <AxisBottom
          scale={xScale}
          top={yMax}
          left={binWidth / 2}
          tickValues={[0, 6, 12, 18, 23]}
          tickFormat={(v) => `${v}:00`}
          stroke="transparent"
          tickStroke="transparent"
          tickLabelProps={() => ({
            fill: THEME_COLORS.foreground,
            fontSize: 10,
            textAnchor: 'middle',
            opacity: 0.6
          })}
        />
      </Group>
    </svg>
  );
};

export const HeatmapChart = withParentSize<HeatmapProps>(HeatmapChartBase);
