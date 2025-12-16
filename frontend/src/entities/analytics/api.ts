import { api } from '@/shared/api/base';
import type {
  StackedPoint,
  HierarchyNode,
  SankeyData,
} from '@/shared/api/types';

interface GeoData {
  [countryCode: string]: number;
}
interface HeatmapData {
  [day: string]: { [hour: string]: number };
}
interface RadarData {
  [metric: string]: number;
}

// Тип для параметров запроса
interface AnalyticsParams {
  campaign_id?: string;
  link_id?: string;
  from?: string;
  to?: string;
  [key: string]: unknown; // Разрешаем доп поля, но не any
}

export const analyticsApi = {
  getHierarchy: async () => {
    return api.get<HierarchyNode>('/campaigns/tree');
  },

  getGeoStats: async (params?: AnalyticsParams) => {
    return api.get<{ data: GeoData }>('/analytics/geo', { params });
  },

  getStats: async (dimension: 'browser' | 'os' | 'device') => {
    return api.get<{ [key: string]: number }>('/analytics/stats', {
      params: { dimension },
    });
  },

  getFlow: async (linkId: string) => {
    return api.get<{ data: SankeyData }>('/analytics/flow', {
      params: { link_id: linkId },
    });
  },

  getStream: async (linkId: string) => {
    return api.get<{ data: StackedPoint[] }>('/analytics/stream', {
      params: { link_id: linkId },
    });
  },

  getHeatmap: async (linkId: string) => {
    return api.get<{ data: HeatmapData }>('/analytics/heatmap', {
      params: { link_id: linkId },
    });
  },

  getQuality: async (linkId: string) => {
    return api.get<{ data: RadarData }>('/analytics/quality', {
      params: { link_id: linkId },
    });
  },
};
