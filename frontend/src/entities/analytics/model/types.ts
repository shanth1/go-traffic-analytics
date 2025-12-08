export interface StackedPoint {
  time: string; // ISO String
  values: Record<string, number>; // { iOS: 10, Android: 5 }
}

export interface SankeyNode {
  id: string;
  layer: number;
}

export interface SankeyLink {
  source: string;
  target: string;
  value: number;
}

export interface SankeyData {
  nodes: SankeyNode[];
  links: SankeyLink[];
}

export interface QualityData {
  high_quality: number;
  suspicious: number;
  bot: number;
}
