import { useEffect, useState } from 'react';
import { SearchIcon, FilterIcon } from 'lucide-react';
import { useSearchParams } from 'react-router-dom';
import { CreateLinkFeature } from '@/features/create-link/CreateLinkFeature';
import { useLinkStore } from '@/entities/link/model/store';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { Input } from '@/shared/ui/input';
import { Pagination } from '@/shared/ui/pagination';
import { ConfirmationModal } from '@/shared/ui/confirmation-modal';
import { LinkCard } from '@/entities/link/ui/LinkCard'; // New Import
import { toast } from '@/entities/notification/store';

export const LinksPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const campaignIdParam = searchParams.get('campaign_id');

  const {
    links,
    meta,
    setFilters,
    setPage,
    fetchLinks,
    isLoading,
    deleteLink,
  } = useLinkStore();
  const { fetchCampaigns, campaigns } = useCampaignStore();
  const [searchTerm, setSearchTerm] = useState('');

  // Delete State
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    fetchCampaigns();
    setFilters({
      campaign_id: campaignIdParam || undefined,
      search: searchTerm,
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [campaignIdParam]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setFilters({ search: searchTerm });
  };

  const handleCampaignFilterChange = (
    e: React.ChangeEvent<HTMLSelectElement>
  ) => {
    const val = e.target.value;
    if (val) {
      setSearchParams({ campaign_id: val });
    } else {
      setSearchParams({});
    }
    setFilters({ campaign_id: val || undefined });
  };

  const handleDeleteConfirm = async () => {
    if (!deleteId) return;

    setIsDeleting(true);
    try {
      await deleteLink(deleteId);
      toast.info(
        'Link Deleted',
        'The short link has been permanently removed.'
      );
      setDeleteId(null);
    } catch (e) {
      console.error('Failed to delete link', e);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500 flex flex-col min-h-[calc(100vh-100px)]">
      <div className="flex flex-col gap-6">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Links</h1>
            <p className="text-muted-foreground">
              Manage and track your shortened links.
            </p>
          </div>
          {/* Refresh list on success */}
          <CreateLinkFeature onSuccess={() => fetchLinks()} />
        </div>

        {/* Filters Toolbar */}
        <div className="bg-card p-4 rounded-xl border border-border flex flex-col md:flex-row gap-4 shadow-sm">
          <form onSubmit={handleSearch} className="relative flex-1">
            <SearchIcon
              className="absolute left-3 top-2.5 text-muted-foreground"
              size={18}
            />
            <Input
              placeholder="Search by slug or URL..."
              className="pl-10 bg-background"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </form>

          <div className="flex items-center gap-2 w-full md:w-auto">
            <FilterIcon className="text-muted-foreground" size={18} />
            <select
              className="flex h-10 w-full md:w-[200px] items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              value={campaignIdParam || ''}
              onChange={handleCampaignFilterChange}
            >
              <option value="">All Campaigns</option>
              {campaigns.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      <div className="flex-1">
        {isLoading && links.length === 0 ? (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-24 bg-muted/50 rounded-xl animate-pulse"
              />
            ))}
          </div>
        ) : links.length === 0 ? (
          <div className="text-center py-20 text-muted-foreground bg-muted/20 rounded-xl border border-dashed border-border">
            No links found matching your criteria.
          </div>
        ) : (
          <div className="grid gap-4">
            {links.map((link) => (
              <LinkCard key={link.id} link={link} onDelete={setDeleteId} />
            ))}
          </div>
        )}
      </div>

      <div className="mt-auto pt-4 flex justify-center">
        <Pagination
          total={meta.total}
          limit={meta.limit}
          offset={meta.offset}
          onChange={setPage}
        />
      </div>

      <ConfirmationModal
        isOpen={!!deleteId}
        onClose={() => setDeleteId(null)}
        onConfirm={handleDeleteConfirm}
        title="Delete Link?"
        description="Are you sure? This short link will stop working immediately, and users will see a 404 error. All analytics data for this link will be permanently removed."
        confirmLabel="Yes, Delete Link"
        isLoading={isDeleting}
      />
    </div>
  );
};
