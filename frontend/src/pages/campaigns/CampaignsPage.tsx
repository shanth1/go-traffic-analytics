import { useEffect } from 'react';
import { FolderIcon, CalendarIcon, ArrowRightIcon } from 'lucide-react';
import { Link } from 'react-router-dom';
import { CreateCampaignFeature } from '@/features/create-campaign/CreateCampaignFeature';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { Card } from '@/shared/ui/card';
import type { Campaign } from '@/shared/api/types';

interface CampaignCardProps {
  data: Campaign;
}

const CampaignCard = ({ data }: CampaignCardProps) => (
  <Card className="hover:shadow-lg transition-shadow duration-300 group cursor-pointer border-slate-200 dark:border-slate-800 relative overflow-hidden">
    <div className="absolute top-0 left-0 w-1 h-full bg-indigo-500 transform scale-y-0 group-hover:scale-y-100 transition-transform origin-bottom" />
    <div className="p-6">
      <div className="flex justify-between items-start mb-4">
        <div className="p-3 bg-indigo-50 dark:bg-indigo-900/20 rounded-lg text-indigo-600 dark:text-indigo-400">
          <FolderIcon size={24} />
        </div>
        <div className="flex space-x-1 items-end h-8">
          {[40, 70, 45, 90, 60].map((h, i) => (
            <div
              key={i}
              className="w-1 bg-slate-200 dark:bg-slate-700 rounded-t-sm"
              style={{ height: `${h}%` }}
            />
          ))}
        </div>
      </div>

      <h3 className="text-lg font-bold text-slate-900 dark:text-slate-100 mb-1">
        {data.name}
      </h3>
      <p className="text-sm text-slate-500 mb-4 flex items-center gap-1">
        <CalendarIcon size={14} />
        {new Date(data.created_at).toLocaleDateString()}
      </p>

      <div className="flex items-center justify-between mt-4 pt-4 border-t border-slate-100 dark:border-slate-800">
        <div className="flex flex-col">
          <span className="text-xs text-slate-400 uppercase font-bold tracking-wider">
            Clicks
          </span>
          <span className="font-mono font-semibold text-slate-700 dark:text-slate-300">
            1,234
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

export const CampaignsPage = () => {
  const { campaigns, fetchCampaigns, isLoading } = useCampaignStore();

  useEffect(() => {
    fetchCampaigns();
  }, [fetchCampaigns]);

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
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

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-48 bg-slate-200 dark:bg-slate-800 rounded-xl animate-pulse"
            />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          {campaigns.map((camp) => (
            <CampaignCard key={camp.id} data={camp} />
          ))}
        </div>
      )}
    </div>
  );
};
