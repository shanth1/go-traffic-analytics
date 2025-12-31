import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Link,
  LinksListResponse,
  LinkResponse,
  CreateLinkReq,
} from '@/shared/api/types';

interface LinkStore {
  links: Link[];
  isLoading: boolean;
  fetchLinks: () => Promise<void>;
  addLink: (campaignId: string, url: string, slug?: string) => Promise<void>;
  deleteLink: (id: string) => Promise<void>;
}

export const useLinkStore = create<LinkStore>((set) => ({
  links: [],
  isLoading: false,
  fetchLinks: async () => {
    set({ isLoading: true });
    try {
      const { data } = await api.get<LinksListResponse>('/links');

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
