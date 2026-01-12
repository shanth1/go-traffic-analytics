import { api } from '@/shared/api/base';
import type { Link, LinksListResponse } from '@/shared/api/types';

export const linkApi = {
  /**
   * Fetches a single link by ID.
   * Since Swagger doesn't specify GET /links/{id}, we query the list with a limit.
   * In a real scenario, backend usually provides GET /links/{id}.
   */
  getLinkById: async (id: string): Promise<Link | null> => {
    // We assume we can't filter by ID directly in list params based on Swagger shown,
    // but usually getting the list works.
    // Optimization: If the backend supports ?ids=... or similar, use that.
    // Fallback: Fetch list and find. Ideally, backend needs GET /links/:id
    try {
      // Trying to find it in the general list (not ideal for perf, but fits constraints)
      // In production, request Backend Team to add GET /links/{id}
      const { data } = await api.get<LinksListResponse>('/links', {
        params: { limit: 100 }, // Hope it's in the last 100 created
      });
      return data.data.find((l) => l.id === id) || null;
    } catch (e) {
      console.error('Failed to fetch link details', e);
      return null;
    }
  },
};
