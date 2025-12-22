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

export interface Plan {
  id: string;
  name: string;
  max_links: number;
  max_clicks_month: number;
  price_cents: number;
  is_active: boolean;
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

export interface CategoryStat {
  name: string;
  value: number;
  share?: number;
}

export interface AnalyticsSummary {
  total_clicks: number;
  top_browsers: CategoryStat[];
  top_os: CategoryStat[];
}

export interface TrafficQuality {
  bot_score: number;
  geo_diversity_score: number;
  human_score: number;
  is_suspicious: boolean;
  mobile_friendly_score: number;
}

// --- Requests ---
export interface LoginReq {
  email: string;
  password?: string;
}

export interface CreateCampaignReq {
  name: string;
}

export interface CreateLinkReq {
  campaign_id: string;
  target_url: string;
}

// --- Responses ---
export interface AuthResponse {
  token: string;
  user: User;
}

export interface ResponseWrapper<T> {
  data: T;
}

export interface GeoPoint {
  country: string;
  value: number;
}

export type CampaignsListResponse = ResponseWrapper<Campaign[]>;
export type CampaignResponse = ResponseWrapper<Campaign>;
export type LinksListResponse = ResponseWrapper<Link[]>;
export type LinkResponse = ResponseWrapper<Link>;
export type AnalyticsSummaryResponse = ResponseWrapper<AnalyticsSummary>;
export type GeoResponse = ResponseWrapper<GeoPoint[]>;
export type TreeResponse = ResponseWrapper<HierarchyNode>;
export type StreamGraphResponse = ResponseWrapper<StackedPoint[]>;
export type SankeyResponse = ResponseWrapper<SankeyData>;
export type QualityResponse = ResponseWrapper<TrafficQuality>;
