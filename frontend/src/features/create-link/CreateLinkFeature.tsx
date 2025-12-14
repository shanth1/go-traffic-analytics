import React, { useState } from 'react';
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
  const campaigns = useCampaignStore((s) => s.campaigns); // Для селекта

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;
    // Если кампания не выбрана, берем первую или дефолтную (в реале валидация)
    const targetCampaign = campaignId || (campaigns[0]?.id ?? 'default');
    await addLink(targetCampaign, url, slug);
    setUrl('');
    setSlug('');
    setIsOpen(false);
  };

  return (
    <>
      <Button
        onClick={() => setIsOpen(true)}
        className="gap-2 bg-indigo-600 hover:bg-indigo-700 text-white"
      >
        <PlusIcon size={18} />
        <span className="hidden sm:inline">Сократить ссылку</span>
        <span className="sm:hidden">Ссылка</span>
      </Button>

      <ResponsiveSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Новая ссылка"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Целевой URL</label>
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
              Свой алиас (необязательно)
            </label>
            <div className="flex items-center gap-2">
              <span className="text-slate-400 text-sm bg-slate-100 px-2 py-2 rounded-md border border-slate-200">
                domain.com/
              </span>
              <Input
                placeholder="my-super-link"
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Кампания</label>
            <select
              className="w-full flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              value={campaignId}
              onChange={(e) => setCampaignId(e.target.value)}
            >
              <option value="">Без кампании</option>
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
            Создать ссылку
          </Button>
        </form>
      </ResponsiveSheet>
    </>
  );
};
