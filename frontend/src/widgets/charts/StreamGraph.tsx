import { useMemo } from 'react';
import { AreaStack } from '@visx/shape';
import { scaleTime, scaleLinear } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import { curveMonotoneX } from '@visx/curve';
import { LinearGradient } from '@visx/gradient';
import type { StackedPoint } from '@/shared/api/types';

interface StreamGraphProps {
  parentWidth?: number;
  parentHeight?: number;
  data: StackedPoint[];
}

const getDate = (d: StackedPoint) => new Date(d.time);
const keys = ['desktop', 'mobile', 'tablet'];

const StreamGraphBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: StreamGraphProps) => {
  const width = parentWidth;
  const height = parentHeight;

  // 1. Сначала вычисляем все данные (Hooks всегда должны вызываться)
  const yMaxVal = useMemo(() => {
    if (data.length === 0) return 0;
    const maxVal = Math.max(
      ...data.map((d) => Object.values(d.values).reduce((a, b) => a + b, 0))
    );
    return maxVal * 1.2;
  }, [data]);

  const xScale = useMemo(
    () =>
      scaleTime({
        range: [0, width], // width может быть 0, это не страшно для создания scale
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

  // 2. Только теперь делаем условный рендеринг
  if (width < 10 || data.length === 0) return null;

  return (
    <svg width={width} height={height}>
      <LinearGradient id="gradient-desktop" from="#4f46e5" to="#818cf8" />
      <LinearGradient id="gradient-mobile" from="#ec4899" to="#f472b6" />
      <LinearGradient id="gradient-tablet" from="#06b6d4" to="#22d3ee" />

      <AreaStack
        data={data}
        keys={keys}
        x={(d) => xScale(getDate(d.data)) ?? 0}
        y0={(d) => yScale(d[0])}
        y1={(d) => yScale(d[1])}
        value={(d, key) => d.values[key] || 0}
        curve={curveMonotoneX}
      >
        {({ stacks, path }) =>
          stacks.map((stack) => (
            <path
              key={`stack-${stack.key}`}
              d={path(stack) || ''}
              fill={`url(#gradient-${stack.key})`}
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
