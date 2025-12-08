import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCampaignStore } from '@/entities/campaign/store';

export const CreatePage = () => {
  const [name, setName] = useState('');
  const addCampaign = useCampaignStore((s) => s.addCampaign);
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name) return;
    addCampaign(name);
    navigate('/campaigns');
  };

  return (
    <div className="pt-10">
      <h1 className="text-2xl font-bold mb-6">New Campaign</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm text-slate-400 mb-2">
            Campaign Name
          </label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full bg-slate-900 border border-slate-800 rounded-lg p-3 text-white focus:outline-none focus:border-blue-500 transition-colors"
            placeholder="e.g. Black Friday 2024"
            autoFocus
          />
        </div>
        <button className="w-full bg-blue-600 hover:bg-blue-500 text-white font-medium py-3 rounded-lg transition-colors shadow-lg shadow-blue-900/20">
          Create & Generate Links
        </button>
      </form>
    </div>
  );
};
