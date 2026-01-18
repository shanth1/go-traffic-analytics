import { useEffect, useState } from 'react';
import {
  CopyIcon,
  ExternalLinkIcon,
  BarChart2Icon,
  TrashIcon,
  SearchIcon,
  FilterIcon,
  QrCodeIcon,
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
import { ConfirmationModal } from '@/shared/ui/confirmation-modal'; // Импорт модалки
import type { Link as LinkType } from '@/shared/api/types';
import { QrCodeModal } from '@/widgets/qr/QrCodeModal';
import { getShortLink } from '@/shared/config';

interface LinkCardProps {
  data: LinkType;
  onDeleteClick: (id: string) => void; // Прокидываем ID наверх
}

const LinkCard = ({ data, onDeleteClick }: LinkCardProps) => {
  const [copied, setCopied] = useState(false);
  const [isQrOpen, setIsQrOpen] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(getShortLink(data.slug));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <>
      <Card className="p-4 flex flex-col sm:flex-row gap-4 items-start sm:items-center bg-card hover:border-primary/50 transition-colors">
        <div
          className={cn(
            'w-12 h-12 rounded-full flex items-center justify-center shrink-0',
            data.is_active
              ? 'bg-chart-2/10 text-chart-2'
              : 'bg-destructive/10 text-destructive'
          )}
        >
          <ExternalLinkIcon size={20} />
        </div>

        <div className="flex-1 min-w-0 grid gap-1">
          <div className="flex items-center gap-2">
            <h4 className="font-bold text-lg truncate text-foreground">
              /{data.slug}
            </h4>
            {!data.is_active && (
              <span className="px-2 py-0.5 rounded text-[10px] bg-destructive/10 text-destructive font-bold uppercase">
                Inactive
              </span>
            )}
          </div>
          <p className="text-sm text-muted-foreground truncate">{data.target_url}</p>
          <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
            <span>{new Date(data.created_at).toLocaleDateString()}</span>
          </div>
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto mt-2 sm:mt-0 border-t sm:border-t-0 pt-3 sm:pt-0 border-border">
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
            className="text-destructive hover:text-destructive hover:bg-destructive/10"
            // Здесь мы больше не удаляем сразу, а вызываем обработчик родителя
            onClick={() => onDeleteClick(data.id)}
          >
            <TrashIcon size={16} />
          </Button>
        </div>
      </Card>

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

  const { links, meta, setFilters, setPage, isLoading, deleteLink } = useLinkStore();
  const { fetchCampaigns, campaigns } = useCampaignStore();
  const [searchTerm, setSearchTerm] = useState('');

  // Состояние для удаления
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

  // Обработчик подтверждения удаления
  const handleDeleteConfirm = async () => {
    if (!deleteId) return;

    setIsDeleting(true);
    try {
      await deleteLink(deleteId);
      setDeleteId(null); // Закрываем модалку при успехе
    } catch (e) {
      console.error("Failed to delete link", e);
      // Опционально: можно добавить toast с ошибкой здесь
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500 flex flex-col min-h-[calc(100vh-100px)]">
      <div className="flex flex-col gap-6">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              Links
            </h1>
            <p className="text-muted-foreground">
              Manage and track your shortened links.
            </p>
          </div>
          <CreateLinkFeature />
        </div>

        {/* Filters Toolbar */}
        <div className="bg-card p-4 rounded-xl border border-border flex flex-col md:flex-row gap-4">
          <form onSubmit={handleSearch} className="relative flex-1">
            <SearchIcon
              className="absolute left-3 top-2.5 text-muted-foreground"
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
            <FilterIcon className="text-muted-foreground" size={18} />
            <select
              className="flex h-10 w-full md:w-[200px] items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
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
                className="h-24 bg-muted rounded-lg animate-pulse"
              />
            ))}
          </div>
        ) : links.length === 0 ? (
          <div className="text-center py-20 text-muted-foreground">
            No links found matching your criteria.
          </div>
        ) : (
          <div className="grid gap-4">
            {links.map((link) => (
              <LinkCard
                key={link.id}
                data={link}
                onDeleteClick={setDeleteId} // Передаем сеттер ID
              />
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

      {/* Модалка подтверждения */}
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
