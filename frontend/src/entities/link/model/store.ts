import { create } from 'zustand';
import { api } from '@/shared/api/base';
import type {
  Link,
  LinksListResponse,
  LinkResponse,
  CreateLinkReq,
  PaginationMeta,
} from '@/shared/api/types';

interface LinkFilters {
  campaign_id?: string;
  search?: string;
  is_active?: boolean;
}

interface LinkState {
  links: Link[];
  meta: PaginationMeta;
  filters: LinkFilters;
  isLoading: boolean;

  // Actions
  fetchLinks: (
    filters?: Partial<LinkFilters>,
    offset?: number
  ) => Promise<void>;
  addLink: (campaignId: string, url: string, slug?: string) => Promise<void>;
  deleteLink: (id: string) => Promise<void>;
  setPage: (offset: number) => void;
  setFilters: (filters: Partial<LinkFilters>) => void;
}

export const useLinkStore = create<LinkState>((set, get) => ({
  links: [],
  meta: {
    total: 0,
    limit: 10,
    offset: 0,
  },
  filters: {},
  isLoading: false,

  fetchLinks: async (newFilters = {}, offset) => {
    set({ isLoading: true });

    const state = get();
    // Merge existing filters with new ones
    const filters = { ...state.filters, ...newFilters };
    const reqOffset = offset ?? state.meta.offset;
    const reqLimit = state.meta.limit;

    try {
      const params = new URLSearchParams();

      if (filters.campaign_id)
        params.append('campaign_id', filters.campaign_id);
      if (filters.search) params.append('search', filters.search);
      if (filters.is_active !== undefined)
        params.append('is_active', filters.is_active.toString());

      params.append('limit', reqLimit.toString());
      params.append('offset', reqOffset.toString());

      const url = `/links${params.toString() ? '?' + params.toString() : ''}`;
      const { data } = await api.get<LinksListResponse>(url);

      set({
        links: data.data,
        meta: data.meta,
        filters: filters, // Save merged filters
        isLoading: false,
      });
    } catch (error) {
      console.error('Failed to fetch links', error);
      set({ isLoading: false });
    }
  },

  addLink: async (campaignId, url) => {
    try {
      const payload: CreateLinkReq = {
        campaign_id: campaignId,
        target_url: url,
        // custom_slug: slug // Assuming backend supports this in body if implemented
      };
      await api.post<LinkResponse>('/links', payload);
      // Refetch current page to see update or go to first page
      await get().fetchLinks(undefined, 0);
    } catch (error) {
      console.error('Failed to create link', error);
      throw error;
    }
  },

  deleteLink: async (id) => {
    try {
      await api.delete(`/links/${id}`);
      // Refetch to update list and total count correctly
      await get().fetchLinks();
    } catch (error) {
      console.error('Failed to delete link', error);
    }
  },

  setPage: (offset: number) => {
    get().fetchLinks(undefined, offset);
  },

  setFilters: (filters: Partial<LinkFilters>) => {
    // When filters change, reset to page 1 (offset 0)
    get().fetchLinks(filters, 0);
  },
}));
