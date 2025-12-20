import { api } from '@/shared/api/base';
import type {
  HierarchyNode,
  StreamGraphResponse,
  GeoResponse,
  SankeyResponse,
  AnalyticsSummaryResponse,
  QualityResponse,
} from '@/shared/api/types';

interface AnalyticsParams {
  campaign_id?: string;
  link_id?: string;
  from?: string;
  to?: string;
  [key: string]: unknown;
}

export const analyticsApi = {
  getHierarchy: async () => {
    const { data } = await api.get<HierarchyNode>('/campaigns/tree');
    return data;
  },

  getGeoStats: async (params?: AnalyticsParams) => {
    const { data } = await api.get<GeoResponse>('/analytics/geo', { params });
    return data.data;
  },

  getStats: async (dimension: 'browser' | 'os' | 'device') => {
    const { data } = await api.get<Record<string, number>>('/analytics/stats', {
      params: { dimension },
    });
    return data;
  },

  getSummary: async () => {
    const { data } =
      await api.get<AnalyticsSummaryResponse>('/analytics/summary');
    return data.data;
  },

  getFlow: async (linkId: string) => {
    const { data } = await api.get<SankeyResponse>('/analytics/flow', {
      params: { link_id: linkId },
    });
    return data.data;
  },

  getStream: async (linkId: string) => {
    const { data } = await api.get<StreamGraphResponse>('/analytics/stream', {
      params: { link_id: linkId },
    });
    return data.data;
  },

  getHeatmap: async (linkId: string) => {
    // Для heatmap тип ответа сложнее (nested object), пока оставляем inline типизацию или Record
    const { data } = await api.get<{
      [day: string]: { [hour: string]: number };
    }>('/analytics/heatmap', { params: { link_id: linkId } });
    return data;
  },

  getQuality: async (linkId: string) => {
    const { data } = await api.get<QualityResponse>('/analytics/quality', {
      params: { link_id: linkId },
    });
    return data.data;
  },
};
