import { Pie } from '@visx/shape';
import { Group } from '@visx/group';
import { withParentSize } from '@visx/responsive';
import { scaleOrdinal } from '@visx/scale';

interface DonutProps {
  parentWidth?: number;
  parentHeight?: number;
  data: Record<string, number>;
}

const DonutChartBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: DonutProps) => {
  const width = parentWidth;
  const height = parentHeight;

  if (width < 10) return null;

  const radius = Math.min(width, height) / 2;
  const innerRadius = radius / 1.6;

  const entries = Object.entries(data).map(([key, value]) => ({
    label: key,
    value,
  }));

  const colorScale = scaleOrdinal({
    domain: entries.map((e) => e.label),
    range: ['#4f46e5', '#ec4899', '#06b6d4', '#f59e0b', '#10b981'],
  });

  return (
    <svg width={width} height={height}>
      <Group top={height / 2} left={width / 2}>
        <Pie
          data={entries}
          pieValue={(d) => d.value}
          outerRadius={radius}
          innerRadius={innerRadius}
          cornerRadius={4}
          padAngle={0.02}
        >
          {({ arcs, path }) => (
            <g>
              {arcs.map((arc, i) => (
                <path
                  key={`arc-${i}`}
                  d={path(arc) || ''}
                  fill={colorScale(arc.data.label)}
                />
              ))}
            </g>
          )}
        </Pie>
        <text
          x={0}
          y={0}
          dy={5}
          textAnchor="middle"
          className="text-2xl font-bold fill-slate-900 dark:fill-white"
        >
          {entries.reduce((acc, v) => acc + v.value, 0)}
        </text>
      </Group>
    </svg>
  );
};

export const DonutChart = withParentSize<DonutProps>(DonutChartBase);
