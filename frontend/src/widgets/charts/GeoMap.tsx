import { useEffect, useState, useMemo } from 'react';
import * as topojson from 'topojson-client';
import { Mercator } from '@visx/geo';
import { scaleQuantize } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import type { Feature, Geometry } from 'geojson';
import countries from 'i18n-iso-countries';
import enLocale from 'i18n-iso-countries/langs/en.json';
import type { GeoPoint } from '@/shared/api/types';

countries.registerLocale(enLocale);

interface GeoMapProps {
  parentWidth?: number;
  parentHeight?: number;
  data: GeoPoint[];
}

interface TopologyData {
  features: Feature<Geometry, { name: string }>[];
}

const GeoMapBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data = [],
}: GeoMapProps) => {
  const width = Math.max(parentWidth, 100);
  const height = Math.max(parentHeight, 100);

  const [world, setWorld] = useState<TopologyData | null>(null);

  useEffect(() => {
    fetch('https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json')
      .then((res) => res.json())
      .then((topology) => {
        const geojson = topojson.feature(
          topology,
          topology.objects.countries
        ) as unknown as TopologyData;
        setWorld(geojson);
      })
      .catch((err) => console.error('Failed to load map topology:', err));
  }, []);

  // ISO-2 ("RU") -> ISO Numeric ("643")
  const dataMap = useMemo(() => {
    const map = new Map<string, number>();

    data.forEach((item) => {
      const numericCode = countries.alpha2ToNumeric(item.country);

      if (numericCode) {
        const normalizedId = String(Number(numericCode));
        map.set(normalizedId, item.value);
      }
    });
    return map;
  }, [data]);

  const colorScale = useMemo(() => {
    const maxVal = data.length > 0 ? Math.max(...data.map((d) => d.value)) : 1;

    return scaleQuantize({
      domain: [0, maxVal],
      range: [
        '#f3f4f6',
        '#c7d2fe',
        '#a5b4fc',
        '#818cf8',
        '#6366f1',
        '#4f46e5',
        '#3730a3',
      ],
    });
  }, [data]);

  if (!world) {
    return (
      <div className="flex items-center justify-center w-full h-full text-slate-400 text-sm animate-pulse">
        Loading Map...
      </div>
    );
  }

  return (
    <svg width={width} height={height} className="touch-none select-none">
      <rect width={width} height={height} fill="transparent" />

      <Mercator<Feature<Geometry, { name: string }>>
        data={world.features}
        scale={width / 6.3}
        translate={[width / 2, height / 1.5]}
      >
        {(mercator) => (
          <g>
            {mercator.features.map((feature, i) => {
              const featureId = String(Number(feature.feature.id));

              const value = dataMap.get(featureId);
              const hasData = value !== undefined && value > 0;
              const countryName = feature.feature.properties?.name || 'Unknown';

              return (
                <path
                  key={`map-feature-${i}`}
                  d={mercator.path(feature.feature) || ''}
                  fill={hasData ? colorScale(value) : '#e2e8f0'}
                  stroke="#ffffff"
                  strokeWidth={0.5}
                  className="transition-all duration-200 outline-none hover:opacity-80"
                  style={{
                    cursor: hasData ? 'pointer' : 'default',
                    fill: hasData ? colorScale(value) : '#f1f5f9',
                  }}
                >
                  <title>{`${countryName}: ${value || 0}`}</title>
                </path>
              );
            })}
          </g>
        )}
      </Mercator>
    </svg>
  );
};

export const GeoMap = withParentSize<GeoMapProps>(GeoMapBase);
