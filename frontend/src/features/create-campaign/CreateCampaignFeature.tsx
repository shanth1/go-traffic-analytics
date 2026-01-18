import React, { useState } from 'react';
import { PlusIcon } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { useCampaignStore } from '@/entities/campaign/model/store';

export const CreateCampaignFeature = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState('');
  const addCampaign = useCampaignStore((s) => s.addCampaign);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name) return;
    await addCampaign(name);
    setName('');
    setIsOpen(false);
  };

  return (
    <>
      <Button
        onClick={() => setIsOpen(true)}
        className="gap-2"
      >
        <PlusIcon size={18} />
        <span className="hidden sm:inline">New Campaign</span>
        <span className="sm:hidden">Create</span>
      </Button>

      <ResponsiveSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Create Campaign"
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">
              Campaign Name
            </label>
            <Input
              placeholder="e.g. Summer Sale 2024"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <p className="text-xs text-muted-foreground">
              Campaigns help group links for overall analytics.
            </p>
          </div>
          <Button type="submit" className="w-full">
            Create
          </Button>
        </form>
      </ResponsiveSheet>
    </>
  );
};
