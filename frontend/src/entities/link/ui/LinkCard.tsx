import { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  CopyIcon,
  ExternalLinkIcon,
  BarChart2Icon,
  TrashIcon,
  QrCodeIcon,
  CalendarIcon,
  CheckIcon,
} from 'lucide-react';
import { cn } from '@/shared/lib/utils';
import { Card } from '@/shared/ui/card';
import { Button } from '@/shared/ui/button';
import { getShortLink } from '@/shared/config';
import { QrCodeModal } from '@/widgets/qr/QrCodeModal';
import type { Link } from '@/shared/api/types';

interface LinkCardProps {
  link: Link;
  onDelete?: (id: string) => void;
  showCampaignInfo?: boolean; // Optional: hide inside campaign details
}

export const LinkCard = ({ link, onDelete }: LinkCardProps) => {
  const [copied, setCopied] = useState(false);
  const [isQrOpen, setIsQrOpen] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(getShortLink(link.slug));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <>
      <Card className="group flex flex-col sm:flex-row gap-4 p-4 items-start sm:items-center bg-card border-border hover:border-primary/50 transition-all duration-200 hover:shadow-sm">
        {/* Icon / Status */}
        <div
          className={cn(
            'w-12 h-12 rounded-xl flex items-center justify-center shrink-0 transition-colors',
            link.is_active
              ? 'bg-primary/10 text-primary'
              : 'bg-destructive/10 text-destructive'
          )}
        >
          <ExternalLinkIcon size={20} />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0 grid gap-1.5">
          <div className="flex items-center gap-2">
            <h4 className="font-bold text-lg truncate text-foreground tracking-tight">
              /{link.slug}
            </h4>
            {!link.is_active && (
              <span className="px-1.5 py-0.5 rounded text-[10px] bg-destructive/10 text-destructive font-bold uppercase tracking-wider">
                Inactive
              </span>
            )}
          </div>

          <a
            href={link.target_url}
            target="_blank"
            rel="noreferrer"
            className="text-sm text-muted-foreground truncate hover:text-foreground transition-colors hover:underline underline-offset-4"
          >
            {link.target_url}
          </a>

          <div className="flex items-center gap-3 mt-1">
             <div className="flex items-center gap-1.5 text-xs text-muted-foreground bg-secondary/50 px-2 py-0.5 rounded-md">
                <CalendarIcon size={12} />
                <span>{new Date(link.created_at).toLocaleDateString()}</span>
             </div>
          </div>
        </div>

        {/* Actions Toolbar */}
        <div className="flex items-center gap-2 w-full sm:w-auto mt-2 sm:mt-0 border-t sm:border-t-0 pt-3 sm:pt-0 border-border justify-end">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsQrOpen(true)}
            title="Show QR Code"
            className="text-muted-foreground hover:text-foreground"
          >
            <QrCodeIcon size={18} />
          </Button>

          <Button
            variant="outline"
            size="sm"
            className={cn("gap-2 min-w-[90px]", copied && "text-green-600 border-green-600/30 bg-green-50")}
            onClick={copyToClipboard}
          >
            {copied ? <CheckIcon size={14} /> : <CopyIcon size={14} />}
            {copied ? 'Copied' : 'Copy'}
          </Button>

          <RouterLink to={`/links/${link.id}`}>
            <Button variant="secondary" size="sm" className="gap-2">
              <BarChart2Icon size={16} />
              <span className="hidden lg:inline">Analytics</span>
            </Button>
          </RouterLink>

          {onDelete && (
            <Button
              variant="ghost"
              size="icon"
              className="text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
              onClick={() => onDelete(link.id)}
            >
              <TrashIcon size={18} />
            </Button>
          )}
        </div>
      </Card>

      <QrCodeModal
        isOpen={isQrOpen}
        onClose={() => setIsQrOpen(false)}
        slug={link.slug}
      />
    </>
  );
};
