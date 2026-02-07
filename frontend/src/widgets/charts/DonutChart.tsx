import { Pie } from '@visx/shape';
import { Group } from '@visx/group';
import { withParentSize } from '@visx/responsive';
import { scaleOrdinal } from '@visx/scale';
import { CHART_COLORS, THEME_COLORS } from '@/shared/config/theme';

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

  // Using theme constants for consistent chart coloring
  const colorScale = scaleOrdinal({
    domain: entries.map((e) => e.label),
    range: CHART_COLORS,
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
          fontSize={24}
          fontWeight="bold"
          fill={THEME_COLORS.foreground}
        >
          {entries.reduce((acc, v) => acc + v.value, 0)}
        </text>
      </Group>
    </svg>
  );
};

export const DonutChart = withParentSize<DonutProps>(DonutChartBase);
