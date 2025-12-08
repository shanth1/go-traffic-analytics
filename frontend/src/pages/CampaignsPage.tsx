import { useCampaignStore } from '@/entities/campaign/store';
import { motion } from 'framer-motion';
import { MousePointer2 } from 'lucide-react';

export const CampaignsPage = () => {
  const { campaigns } = useCampaignStore();

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="space-y-4"
    >
      <div>
        <h1 className="text-2xl font-bold">My Campaigns</h1>
        <p className="text-slate-400 text-sm">Active links overview</p>
      </div>

      <div className="space-y-3">
        {campaigns.map((c) => (
          <div
            key={c.id}
            className="bg-surface border border-slate-800 rounded-xl p-4 flex justify-between items-center group active:scale-95 transition-transform"
          >
            <div>
              <h3 className="font-semibold text-slate-200">{c.name}</h3>
              <div className="flex items-center gap-1 text-xs text-slate-500 mt-1">
                <span
                  className={
                    c.status === 'active'
                      ? 'text-emerald-500'
                      : 'text-amber-500'
                  }
                >
                  ● {c.status.toUpperCase()}
                </span>
                <span>• {new Date().toLocaleDateString()}</span>
              </div>
            </div>
            <div className="text-right">
              <div className="text-xl font-bold text-white flex items-center justify-end gap-1">
                {c.clicks.toLocaleString()}{' '}
                <MousePointer2 size={14} className="text-slate-500" />
              </div>
              <div className="text-xs text-slate-500">total clicks</div>
            </div>
          </div>
        ))}
      </div>
    </motion.div>
  );
};
