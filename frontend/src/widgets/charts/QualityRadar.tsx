import { useMemo } from 'react';
import { Group } from '@visx/group';
import { scaleLinear } from '@visx/scale';
import { Point } from '@visx/point';
import { withParentSize } from '@visx/responsive';
import type { TrafficQuality } from '@/shared/api/types';

interface QualityRadarProps {
  parentWidth?: number;
  parentHeight?: number;
  data: TrafficQuality;
}

const METRICS = [
  { key: 'human_score', label: 'Human' },
  { key: 'bot_score', label: 'Bot Rep.' }, // Inverse logic often applies, but we visualize raw score
  { key: 'geo_diversity_score', label: 'Geo' },
  { key: 'mobile_friendly_score', label: 'Mobile' },
  { key: 'is_suspicious', label: 'Safety' }, // Boolean needs mapping
] as const;

const QualityRadarBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: QualityRadarProps) => {
  const width = parentWidth;
  const height = parentHeight;
  const margin = { top: 30, left: 30, right: 30, bottom: 30 };

  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;
  const radius = Math.min(xMax, yMax) / 2;

  // Map data to radial points
  const radarData = useMemo(() => {
    return METRICS.map((m) => {
      let val = 0;
      if (m.key === 'is_suspicious') {
        // Invert: if suspicious (true), score is 0 (bad). If false, score 100 (good).
        val = data[m.key] ? 20 : 100;
      } else {
        val = Number(data[m.key as keyof TrafficQuality]) || 0;
      }
      return { label: m.label, value: val };
    });
  }, [data]);

  const yScale = scaleLinear<number>({
    range: [0, radius],
    domain: [0, 100],
  });

  // Calculate polygon points
  const angleStep = (Math.PI * 2) / radarData.length;
  const points = radarData.map((_, i) => {
    const angle = i * angleStep - Math.PI / 2;
    return new Point({
      x: Math.cos(angle) * radius,
      y: Math.sin(angle) * radius,
    });
  });

  const valuePoints = radarData.map((d, i) => {
    const r = yScale(d.value);
    const angle = i * angleStep - Math.PI / 2;
    return new Point({
      x: Math.cos(angle) * r,
      y: Math.sin(angle) * r,
    });
  });

  if (width < 10) return null;

  return (
    <svg width={width} height={height}>
      <Group top={height / 2} left={width / 2}>
        {/* Background Grid Circles */}
        {[20, 40, 60, 80, 100].map((tick) => (
          <circle
            key={`grid-${tick}`}
            r={yScale(tick)}
            fill="none"
            stroke="#e2e8f0"
            strokeWidth={1}
            className="dark:stroke-slate-800"
          />
        ))}

        {/* Axes */}
        {points.map((p, i) => (
          <line
            key={`axis-${i}`}
            x1={0}
            y1={0}
            x2={p.x}
            y2={p.y}
            stroke="#cbd5e1"
            strokeWidth={1}
          />
        ))}

        {/* Labels */}
        {points.map((p, i) => {
          // Push label out slightly
          const labelX = p.x * 1.2;
          const labelY = p.y * 1.2;
          return (
            <text
              key={`label-${i}`}
              x={labelX}
              y={labelY}
              dy={5}
              fontSize={11}
              textAnchor="middle"
              className="fill-slate-500 dark:fill-slate-400 font-medium"
            >
              {radarData[i].label}
            </text>
          );
        })}

        {/* The Data Shape */}
        <polygon
          points={valuePoints.map((p) => `${p.x},${p.y}`).join(' ')}
          fill="rgba(79, 70, 229, 0.2)"
          stroke="#4f46e5"
          strokeWidth={2}
        />

        {/* Data Points */}
        {valuePoints.map((p, i) => (
          <circle key={`point-${i}`} cx={p.x} cy={p.y} r={3} fill="#4f46e5" />
        ))}

        {/* Score in Center */}
        <text
          x={0}
          y={0}
          dy={4}
          textAnchor="middle"
          className="text-xs font-bold fill-indigo-600 dark:fill-indigo-400 opacity-0 hover:opacity-100 transition-opacity"
        >
          {data.human_score}
        </text>
      </Group>
    </svg>
  );
};

export const QualityRadar = withParentSize<QualityRadarProps>(QualityRadarBase);
