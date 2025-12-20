import { useMemo } from 'react';
import { Group } from '@visx/group';
import { Cluster, hierarchy } from '@visx/hierarchy';
import { LinkVertical } from '@visx/shape';
import { withParentSize } from '@visx/responsive';
import type {
  HierarchyPointNode,
  HierarchyPointLink,
} from '@visx/hierarchy/lib/types'; // Импорт типов
import type { HierarchyNode } from '@/shared/api/types';

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

  // !!! Исправление 3 (Hooks): Перемещаем проверку width ниже хуков,
  // но т.к. Cluster сам внутри использует хуки, просто вернем null внутри SVG или
  // позволим visx обработать это.
  // Лучше всего вернуть null в самом низу рендера, если width 0.

  return (
    <svg width={width} height={height}>
      {width > 10 && (
        <Cluster<HierarchyNode> root={root} size={[width - 40, height - 100]}>
          {(rootNode) => {
            // Исправление 2: Cluster передает сам rootNode (HierarchyPointNode), а не объект
            const nodes = rootNode.descendants();
            const links = rootNode.links();

            return (
              <Group top={40} left={20}>
                {links.map(
                  (link: HierarchyPointLink<HierarchyNode>, i: number) => (
                    <LinkVertical
                      key={`link-${i}`}
                      data={link}
                      stroke="#cbd5e1"
                      strokeWidth="1"
                      fill="none"
                    />
                  )
                )}
                {nodes.map(
                  (node: HierarchyPointNode<HierarchyNode>, i: number) => {
                    const isRoot = node.depth === 0;
                    const isLink = !node.children;

                    return (
                      <Group top={node.y} left={node.x} key={`node-${i}`}>
                        <circle
                          r={isRoot ? 8 : isLink ? 4 : 6}
                          fill={
                            isRoot ? '#4f46e5' : isLink ? '#ec4899' : '#0ea5e9'
                          }
                          stroke="white"
                          strokeWidth={2}
                        />
                        <text
                          dy={isLink ? 5 : -15}
                          dx={isLink ? 8 : 0}
                          fontSize={10}
                          textAnchor={isLink ? 'start' : 'middle'}
                          fill="#64748b"
                          className="font-medium"
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
