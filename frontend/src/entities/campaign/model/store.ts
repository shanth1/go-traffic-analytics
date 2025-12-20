import { create } from 'zustand';
import type { Campaign } from '@/shared/api/types';

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
    // Mock API delay
    await new Promise((r) => setTimeout(r, 800));
    set({
      isLoading: false,
      campaigns: [
        {
          id: '1',
          name: 'Summer Sale 2024',
          user_id: 'u1',
          created_at: '2024-05-01T10:00:00Z',
        },
        {
          id: '2',
          name: 'Instagram Bio',
          user_id: 'u1',
          created_at: '2024-05-05T12:30:00Z',
        },
        {
          id: '3',
          name: 'Tech Blog Promo',
          user_id: 'u1',
          created_at: '2024-06-01T09:15:00Z',
        },
        {
          id: '4',
          name: 'Q3 Report Share',
          user_id: 'u1',
          created_at: '2024-07-10T14:20:00Z',
        },
      ],
    });
  },
  addCampaign: async (name) => {
    // Mock Create
    const newCamp: Campaign = {
      id: Math.random().toString(36).substr(2, 9),
      name,
      user_id: 'u1',
      created_at: new Date().toISOString(),
    };
    set((state) => ({ campaigns: [newCamp, ...state.campaigns] }));
  },
}));
