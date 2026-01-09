import { useEffect } from 'react';
import { FolderIcon, CalendarIcon, ArrowRightIcon } from 'lucide-react';
import { Link } from 'react-router-dom';
import { CreateCampaignFeature } from '@/features/create-campaign/CreateCampaignFeature';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { Card } from '@/shared/ui/card';
import { Pagination } from '@/shared/ui/pagination';
import type { Campaign } from '@/shared/api/types';

interface CampaignCardProps {
  data: Campaign;
}

const CampaignCard = ({ data }: CampaignCardProps) => (
  <Card className="hover:shadow-lg transition-shadow duration-300 group cursor-pointer border-slate-200 dark:border-slate-800 relative overflow-hidden flex flex-col h-full">
    <div className="absolute top-0 left-0 w-1 h-full bg-indigo-500 transform scale-y-0 group-hover:scale-y-100 transition-transform origin-bottom" />
    <div className="p-6 flex flex-col flex-1">
      <div className="flex justify-between items-start mb-4">
        <div className="p-3 bg-indigo-50 dark:bg-indigo-900/20 rounded-lg text-indigo-600 dark:text-indigo-400">
          <FolderIcon size={24} />
        </div>
      </div>

      <h3
        className="text-lg font-bold text-slate-900 dark:text-slate-100 mb-1 line-clamp-1"
        title={data.name}
      >
        {data.name}
      </h3>
      <p className="text-sm text-slate-500 mb-4 flex items-center gap-1">
        <CalendarIcon size={14} />
        {new Date(data.created_at).toLocaleDateString()}
      </p>

      <div className="mt-auto pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
        <div className="flex flex-col">
          <span className="text-xs text-slate-400 uppercase font-bold tracking-wider">
            Status
          </span>
          <span className="text-xs font-medium text-green-600 dark:text-green-400">
            Active
          </span>
        </div>
        <Link
          to={`/links?campaign_id=${data.id}`}
          className="text-sm font-medium text-indigo-600 hover:text-indigo-700 flex items-center gap-1"
        >
          Details <ArrowRightIcon size={14} />
        </Link>
      </div>
    </div>
  </Card>
);

const CampaignsSkeleton = () => (
  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    {[1, 2, 3, 4, 5, 6].map((i) => (
      <div
        key={i}
        className="h-64 bg-slate-200 dark:bg-slate-800 rounded-xl animate-pulse"
      />
    ))}
  </div>
);

export const CampaignsPage = () => {
  const { campaigns, meta, fetchCampaigns, isLoading, setPage } =
    useCampaignStore();

  console.log('CAMP:', campaigns);

  useEffect(() => {
    // Initial fetch only if empty or explicitly needed
    fetchCampaigns();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="space-y-8 animate-in fade-in duration-500 min-h-[calc(100vh-100px)] flex flex-col">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
            Campaigns
          </h1>
          <p className="text-slate-500 mt-1">
            Manage link groups and track their effectiveness.
          </p>
        </div>
        <CreateCampaignFeature />
      </div>

      <div className="flex-1">
        {isLoading && campaigns.length === 0 ? (
          <CampaignsSkeleton />
        ) : (
          <>
            {campaigns.length === 0 ? (
              <div className="text-center py-20 bg-slate-50 dark:bg-slate-900/50 rounded-xl border border-dashed border-slate-300 dark:border-slate-700">
                <FolderIcon className="mx-auto h-12 w-12 text-slate-400 mb-4" />
                <h3 className="text-lg font-medium text-slate-900 dark:text-slate-100">
                  No campaigns yet
                </h3>
                <p className="text-slate-500 mb-6">
                  Create your first campaign to organize your links.
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {campaigns.map((camp) => (
                  <CampaignCard key={camp.id} data={camp} />
                ))}
              </div>
            )}
          </>
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
    </div>
  );
};
