import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Link,
  LinksListResponse,
  LinkResponse,
  CreateLinkReq,
} from '@/shared/api/types';

interface LinkFilters {
  campaign_id?: string;
  search?: string;
  is_active?: boolean;
  limit?: number;
  offset?: number;
}

interface LinkStore {
  links: Link[];
  isLoading: boolean;
  fetchLinks: (filters?: LinkFilters) => Promise<void>;
  addLink: (campaignId: string, url: string, slug?: string) => Promise<void>;
  deleteLink: (id: string) => Promise<void>;
}

export const useLinkStore = create<LinkStore>((set) => ({
  links: [],
  isLoading: false,
  fetchLinks: async (filters?: LinkFilters) => {
    set({ isLoading: true });
    try {
      const params = new URLSearchParams();
      if (filters?.campaign_id)
        params.append('campaign_id', filters.campaign_id);
      if (filters?.search) params.append('search', filters.search);
      if (filters?.is_active !== undefined)
        params.append('is_active', filters.is_active.toString());
      if (filters?.limit) params.append('limit', filters.limit.toString());
      if (filters?.offset) params.append('offset', filters.offset.toString());

      const url = `/links${params.toString() ? '?' + params.toString() : ''}`;
      const { data } = await api.get<LinksListResponse>(url);

      // Сортируем по дате создания (новые сверху)
      data.data.sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );

      set({ links: data.data });
    } catch (error) {
      console.error('Failed to fetch links', error);
    } finally {
      set({ isLoading: false });
    }
  },
  addLink: async (campaignId, url) => {
    try {
      const payload: CreateLinkReq = {
        campaign_id: campaignId,
        target_url: url,
      };
      const { data } = await api.post<LinkResponse>('/links', payload);
      set((state) => ({ links: [data.data, ...state.links] }));
    } catch (error) {
      console.error('Failed to create link', error);
      throw error;
    }
  },
  deleteLink: async (id) => {
    try {
      await api.delete(`/links/${id}`);
      set((state) => ({ links: state.links.filter((l) => l.id !== id) }));
    } catch (error) {
      console.error('Failed to delete link', error);
    }
  },
}));
