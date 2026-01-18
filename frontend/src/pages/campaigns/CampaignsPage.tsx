import { useEffect } from 'react';
import { FolderIcon, CalendarIcon, LayoutDashboardIcon } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { CreateCampaignFeature } from '@/features/create-campaign/CreateCampaignFeature';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { Card } from '@/shared/ui/card';
import { Pagination } from '@/shared/ui/pagination';
import type { Campaign } from '@/shared/api/types';
import { Button } from '@/shared/ui/button';

interface CampaignCardProps {
  data: Campaign;
}

const CampaignCard = ({ data }: CampaignCardProps) => {
  const navigate = useNavigate();

  return (
    <Card
      className="group hover:shadow-lg transition-all duration-300 border-border relative overflow-hidden flex flex-col h-full cursor-pointer"
      onClick={() => navigate(`/campaigns/${data.id}`)}
    >
      <div className="absolute top-0 left-0 w-1 h-full bg-primary transform scale-y-0 group-hover:scale-y-100 transition-transform origin-bottom duration-300" />

      <div className="p-6 flex flex-col flex-1">
        <div className="flex justify-between items-start mb-4">
          <div className="p-3 bg-primary/10 rounded-lg text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
            <FolderIcon size={24} />
          </div>
        </div>

        <h3
          className="text-lg font-bold text-foreground mb-2 line-clamp-1 group-hover:text-primary transition-colors"
          title={data.name}
        >
          {data.name}
        </h3>

        <p className="text-sm text-muted-foreground mb-6 flex items-center gap-1.5">
          <CalendarIcon size={14} />
          <span>Created {new Date(data.created_at).toLocaleDateString()}</span>
        </p>

        <div className="mt-auto pt-4 border-t border-border flex items-center justify-between">
          <div className="flex flex-col">
            <span className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">
              Status
            </span>
            <span className="text-xs font-medium text-chart-2">
              Active
            </span>
          </div>

          <Button
            variant="secondary"
            size="sm"
            className="gap-2 pointer-events-none group-hover:bg-primary/10 group-hover:text-primary"
          >
            <LayoutDashboardIcon size={14} /> Dashboard
          </Button>
        </div>
      </div>
    </Card>
  );
};

const CampaignsSkeleton = () => (
  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    {[1, 2, 3, 4, 5, 6].map((i) => (
      <div
        key={i}
        className="h-64 bg-muted rounded-xl animate-pulse"
      />
    ))}
  </div>
);

export const CampaignsPage = () => {
  const { campaigns, meta, fetchCampaigns, isLoading, setPage } =
    useCampaignStore();

  useEffect(() => {
    fetchCampaigns();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="space-y-8 animate-in fade-in duration-500 min-h-[calc(100vh-100px)] flex flex-col">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">
            Campaigns
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage your marketing campaigns and track their aggregated
            performance.
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
              <div className="text-center py-20 bg-muted/20 rounded-xl border border-dashed border-border">
                <FolderIcon className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium text-foreground">
                  No campaigns yet
                </h3>
                <p className="text-muted-foreground mb-6">
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
