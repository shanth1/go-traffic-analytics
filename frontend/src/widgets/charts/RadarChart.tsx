import { Group } from '@visx/group';
import { scaleLinear } from '@visx/scale';
import { Point } from '@visx/point';
import { ParentSize } from '@visx/responsive';

const data = [
  { label: 'Quality', value: 90 },
  { label: 'Suspicious', value: 20 },
  { label: 'Bots', value: 10 },
  { label: 'Proxies', value: 15 },
  { label: 'Repeated', value: 30 },
];

const Radar = ({ width, height }: { width: number; height: number }) => {
  const radius = Math.min(width, height) / 2 - 30;
  const angleScale = (i: number) => (i * (Math.PI * 2)) / data.length;
  const radiusScale = scaleLinear({ range: [0, radius], domain: [0, 100] });

  const points = data.map((d, i) => {
    const angle = angleScale(i) - Math.PI / 2;
    const r = radiusScale(d.value);
    return new Point({ x: r * Math.cos(angle), y: r * Math.sin(angle) });
  });

  return (
    <svg width={width} height={height} className="overflow-visible">
      <Group top={height / 2} left={width / 2}>
        {/* Сетка */}
        {[25, 50, 75, 100].map((r) => (
          <circle
            key={r}
            r={radiusScale(r)}
            fill="none"
            stroke="#334155"
            strokeDasharray="4"
          />
        ))}
        {/* Оси */}
        {data.map((d, i) => {
          const angle = angleScale(i) - Math.PI / 2;
          const end = new Point({
            x: radius * Math.cos(angle),
            y: radius * Math.sin(angle),
          });
          return (
            <g key={i}>
              <line x1={0} y1={0} x2={end.x} y2={end.y} stroke="#475569" />
              <text
                x={end.x * 1.15}
                y={end.y * 1.15}
                textAnchor="middle"
                fill="#94a3b8"
                fontSize={10}
                dy=".3em"
              >
                {d.label}
              </text>
            </g>
          );
        })}
        {/* Полигон данных */}
        <polygon
          points={points.map((p) => `${p.x},${p.y}`).join(' ')}
          fill="rgba(16, 185, 129, 0.2)"
          stroke="#10b981"
          strokeWidth={2}
        />
        {points.map((p, i) => (
          <circle key={i} cx={p.x} cy={p.y} r={3} fill="#10b981" />
        ))}
      </Group>
    </svg>
  );
};

export const RadarWidget = () => (
  <div className="h-64 w-full">
    <ParentSize>
      {({ width, height }) => <Radar width={width} height={height} />}
    </ParentSize>
  </div>
);
