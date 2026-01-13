import { useEffect, useState } from 'react';
import {
  CopyIcon,
  ExternalLinkIcon,
  BarChart2Icon,
  TrashIcon,
  SearchIcon,
  FilterIcon,
  QrCodeIcon, // [ADDED]
} from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { CreateLinkFeature } from '@/features/create-link/CreateLinkFeature';
import { useLinkStore } from '@/entities/link/model/store';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { cn } from '@/shared/lib/utils';
import { Card } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Pagination } from '@/shared/ui/pagination';
import type { Link as LinkType } from '@/shared/api/types';
import { QrCodeModal } from '@/widgets/qr/QrCodeModal'; // [ADDED]
import { getShortLink } from '@/shared/config'; // [ADDED]

interface LinkCardProps {
  data: LinkType;
}

const LinkCard = ({ data }: LinkCardProps) => {
  const { deleteLink } = useLinkStore();
  const [copied, setCopied] = useState(false);

  // [ADDED] State for QR Modal
  const [isQrOpen, setIsQrOpen] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(getShortLink(data.slug));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <>
      <Card className="p-4 flex flex-col sm:flex-row gap-4 items-start sm:items-center bg-white dark:bg-slate-950 hover:border-indigo-300 transition-colors">
        <div
          className={cn(
            'w-12 h-12 rounded-full flex items-center justify-center shrink-0',
            data.is_active
              ? 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400'
              : 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400'
          )}
        >
          <ExternalLinkIcon size={20} />
        </div>

        <div className="flex-1 min-w-0 grid gap-1">
          <div className="flex items-center gap-2">
            <h4 className="font-bold text-lg truncate text-slate-900 dark:text-slate-100">
              /{data.slug}
            </h4>
            {!data.is_active && (
              <span className="px-2 py-0.5 rounded text-[10px] bg-red-100 text-red-600 font-bold uppercase">
                Inactive
              </span>
            )}
          </div>
          <p className="text-sm text-slate-500 truncate">{data.target_url}</p>
          <div className="flex items-center gap-4 mt-1 text-xs text-slate-400">
            <span>{new Date(data.created_at).toLocaleDateString()}</span>
          </div>
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto mt-2 sm:mt-0 border-t sm:border-t-0 pt-3 sm:pt-0">
          {/* [ADDED] QR Button */}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsQrOpen(true)}
            title="Show QR Code"
          >
            <QrCodeIcon size={16} />
          </Button>

          <Button
            variant="outline"
            size="sm"
            className="flex-1 sm:flex-none gap-2"
            onClick={copyToClipboard}
          >
            <CopyIcon size={14} /> {copied ? 'Copied!' : 'Copy'}
          </Button>
          <Link to={`/links/${data.id}`} className="flex-1 sm:flex-none">
            <Button variant="secondary" size="sm" className="w-full gap-2">
              <BarChart2Icon size={14} /> Analytics
            </Button>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            className="text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
            onClick={() => deleteLink(data.id)}
          >
            <TrashIcon size={16} />
          </Button>
        </div>
      </Card>

      {/* [ADDED] Modal */}
      <QrCodeModal
        isOpen={isQrOpen}
        onClose={() => setIsQrOpen(false)}
        slug={data.slug}
      />
    </>
  );
};

export const LinksPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const campaignIdParam = searchParams.get('campaign_id');

  const { links, meta, setFilters, setPage, isLoading } = useLinkStore();
  const { fetchCampaigns, campaigns } = useCampaignStore();

  const [searchTerm, setSearchTerm] = useState('');

  // Initial Load & URL Sync
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

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500 flex flex-col min-h-[calc(100vh-100px)]">
      <div className="flex flex-col gap-6">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
              Links
            </h1>
            <p className="text-slate-500">
              Manage and track your shortened links.
            </p>
          </div>
          <CreateLinkFeature />
        </div>

        {/* Filters Toolbar */}
        <div className="bg-white dark:bg-slate-950 p-4 rounded-xl border border-slate-200 dark:border-slate-800 flex flex-col md:flex-row gap-4">
          <form onSubmit={handleSearch} className="relative flex-1">
            <SearchIcon
              className="absolute left-3 top-2.5 text-slate-400"
              size={18}
            />
            <Input
              placeholder="Search by slug or URL..."
              className="pl-10"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </form>

          <div className="flex items-center gap-2 w-full md:w-auto">
            <FilterIcon className="text-slate-400" size={18} />
            <select
              className="flex h-10 w-full md:w-[200px] items-center justify-between rounded-md border border-slate-200 bg-white px-3 py-2 text-sm placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-800 dark:bg-slate-950 dark:ring-offset-slate-950 dark:placeholder:text-slate-400 dark:focus:ring-indigo-500"
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
                className="h-24 bg-slate-100 dark:bg-slate-900 rounded-lg animate-pulse"
              />
            ))}
          </div>
        ) : links.length === 0 ? (
          <div className="text-center py-20 text-slate-500">
            No links found matching your criteria.
          </div>
        ) : (
          <div className="grid gap-4">
            {links.map((link) => (
              <LinkCard key={link.id} data={link} />
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
    </div>
  );
};
