import { useMemo } from 'react';
import { AreaStack } from '@visx/shape';
import { scaleTime, scaleLinear, scaleOrdinal } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import { curveMonotoneX } from '@visx/curve';
import type { StreamChartData } from '@/shared/api/types';

interface StreamGraphProps {
  parentWidth?: number;
  parentHeight?: number;
  data: StreamChartData[];
  keys: string[];
}

const getDate = (d: StreamChartData) => d.time;

const colorScale = scaleOrdinal({
  range: ['#4f46e5', '#ec4899', '#06b6d4', '#fbbf24', '#34d399'],
});

const StreamGraphBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
  keys,
}: StreamGraphProps) => {
  const width = parentWidth;
  const height = parentHeight;

  const yMaxVal = useMemo(() => {
    if (data.length === 0) return 0;
    const maxVal = Math.max(
      ...data.map((d) => {
        return keys.reduce((acc, key) => {
          const val = d[key];
          return acc + (typeof val === 'number' ? val : 0);
        }, 0);
      })
    );
    return maxVal * 1.2;
  }, [data, keys]);

  const xScale = useMemo(
    () =>
      scaleTime({
        range: [0, width],
        domain: [
          Math.min(...data.map((d) => getDate(d).getTime())),
          Math.max(...data.map((d) => getDate(d).getTime())),
        ],
      }),
    [width, data]
  );

  const yScale = useMemo(
    () =>
      scaleLinear({
        range: [height, 0],
        domain: [0, yMaxVal],
      }),
    [height, yMaxVal]
  );

  if (width < 10 || data.length === 0) return null;

  return (
    <svg width={width} height={height}>
      <AreaStack
        data={data}
        keys={keys}
        x={(d) => xScale(getDate(d.data)) ?? 0}
        y0={(d) => yScale(d[0])}
        y1={(d) => yScale(d[1])}
        value={(d, key) => (d[key] as number) || 0}
        curve={curveMonotoneX}
      >
        {({ stacks, path }) =>
          stacks.map((stack) => (
            <path
              key={`stack-${stack.key}`}
              d={path(stack) || ''}
              fill={colorScale(stack.key)}
              stroke="white"
              strokeWidth={1}
              opacity={0.8}
            />
          ))
        }
      </AreaStack>
    </svg>
  );
};

export const StreamGraph = withParentSize<StreamGraphProps>(StreamGraphBase);
