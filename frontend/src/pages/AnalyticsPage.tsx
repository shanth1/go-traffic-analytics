import { motion } from 'framer-motion';
import { StreamGraphWidget } from '@/widgets/charts/StreamGraph';
import { RadarWidget } from '@/widgets/charts/RadarChart';

export const AnalyticsPage = () => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div>
        <h1 className="text-2xl font-bold text-white">Deep Analytics</h1>
        <p className="text-slate-400 text-sm">Real-time traffic monitoring</p>
      </div>

      {/* Stream Graph Card */}
      <div className="bg-surface border border-slate-800 rounded-2xl p-4 shadow-xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-sm font-medium text-slate-300">
            Traffic Stream (Devices)
          </h3>
          <span className="text-xs px-2 py-1 bg-slate-800 rounded text-slate-400">
            7 Days
          </span>
        </div>
        <StreamGraphWidget />
      </div>

      {/* Grid for smaller charts */}
      <div className="grid grid-cols-2 gap-4">
        <div className="col-span-2 bg-surface border border-slate-800 rounded-2xl p-4 shadow-xl">
          <h3 className="text-sm font-medium text-slate-300 mb-2">
            Traffic Quality Score
          </h3>
          <RadarWidget />
        </div>
      </div>
    </motion.div>
  );
};
