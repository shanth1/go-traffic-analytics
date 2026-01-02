import { useEffect } from 'react';
import {
  CopyIcon,
  ExternalLinkIcon,
  BarChart2Icon,
  TrashIcon,
} from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { CreateLinkFeature } from '@/features/create-link/CreateLinkFeature';
import { useLinkStore } from '@/entities/link/model/store';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { cn } from '@/shared/lib/utils';
import { Card } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import type { Link as LinkType } from '@/shared/api/types';

interface LinkPageProps {
  data: LinkType;
}

const LinkCard = ({ data }: LinkPageProps) => {
  const { deleteLink } = useLinkStore();

  const copyToClipboard = () => {
    navigator.clipboard.writeText(`http://short.ly/${data.slug}`);
  };

  return (
    <Card className="p-4 flex flex-col sm:flex-row gap-4 items-start sm:items-center bg-white dark:bg-slate-950 hover:border-indigo-300 transition-colors">
      <div
        className={cn(
          'w-12 h-12 rounded-full flex items-center justify-center shrink-0',
          data.is_active
            ? 'bg-green-100 text-green-600'
            : 'bg-red-100 text-red-600'
        )}
      >
        <ExternalLinkIcon size={20} />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <h4 className="font-bold text-lg truncate">/{data.slug}</h4>
          {!data.is_active && (
            <span className="px-2 py-0.5 rounded text-[10px] bg-red-100 text-red-600 font-bold uppercase">
              Inactive
            </span>
          )}
        </div>
        <p className="text-sm text-slate-500 truncate">{data.target_url}</p>
        <div className="flex items-center gap-4 mt-2 text-xs text-slate-400">
          <span>{new Date(data.created_at).toLocaleDateString()}</span>
          <span className="flex items-center gap-1">
            <BarChart2Icon size={12} /> 245 clicks
          </span>
        </div>
      </div>

      <div className="flex items-center gap-2 w-full sm:w-auto mt-2 sm:mt-0 border-t sm:border-t-0 pt-3 sm:pt-0">
        <Button
          variant="outline"
          size="sm"
          className="flex-1 sm:flex-none gap-2"
          onClick={copyToClipboard}
        >
          <CopyIcon size={14} /> Copy
        </Button>
        <Link to={`/links/${data.id}`}>
          <Button variant="secondary" size="sm" className="flex-1 sm:flex-none">
            Analytics
          </Button>
        </Link>
        <Button
          variant="ghost"
          size="icon"
          className="text-red-400 hover:text-red-600"
          onClick={() => deleteLink(data.id)}
        >
          <TrashIcon size={16} />
        </Button>
      </div>
    </Card>
  );
};

export const LinksPage = () => {
  const [searchParams] = useSearchParams();
  const campaignId = searchParams.get('campaign_id');

  const { links, fetchLinks, isLoading } = useLinkStore();
  const { fetchCampaigns } = useCampaignStore();

  useEffect(() => {
    fetchLinks(campaignId ? { campaign_id: campaignId } : undefined);
    fetchCampaigns();
  }, [fetchLinks, fetchCampaigns, campaignId]);

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
            Ссылки{campaignId ? ` кампании ${campaignId}` : ''}
          </h1>
          <p className="text-slate-500">
            {campaignId
              ? 'Ссылки выбранной кампании.'
              : 'Все ваши сокращенные ссылки в одном месте.'}
          </p>
        </div>
        <CreateLinkFeature />
      </div>

      <div className="grid gap-4">
        {isLoading
          ? [1, 2, 3, 4].map((i) => (
              <div
                key={i}
                className="h-24 bg-slate-100 dark:bg-slate-900 rounded-lg animate-pulse"
              />
            ))
          : links.map((link) => <LinkCard key={link.id} data={link} />)}
      </div>
    </div>
  );
};
