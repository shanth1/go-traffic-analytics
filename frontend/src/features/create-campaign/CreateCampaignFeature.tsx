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
        className="gap-2 bg-indigo-600 hover:bg-indigo-700 text-white"
      >
        <PlusIcon size={18} />
        <span className="hidden sm:inline">Новая кампания</span>
        <span className="sm:hidden">Создать</span>
      </Button>

      <ResponsiveSheet
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Создать кампанию"
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
              Название кампании
            </label>
            <Input
              placeholder="Например: Summer Sale 2024"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
            <p className="text-xs text-slate-500">
              Кампании помогают группировать ссылки для общей аналитики.
            </p>
          </div>
          <Button type="submit" className="w-full bg-indigo-600 text-white">
            Создать
          </Button>
        </form>
      </ResponsiveSheet>
    </>
  );
};
