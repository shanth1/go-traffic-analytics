// --- Domain Models ---
export interface User {
  id: string;
  email: string;
  role: string;
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

export interface StreamChartData {
  time: Date;
  [key: string]: number | Date;
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

export interface HeatmapPoint {
  day: string;
  hour: string;
  count: number;
}

export interface CountryStat {
  country: string;
  value: number;
}

export interface CityStat {
  city: string;
  value: number;
}

export interface GeoStats {
  countries: CountryStat[];
  cities: CityStat[];
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

// --- Responses & Metadata ---

export interface PaginationMeta {
  total: number;
  limit: number;
  offset: number;
}

export interface ResponseWrapper<T> {
  data: T;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

export interface LoginData {
  token: string;
  user: User;
}

// Typed Responses
export type AuthResponse = ResponseWrapper<LoginData>;
export type CampaignResponse = ResponseWrapper<Campaign>;
export type LinkResponse = ResponseWrapper<Link>;

// List Responses
export type CampaignsListResponse = PaginatedResponse<Campaign>;
export type LinksListResponse = PaginatedResponse<Link>;

// Analytics Responses
export type AnalyticsSummaryResponse = ResponseWrapper<AnalyticsSummary>;
export type GeoResponse = ResponseWrapper<GeoStats>;
export type TreeResponse = ResponseWrapper<HierarchyNode>;
export type StreamGraphResponse = ResponseWrapper<StackedPoint[]>;
export type SankeyResponse = ResponseWrapper<SankeyData>;
export type QualityResponse = ResponseWrapper<TrafficQuality>;
export type StatsResponse = ResponseWrapper<CategoryStat[]>;
export type HeatmapResponse = ResponseWrapper<HeatmapPoint[]>;
