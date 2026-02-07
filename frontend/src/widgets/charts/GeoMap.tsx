import { useEffect, useState, useMemo } from 'react';
import * as topojson from 'topojson-client';
import { Mercator } from '@visx/geo';
import { scaleLinear } from '@visx/scale';
import { withParentSize } from '@visx/responsive';
import type { Feature, Geometry } from 'geojson';
import countries from 'i18n-iso-countries';
import enLocale from 'i18n-iso-countries/langs/en.json';
import type { GeoPoint } from '@/shared/api/types';
import { THEME_COLORS } from '@/shared/config/theme';

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

  // FIX: Используем scaleLinear для прозрачности, чтобы не хардкодить цвета
  const opacityScale = useMemo(() => {
    const maxVal = data.length > 0 ? Math.max(...data.map((d) => d.value)) : 1;
    return scaleLinear<number>({
      domain: [0, maxVal],
      range: [0.2, 1], // Страны с данными будут иметь opacity от 0.2 до 1
      clamp: true,
    });
  }, [data]);

  if (!world) {
    return (
      <div className="flex items-center justify-center w-full h-full text-muted-foreground text-sm animate-pulse">
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
                  // FIX: Основной цвет = Primary, фон = Muted
                  fill={hasData ? THEME_COLORS.primary : THEME_COLORS.muted}
                  // FIX: Прозрачность зависит от данных
                  fillOpacity={hasData ? opacityScale(value) : 1}
                  stroke={THEME_COLORS.background}
                  strokeWidth={0.5}
                  className="transition-all duration-200 outline-none hover:opacity-80"
                  style={{
                    cursor: hasData ? 'pointer' : 'default',
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
