import React, { useMemo, useCallback } from 'react';
import { AreaStack, Line, Bar } from '@visx/shape';
import { scaleTime, scaleLinear, scaleOrdinal } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import { curveMonotoneX } from '@visx/curve';
import { GridRows, GridColumns } from '@visx/grid';
import { AxisBottom, AxisRight } from '@visx/axis';
import { useTooltip, useTooltipInPortal, defaultStyles } from '@visx/tooltip';
import { localPoint } from '@visx/event';
import { LegendOrdinal } from '@visx/legend';
import { bisector } from 'd3-array';
import { timeFormat } from 'd3-time-format';

import type { StreamChartData } from '@/shared/api/types';
import { CHART_COLORS, THEME_COLORS } from '@/shared/config/theme';

const margin = { top: 20, right: 30, bottom: 50, left: 0 };
const formatDate = timeFormat('%d %b');

interface StreamGraphProps {
  parentWidth?: number;
  parentHeight?: number;
  data: StreamChartData[];
  keys: string[];
}

const getDate = (d: StreamChartData) => d.time;
const bisectDate = bisector<StreamChartData, Date>((d) => d.time).left;

const StreamGraphBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
  keys,
}: StreamGraphProps) => {
  const {
    tooltipOpen,
    tooltipLeft,
    tooltipTop,
    tooltipData,
    hideTooltip,
    showTooltip,
  } = useTooltip<StreamChartData>();

  const { containerRef, TooltipInPortal } = useTooltipInPortal({
    scroll: true,
  });

  const width = parentWidth;
  const height = parentHeight;

  const xMax = width - margin.left - margin.right;
  const yMax = height - margin.top - margin.bottom;

  const xScale = useMemo(
    () =>
      scaleTime({
        range: [0, xMax],
        domain: [
          Math.min(...data.map(getDate).map((d) => d.getTime())),
          Math.max(...data.map(getDate).map((d) => d.getTime())),
        ],
      }),
    [xMax, data]
  );

  const yMaxVal = useMemo(() => {
    if (data.length === 0) return 0;
    return Math.max(
      ...data.map((d) =>
        keys.reduce((acc, k) => {
          const val = d[k];
          return acc + (typeof val === 'number' ? val : 0);
        }, 0)
      )
    );
  }, [data, keys]);

  const yScale = useMemo(
    () =>
      scaleLinear({
        range: [yMax, 0],
        domain: [0, yMaxVal * 1.1],
        nice: true,
      }),
    [yMax, yMaxVal]
  );

  const colorScale = useMemo(
    () => scaleOrdinal({ domain: keys, range: CHART_COLORS }),
    [keys]
  );

  const handleTooltip = useCallback(
    (
      event: React.TouchEvent<SVGRectElement> | React.MouseEvent<SVGRectElement>
    ) => {
      const { x } = localPoint(event) || { x: 0 };
      const x0 = xScale.invert(x - margin.left);

      const index = bisectDate(data, x0, 1);
      const d0 = data[index - 1];
      const d1 = data[index];
      let d = d0;
      if (d1 && d0) {
        d =
          x0.valueOf() - d0.time.valueOf() > d1.time.valueOf() - x0.valueOf()
            ? d1
            : d0;
      }

      if (d) {
        showTooltip({
          tooltipData: d,
          tooltipLeft: xScale(d.time) + margin.left,
          tooltipTop: yMax + margin.top,
        });
      }
    },
    [showTooltip, xScale, data, yMax]
  );

  if (width < 10) return null;

  return (
    <div className="relative">
      <svg ref={containerRef} width={width} height={height}>
        <rect x={0} y={0} width={width} height={height} fill="transparent" />

        <g transform={`translate(${margin.left},${margin.top})`}>
          <GridRows
            scale={yScale}
            width={xMax}
            strokeDasharray="3,3"
            stroke={THEME_COLORS.border}
          />
          <GridColumns
            scale={xScale}
            height={yMax}
            strokeDasharray="3,3"
            stroke={THEME_COLORS.border}
          />

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
                  stroke={THEME_COLORS.background}
                  strokeWidth={0.5}
                  fillOpacity={0.85}
                />
              ))
            }
          </AreaStack>

          <AxisBottom
            top={yMax}
            scale={xScale}
            tickFormat={(d) => formatDate(d as Date)}
            stroke={THEME_COLORS.border}
            tickStroke={THEME_COLORS.border}
            tickLabelProps={() => ({
              fill: THEME_COLORS.foreground,
              fontSize: 11,
              textAnchor: 'middle',
              opacity: 0.6,
            })}
          />

          <AxisRight
            scale={yScale}
            left={xMax}
            numTicks={5}
            stroke="transparent"
            tickStroke="transparent"
            tickLabelProps={() => ({
              fill: THEME_COLORS.foreground,
              fontSize: 10,
              textAnchor: 'start',
              dx: 4,
              opacity: 0.5,
            })}
          />

          {tooltipOpen && tooltipLeft !== undefined && (
            <Line
              from={{ x: tooltipLeft - margin.left, y: 0 }}
              to={{ x: tooltipLeft - margin.left, y: yMax }}
              stroke={THEME_COLORS.foreground}
              strokeWidth={1}
              pointerEvents="none"
              strokeDasharray="5,2"
              opacity={0.5}
            />
          )}

          <Bar
            x={0}
            y={0}
            width={xMax}
            height={yMax}
            fill="transparent"
            rx={14}
            onTouchStart={handleTooltip}
            onTouchMove={handleTooltip}
            onMouseMove={handleTooltip}
            onMouseLeave={() => hideTooltip()}
          />
        </g>
      </svg>

      {tooltipOpen && tooltipData && (
        <TooltipInPortal
          key={tooltipData.time.toString()}
          top={tooltipTop}
          left={tooltipLeft}
          style={{
            ...defaultStyles,
            backgroundColor: THEME_COLORS.background,
            color: THEME_COLORS.foreground,
            border: `1px solid ${THEME_COLORS.border}`,
            minWidth: 120,
            zIndex: 100,
          }}
        >
          <div className="text-xs font-bold mb-1">
            {formatDate(tooltipData.time)}
          </div>
          <div className="space-y-1">
            {keys.map((key) => {
              const val = tooltipData[key];
              if (typeof val !== 'number') return null;
              if (val === 0) return null;

              return (
                <div
                  key={key}
                  className="flex justify-between items-center text-xs"
                >
                  <span className="flex items-center gap-1">
                    <span
                      className="w-2 h-2 rounded-full"
                      style={{ background: colorScale(key) }}
                    ></span>
                    {key}
                  </span>
                  <span className="font-mono">{val}</span>
                </div>
              );
            })}
          </div>
        </TooltipInPortal>
      )}

      <div className="flex flex-wrap justify-center gap-4 mt-2">
        <LegendOrdinal
          scale={colorScale}
          direction="row"
          labelMargin="0 15px 0 0"
        >
          {(labels) => (
            <div className="flex flex-wrap gap-4">
              {labels.map((label, i) => (
                <div
                  key={`legend-${i}`}
                  className="flex items-center cursor-pointer hover:opacity-70 transition-opacity"
                >
                  <div
                    className="w-3 h-3 rounded-full mr-2"
                    style={{ backgroundColor: label.value }}
                  />
                  <span className="text-xs text-muted-foreground font-medium">
                    {label.text}
                  </span>
                </div>
              ))}
            </div>
          )}
        </LegendOrdinal>
      </div>
    </div>
  );
};

export const StreamGraph = withParentSize<StreamGraphProps>(StreamGraphBase);
