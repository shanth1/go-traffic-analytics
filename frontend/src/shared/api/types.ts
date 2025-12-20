// --- Domain Models ---
export interface User {
  id: string;
  email: string;
  role: 'admin' | 'client';
  is_active: boolean;
  clicks_current_month: number;
  plan_id: string;
  created_at: string;
}

export interface Campaign {
  id: string;
  name: string;
  user_id: string;
  created_at: string;
}

export interface Link {
  id: string;
  campaign_id: string;
  user_id: string;
  slug: string;
  target_url: string;
  is_active: boolean;
  created_at: string;
}

// --- Analytics Models ---
export interface StackedPoint {
  time: string;
  values: Record<string, number>;
}

export interface HierarchyNode {
  name: string;
  type: 'root' | 'campaign' | 'link';
  value?: number;
  children?: HierarchyNode[];
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

// --- Responses ---
export interface AuthResponse {
  token: string;
  user: User;
}
