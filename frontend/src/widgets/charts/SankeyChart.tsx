import { useMemo } from 'react';
import { sankey, sankeyLinkHorizontal, sankeyLeft } from 'd3-sankey';
import { Group } from '@visx/group';
import { withParentSize } from '@visx/responsive';
import type { SankeyData } from '@/shared/api/types';
import { CHART_COLORS, THEME_COLORS } from '@/shared/config/theme';

type NodeDatum = {
  id: string;
  layer: number;
};

type LinkDatum = {
  source: number;
  target: number;
  value: number;
};

type ExtendedNode = NodeDatum & {
  x0: number;
  y0: number;
  x1: number;
  y1: number;
  value: number;
};

type ExtendedLink = {
  source: ExtendedNode;
  target: ExtendedNode;
  value: number;
  width: number;
};

interface SankeyChartProps {
  parentWidth?: number;
  parentHeight?: number;
  data: SankeyData;
}

const SankeyChartBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: SankeyChartProps) => {
  const width = Math.max(parentWidth, 1);
  const height = Math.max(parentHeight, 1);

  const { nodes, links } = useMemo(() => {
    const margin = { top: 20, left: 20, right: 20, bottom: 20 };
    if (!data.nodes.length || !data.links.length) {
      return { nodes: [], links: [] };
    }
    const nodes = data.nodes.map((n) => ({ ...n }));
    const links = data.links.map((l) => ({ ...l }));
    const idToIndex = new Map(nodes.map((n, i) => [n.id, i]));
    const indexedLinks = links
      .map((link) => ({
        source: idToIndex.get(link.source) ?? 0,
        target: idToIndex.get(link.target) ?? 0,
        value: link.value,
      }))
      .filter((l) => l.value > 0);

    const sankeyGenerator = sankey<NodeDatum, LinkDatum>()
      .nodeWidth(15)
      .nodePadding(20)
      .extent([
        [margin.left, margin.top],
        [width - margin.right, height - margin.bottom],
      ])
      .nodeAlign(sankeyLeft);

    try {
      return sankeyGenerator({
        nodes,
        links: indexedLinks,
      }) as { nodes: ExtendedNode[]; links: ExtendedLink[] };
    } catch (e) {
      console.error('Sankey Layout Error', e);
      return { nodes: [], links: [] };
    }
  }, [data, width, height]);

  if (width < 50 || nodes.length === 0) return null;

  return (
    <svg width={width} height={height}>
      <Group>
        {links.map((link, i) => (
          <path
            key={`link-${i}`}
            d={sankeyLinkHorizontal()(link) || undefined}
            stroke={THEME_COLORS.border}
            strokeWidth={Math.max(1, link.width || 0)}
            fill="none"
            strokeOpacity={0.5}
            className="hover:stroke-primary transition-colors"
          >
            <title>{`${link.source.id} → ${link.target.id}: ${link.value}`}</title>
          </path>
        ))}
      </Group>

      <Group>
        {nodes.map((node, i) => (
          <g key={`node-${i}`}>
            <rect
              x={node.x0}
              y={node.y0}
              width={(node.x1 || 0) - (node.x0 || 0)}
              height={(node.y1 || 0) - (node.y0 || 0)}
              fill={CHART_COLORS[node.layer % CHART_COLORS.length]}
              rx={2}
              opacity={0.9}
            >
              <title>{`${node.id}: ${node.value}`}</title>
            </rect>
            <text
              x={
                node.x0 && node.x0 < width / 2
                  ? (node.x1 || 0) + 6
                  : (node.x0 || 0) - 6
              }
              y={(node.y1! + node.y0!) / 2}
              dy="0.35em"
              textAnchor={node.x0 && node.x0 < width / 2 ? 'start' : 'end'}
              fontSize={10}
              fill={THEME_COLORS.foreground}
              className="font-medium pointer-events-none truncate opacity-80"
            >
              {node.id}
            </text>
          </g>
        ))}
      </Group>
    </svg>
  );
};

export const SankeyChart = withParentSize<SankeyChartProps>(SankeyChartBase);
