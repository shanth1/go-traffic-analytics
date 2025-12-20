import { useEffect, useState } from 'react';
import * as topojson from 'topojson-client';
import { Mercator } from '@visx/geo';
import { scaleQuantize } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import type { Feature, Geometry } from 'geojson'; // Типы из geojson (обычно есть в @types/geojson)

interface GeoMapProps {
  parentWidth?: number;
  parentHeight?: number;
  data: { [country: string]: number };
}

// Тип для состояния world data
interface TopologyData {
  features: Feature<Geometry, { name: string }>[];
}

const GeoMapBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: GeoMapProps) => {
  const width = parentWidth;
  const height = parentHeight;

  const [world, setWorld] = useState<TopologyData | null>(null);

  useEffect(() => {
    fetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json')
      .then((res) => res.json())
      .then((topology) => {
        // Явное приведение типов для topojson
        const geojson = topojson.feature(
          topology,
          topology.objects.countries
        ) as unknown as TopologyData;
        setWorld(geojson);
      });
  }, []);

  const maxVal = Math.max(...Object.values(data), 1);
  const colorScale = scaleQuantize({
    domain: [0, maxVal],
    range: [
      '#e0e7ff',
      '#c7d2fe',
      '#a5b4fc',
      '#818cf8',
      '#6366f1',
      '#4f46e5',
      '#4338ca',
    ],
  });

  if (!world || width < 10)
    return <div className="animate-pulse bg-slate-100 w-full h-full rounded" />;

  return (
    <svg width={width} height={height}>
      <rect width={width} height={height} fill="transparent" />
      {/*
         @ts-ignore: visx types mismatch with recent React versions or generic GeoJSON types
      */}
      <Mercator<Feature<Geometry, { name: string }>>
        data={world.features}
        scale={width / 6.5}
        translate={[width / 2, height / 1.5]}
      >
        {(mercator) => (
          <g>
            {mercator.features.map((feature, i) => {
              // Безопасный доступ к свойствам
              const countryName = feature.feature.properties?.name || '';
              const value = data[countryName] || 0;

              return (
                <path
                  key={`map-feature-${i}`}
                  d={mercator.path(feature.feature) || ''}
                  fill={value ? colorScale(value) : '#f1f5f9'}
                  stroke="#cbd5e1"
                  strokeWidth={0.5}
                  className="transition-colors duration-300 hover:fill-indigo-400"
                />
              );
            })}
          </g>
        )}
      </Mercator>
    </svg>
  );
};

export const GeoMap = withParentSize<GeoMapProps>(GeoMapBase);
