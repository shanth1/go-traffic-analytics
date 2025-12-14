import { create } from 'zustand';
import type { Link } from '@/shared/api/types';

interface LinkStore {
  links: Link[];
  isLoading: boolean;
  fetchLinks: () => Promise<void>;
  addLink: (campaignId: string, url: string, slug?: string) => Promise<void>;
  deleteLink: (id: string) => void;
}

export const useLinkStore = create<LinkStore>((set) => ({
  links: [],
  isLoading: false,
  fetchLinks: async () => {
    set({ isLoading: true });
    await new Promise((r) => setTimeout(r, 600));
    set({
      isLoading: false,
      links: [
        {
          id: 'l1',
          campaign_id: '1',
          user_id: 'u1',
          slug: 'summer-sale',
          target_url: 'https://shop.com/sale',
          is_active: true,
          created_at: '2024-05-01',
        },
        {
          id: 'l2',
          campaign_id: '2',
          user_id: 'u1',
          slug: 'my-insta',
          target_url: 'https://instagram.com/me',
          is_active: true,
          created_at: '2024-05-06',
        },
        {
          id: 'l3',
          campaign_id: '1',
          user_id: 'u1',
          slug: 'discount-50',
          target_url: 'https://shop.com/discount',
          is_active: false,
          created_at: '2024-05-02',
        },
      ],
    });
  },
  addLink: async (campaignId, url, slug) => {
    const newLink: Link = {
      id: Math.random().toString(36).substr(2, 9),
      campaign_id: campaignId,
      user_id: 'u1',
      target_url: url,
      slug: slug || Math.random().toString(36).substr(2, 6),
      is_active: true,
      created_at: new Date().toISOString(),
    };
    set((state) => ({ links: [newLink, ...state.links] }));
  },
  deleteLink: (id) =>
    set((state) => ({ links: state.links.filter((l) => l.id !== id) })),
}));
