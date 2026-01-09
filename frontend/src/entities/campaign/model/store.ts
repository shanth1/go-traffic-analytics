import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Campaign,
  CampaignsListResponse,
  CampaignResponse,
  CreateCampaignReq,
  PaginationMeta,
} from '@/shared/api/types';

interface CampaignState {
  campaigns: Campaign[];
  meta: PaginationMeta;
  isLoading: boolean;

  // Actions
  fetchCampaigns: (limit?: number, offset?: number) => Promise<void>;
  addCampaign: (name: string) => Promise<void>;
  setPage: (offset: number) => void;
}

export const useCampaignStore = create<CampaignState>((set, get) => ({
  campaigns: [],
  meta: {
    total: 0,
    limit: 9, // Grid layout: 3x3 works well
    offset: 0,
  },
  isLoading: false,

  fetchCampaigns: async (limit, offset) => {
    set({ isLoading: true });

    // Use current state if params not provided
    const currentMeta = get().meta;
    const reqLimit = limit ?? currentMeta.limit;
    const reqOffset = offset ?? currentMeta.offset;

    try {
      const { data } = await api.get<CampaignsListResponse>('/campaigns', {
        params: {
          limit: reqLimit,
          offset: reqOffset,
        },
      });

      set({
        campaigns: data.data,
        meta: data.meta,
        isLoading: false,
      });
    } catch (error) {
      console.error('Failed to fetch campaigns', error);
      set({ isLoading: false });
    }
  },

  addCampaign: async (name) => {
    try {
      const payload: CreateCampaignReq = { name };
      await api.post<CampaignResponse>('/campaigns', payload);

      await get().fetchCampaigns(get().meta.limit, 0);
    } catch (error) {
      console.error('Failed to create campaign', error);
      throw error;
    }
  },

  setPage: (offset: number) => {
    get().fetchCampaigns(undefined, offset);
  },
}));
