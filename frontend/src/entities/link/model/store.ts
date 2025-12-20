import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Link,
  LinksListResponse,
  LinkResponse,
  CreateLinkReq,
  CampaignsListResponse,
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
      // 1. Получаем все кампании
      const campaignsRes = await api.get<CampaignsListResponse>('/campaigns');
      const campaigns = campaignsRes.data.data;

      // 2. Для каждой кампании запрашиваем ссылки
      const linksPromises = campaigns.map((c) =>
        api.get<LinksListResponse>(`/campaigns/${c.id}/links`)
      );

      const responses = await Promise.all(linksPromises);

      // 3. Объединяем все ссылки в один массив
      const allLinks = responses.flatMap((res) => res.data.data);

      // Сортируем по дате создания (новые сверху)
      allLinks.sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );

      set({ links: allLinks });
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
