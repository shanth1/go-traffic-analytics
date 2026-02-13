import { api } from '@/shared/api/base';
import type {
  StreamGraphResponse,
  GeoResponse,
  SankeyResponse,
  AnalyticsSummaryResponse,
  QualityResponse,
  TreeResponse,
  StatsResponse,
  HeatmapResponse,
} from '@/shared/api/types';

export interface AnalyticsParams {
  campaign_id?: string;
  link_id?: string;
  from?: string; // RFC3339
  to?: string; // RFC3339
  [key: string]: unknown;
}

export const analyticsApi = {
  getHierarchy: async () => {
    const { data } = await api.get<TreeResponse>('/users/tree');
    return data.data;
  },

  getCountryStats: async (params?: AnalyticsParams) => {
    const { data } = await api.get<GeoResponse>('/analytics/geo', { params });
    return data.data;
  },

  getStats: async (
    dimension: 'browser' | 'os' | 'device',
    params?: AnalyticsParams
  ) => {
    const { data } = await api.get<StatsResponse>('/analytics/stats', {
      params: { dimension, ...params },
    });
    return data.data;
  },

  getSummary: async (params?: AnalyticsParams) => {
    const { data } = await api.get<AnalyticsSummaryResponse>(
      '/analytics/summary',
      { params }
    );
    return data.data;
  },

  getFlow: async (linkId?: string, params?: AnalyticsParams) => {
    const { data } = await api.get<SankeyResponse>('/analytics/flow', {
      params: { link_id: linkId, ...params },
    });
    return data.data;
  },

  getStream: async (linkId?: string, params?: AnalyticsParams) => {
    const { data } = await api.get<StreamGraphResponse>('/analytics/stream', {
      params: { link_id: linkId, ...params },
    });
    return data.data;
  },

  getHeatmap: async (linkId?: string, params?: AnalyticsParams) => {
    const { data } = await api.get<HeatmapResponse>('/analytics/heatmap', {
      params: { link_id: linkId, ...params },
    });
    return data.data;
  },

  getQuality: async (linkId?: string, params?: AnalyticsParams) => {
    const { data } = await api.get<QualityResponse>('/analytics/quality', {
      params: { link_id: linkId, ...params },
    });
    return data.data;
  },

  exportExcel: async (params: {
    from?: string;
    to?: string;
    campaign_id?: string;
  }) => {
    const response = await api.get('/analytics/export', {
      params,
      responseType: 'blob',
    });
    return response.data;
  },
};
