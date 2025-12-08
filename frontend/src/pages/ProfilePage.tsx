import { motion } from 'framer-motion';

export const ProfilePage = () => {
  return (
    <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
      <div className="bg-linear-to-br from-slate-900 to-slate-800 rounded-2xl p-6 border border-slate-700 text-center mb-6">
        <div className="w-20 h-20 bg-linear-to-tr from-blue-500 to-purple-600 rounded-full mx-auto mb-4 border-4 border-slate-950" />
        <h2 className="text-xl font-bold">Alex Johnson</h2>
        <p className="text-slate-400 text-sm">Pro Plan • Expires in 12 days</p>
      </div>

      <div className="grid grid-cols-2 gap-3 mb-6">
        <div className="bg-surface p-4 rounded-xl border border-slate-800">
          <div className="text-3xl font-bold text-blue-500">12.5k</div>
          <div className="text-xs text-slate-400">Total Clicks</div>
        </div>
        <div className="bg-surface p-4 rounded-xl border border-slate-800">
          <div className="text-3xl font-bold text-emerald-500">85%</div>
          <div className="text-xs text-slate-400">Quality Score</div>
        </div>
      </div>
    </motion.div>
  );
};
