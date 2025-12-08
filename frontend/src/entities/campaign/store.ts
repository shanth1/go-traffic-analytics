import { create } from 'zustand';

export interface Campaign {
  id: string;
  name: string;
  clicks: number;
  status: 'active' | 'paused';
}

interface CampaignStore {
  campaigns: Campaign[];
  addCampaign: (name: string) => void;
}

// MOCK DATA для демонстрации
export const useCampaignStore = create<CampaignStore>((set) => ({
  campaigns: [
    { id: '1', name: 'Summer Sale 2024', clicks: 1240, status: 'active' },
    { id: '2', name: 'Instagram Promo', clicks: 850, status: 'paused' },
    { id: '3', name: 'Tech Blog Referral', clicks: 3200, status: 'active' },
  ],
  addCampaign: (name) =>
    set((state) => ({
      campaigns: [
        ...state.campaigns,
        {
          id: Math.random().toString(36).substr(2, 9),
          name,
          clicks: 0,
          status: 'active',
        },
      ],
    })),
}));
