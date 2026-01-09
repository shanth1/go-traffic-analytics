import React, { useState, useEffect } from 'react';
import { PlusIcon, Link2Icon } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { useLinkStore } from '@/entities/link/model/store';
import { useCampaignStore } from '@/entities/campaign/model/store';

export const CreateLinkFeature = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [url, setUrl] = useState('');
  const [slug, setSlug] = useState('');
  const [campaignId, setCampaignId] = useState('');

  const addLink = useLinkStore((s) => s.addLink);
  const { campaigns, fetchCampaigns } = useCampaignStore();

  useEffect(() => {
    if (isOpen) {
      // Fetch list for dropdown (fetching page 1 usually enough for dropdown in simple cases,
      // ideally should have specific endpoint for dropdowns or infinite scroll)
      fetchCampaigns(100, 0);
    }
  }, [isOpen, fetchCampaigns]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;
    const targetCampaign = campaignId || (campaigns[0]?.id ?? 'default');
    await addLink(targetCampaign, url, slug);
    setUrl('');
    setSlug('');
    setIsOpen(false);
  };

  // ... (rest of the component remains similar, ensure Select uses campaigns array)

  return (
    <>
      <Button
        onClick={() => setIsOpen(true)}
        className="gap-2 bg-indigo-600 hover:bg-indigo-700 text-white"
      >
        <PlusIcon size={18} />
        <span className="hidden sm:inline">Shorten Link</span>
        <span className="sm:hidden">Link</span>
      </Button>

      <ResponsiveSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="New Link"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Target URL</label>
            <div className="relative">
              <Link2Icon
                className="absolute left-3 top-3 text-slate-400"
                size={16}
              />
              <Input
                className="pl-9"
                placeholder="https://very-long-url.com/..."
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                type="url"
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">
              Custom Alias (optional)
            </label>
            <div className="flex items-center gap-2">
              <span className="text-slate-400 text-sm bg-slate-100 px-2 py-2 rounded-md border border-slate-200">
                /
              </span>
              <Input
                placeholder="my-super-link"
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Campaign</label>
            <select
              className="w-full flex h-10 rounded-md border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
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

          <Button
            type="submit"
            className="w-full bg-indigo-600 text-white mt-4"
          >
            Create Link
          </Button>
        </form>
      </ResponsiveSheet>
    </>
  );
};
