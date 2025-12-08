import { useMemo } from 'react';
import { Stack } from '@visx/shape';
import { scaleLinear, scaleTime } from '@visx/scale';
import { curveNatural } from '@visx/curve'; // Плавные линии
import { Group } from '@visx/group';
import { ParentSize } from '@visx/responsive';
import { motion } from 'framer-motion';

// Генерация фейковых данных для демо
const generateData = () => {
  const points = [];
  const now = Date.now();
  for (let i = 0; i < 20; i++) {
    points.push({
      time: new Date(now - (20 - i) * 86400000),
      iOS: Math.floor(Math.random() * 50) + 20,
      Android: Math.floor(Math.random() * 60) + 30,
      Desktop: Math.floor(Math.random() * 40) + 10,
    });
  }
  return points;
};

const data = generateData();
const keys = ['iOS', 'Android', 'Desktop'];
const colors: Record<string, string> = {
  iOS: '#3b82f6',
  Android: '#10b981',
  Desktop: '#8b5cf6',
};

const Graph = ({ width, height }: { width: number; height: number }) => {
  // Scales
  const xScale = useMemo(
    () =>
      scaleTime({
        range: [0, width],
        domain: [
          Math.min(...data.map((d) => d.time.getTime())),
          Math.max(...data.map((d) => d.time.getTime())),
        ],
      }),
    [width]
  );

  const yScale = useMemo(
    () =>
      scaleLinear({
        range: [height, 0],
        domain: [0, 200], // Примерный макс
      }),
    [height]
  );

  return (
    <svg width={width} height={height}>
      <Group>
        <Stack
          data={data}
          keys={keys}
          x={(d) => xScale(d.data.time) ?? 0}
          y0={(d) => yScale(d[0]) ?? 0}
          y1={(d) => yScale(d[1]) ?? 0}
          curve={curveNatural}
          offset="wiggle" // Эффект "потока"
        >
          {({ stacks, path }) =>
            stacks.map((stack, i) => (
              <motion.path
                key={`stack-${stack.key}`}
                initial={{ opacity: 0, pathLength: 0 }}
                animate={{ opacity: 0.8, pathLength: 1 }}
                transition={{ duration: 1.5, delay: i * 0.2 }}
                d={path(stack) || ''}
                fill={colors[stack.key]}
                stroke="rgba(255,255,255,0.1)"
                whileHover={{ opacity: 1, scale: 1.01 }}
              />
            ))
          }
        </Stack>
      </Group>
    </svg>
  );
};

export const StreamGraphWidget = () => (
  <div className="h-64 w-full">
    <ParentSize>
      {({ width, height }) => <Graph width={width} height={height} />}
    </ParentSize>
  </div>
);
