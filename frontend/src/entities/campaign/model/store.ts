import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Campaign,
  CampaignsListResponse,
  CampaignResponse,
  CreateCampaignReq,
} from '@/shared/api/types';

interface CampaignStore {
  campaigns: Campaign[];
  isLoading: boolean;
  fetchCampaigns: () => Promise<void>;
  addCampaign: (name: string) => Promise<void>;
}

export const useCampaignStore = create<CampaignStore>((set) => ({
  campaigns: [],
  isLoading: false,
  fetchCampaigns: async () => {
    set({ isLoading: true });
    try {
      const { data } = await api.get<CampaignsListResponse>('/campaigns');
      // Swagger возвращает { data: Campaign[] }
      set({ campaigns: data.data });
    } catch (error) {
      console.error('Failed to fetch campaigns', error);
    } finally {
      set({ isLoading: false });
    }
  },
  addCampaign: async (name) => {
    try {
      const payload: CreateCampaignReq = { name };
      const { data } = await api.post<CampaignResponse>('/campaigns', payload);
      set((state) => ({ campaigns: [data.data, ...state.campaigns] }));
    } catch (error) {
      console.error('Failed to create campaign', error);
      throw error;
    }
  },
}));
