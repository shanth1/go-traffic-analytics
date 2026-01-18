import React, { useState } from 'react';
import { PlusIcon, Link2Icon } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { useLinkStore } from '@/entities/link/model/store';
import { useCampaignStore } from '@/entities/campaign/model/store';
import { toast } from '@/entities/notification/store';

type Props = {
  selectedCampaignId?: string;
  onSuccess?: () => void;
};

export const CreateLinkFeature: React.FC<Props> = ({
  selectedCampaignId,
  onSuccess,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [url, setUrl] = useState('');
  const [campaignId, setCampaignId] = useState(selectedCampaignId || '');

  const addLink = useLinkStore((s) => s.addLink);
  const { campaigns, fetchCampaigns } = useCampaignStore();

  // Handler to open modal and initialize state
  const handleOpen = () => {
    // 1. Fetch fresh data
    fetchCampaigns(100, 0);

    // 2. Reset form state based on current props
    setCampaignId(selectedCampaignId || '');
    setUrl('');

    // 3. Open modal
    setIsOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) {
      toast.warn("Validation", "Please enter a valid URL");
      return;
    }

    const targetCampaign = campaignId || (campaigns[0]?.id ?? 'default');

    try {
      await addLink(targetCampaign, url);
      toast.info("Success", "Short link created successfully."); // Success message
      setIsOpen(false);
      if (onSuccess) onSuccess();
    } catch (e) {
      console.error("Error creating link:", e);
    }
  };

  return (
    <>
      <Button
        onClick={handleOpen}
        className="gap-2 shadow-sm"
      >
        <PlusIcon size={18} />
        <span className="hidden sm:inline">New Link</span>
        <span className="sm:hidden">Add</span>
      </Button>

      <ResponsiveSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="New Link"
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">Target URL</label>
            <div className="relative">
              <Link2Icon
                className="absolute left-3 top-3 text-muted-foreground"
                size={16}
              />
              <Input
                className="pl-9 bg-background"
                placeholder="https://your-service.com/..."
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                type="url"
                required
                autoFocus
              />
            </div>
            <p className="text-xs text-muted-foreground">
              Paste the full HTTPS URL where you want users to land.
            </p>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">Campaign</label>
            <select
              className="w-full flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              value={campaignId}
              onChange={(e) => setCampaignId(e.target.value)}
            >
              <option value="">Select Campaign...</option>
              {campaigns.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </div>

          <Button type="submit" className="w-full mt-4">
            Create Link
          </Button>
        </form>
      </ResponsiveSheet>
    </>
  );
};
