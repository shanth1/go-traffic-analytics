import { useMemo } from 'react';
import { Group } from '@visx/group';
import { Cluster, hierarchy } from '@visx/hierarchy';
import { LinkVertical } from '@visx/shape';
import { withParentSize } from '@visx/responsive';
import type {
  HierarchyPointNode,
  HierarchyPointLink,
} from '@visx/hierarchy/lib/types';
import type { HierarchyNode } from '@/shared/api/types';
import { THEME_COLORS, CHART_COLORS } from '@/shared/config/theme';

interface TreeProps {
  parentWidth?: number;
  parentHeight?: number;
  data: HierarchyNode;
}

const HierarchyTreeBase = ({
  parentWidth = 0,
  parentHeight = 0,
  data,
}: TreeProps) => {
  const width = parentWidth;
  const height = parentHeight;
  const root = useMemo(() => hierarchy<HierarchyNode>(data), [data]);

  return (
    <svg width={width} height={height}>
      {width > 10 && (
        <Cluster<HierarchyNode> root={root} size={[width - 40, height - 100]}>
          {(rootNode) => {
            const nodes = rootNode.descendants();
            const links = rootNode.links();

            return (
              <Group top={40} left={20}>
                {links.map(
                  (link: HierarchyPointLink<HierarchyNode>, i: number) => (
                    <LinkVertical
                      key={`link-${i}`}
                      data={link}
                      stroke={THEME_COLORS.border}
                      strokeWidth="1"
                      fill="none"
                    />
                  )
                )}
                {nodes.map(
                  (node: HierarchyPointNode<HierarchyNode>, i: number) => {
                    const isRoot = node.depth === 0;
                    const isLink = !node.children;
                    const fill = isRoot
                      ? CHART_COLORS[0]
                      : isLink
                        ? CHART_COLORS[3]
                        : CHART_COLORS[4];

                    return (
                      <Group top={node.y} left={node.x} key={`node-${i}`}>
                        <circle
                          r={isRoot ? 8 : isLink ? 4 : 6}
                          fill={fill}
                          stroke={THEME_COLORS.background}
                          strokeWidth={2}
                        />
                        <text
                          dy={isLink ? 5 : -15}
                          dx={isLink ? 8 : 0}
                          fontSize={10}
                          textAnchor={isLink ? 'start' : 'middle'}
                          fill={THEME_COLORS.foreground}
                          className="font-medium opacity-60"
                        >
                          {node.data.name}
                        </text>
                      </Group>
                    );
                  }
                )}
              </Group>
            );
          }}
        </Cluster>
      )}
    </svg>
  );
};

export const HierarchyTree = withParentSize<TreeProps>(HierarchyTreeBase);
