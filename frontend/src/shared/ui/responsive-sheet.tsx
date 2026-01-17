import { createPortal } from 'react-dom';
import { AnimatePresence, motion } from 'framer-motion';
import type { PanInfo } from 'framer-motion';
import { X } from 'lucide-react';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';
import { useScrollLock } from '@/shared/lib/hooks/useScrollLock';
import { cn } from '@/shared/lib/utils';
import { useState } from 'react';

interface ResponsiveSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}

const Backdrop = ({ onClose }: { onClose: () => void }) => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
    onClick={onClose}
    className="fixed inset-0 bg-black/60 z-99 backdrop-blur-[2px]"
  />
);

export const ResponsiveSheet: React.FC<ResponsiveSheetProps> = ({
  isOpen,
  onClose,
  title,
  children,
}) => {
  const isDesktop = useIsDesktop();
  const [mounted] = useState(true);

  // Lock scroll when open
  useScrollLock(isOpen);

  // No effect needed; mounted state initialized to true

  // Handle Drag on Mobile (Swipe down to close)
  const onDragEnd = (_: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => {
    if (info.offset.y > 100) {
      onClose();
    }
  };

  if (!mounted) return null;

  // Use Portal to ensure the modal is always on top of everything (including fixed navbars)
  return createPortal(
    <AnimatePresence>
      {isOpen && (
        <>
          <Backdrop onClose={onClose} />

          {isDesktop ? (
            /* --- Desktop Modal (Centered) --- */
            <div className="fixed inset-0 z-100 flex items-center justify-center pointer-events-none p-4">
              <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 10 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 10 }}
                transition={{ duration: 0.2 }}
                className="w-full max-w-lg bg-white dark:bg-slate-950 rounded-2xl shadow-2xl border border-slate-200 dark:border-slate-800 pointer-events-auto flex flex-col max-h-[85vh]"
              >
                <div className="flex items-center justify-between p-5 border-b border-slate-100 dark:border-slate-800 shrink-0">
                  <h2 className="text-xl font-bold tracking-tight">{title}</h2>
                  <button
                    onClick={onClose}
                    className="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors text-slate-500"
                  >
                    <X size={20} />
                  </button>
                </div>
                <div className="p-6 overflow-y-auto custom-scrollbar">
                    {children}
                </div>
              </motion.div>
            </div>
          ) : (
            /* --- Mobile Drawer (iOS Style) --- */
            <div className="fixed inset-x-0 bottom-0 z-100 flex justify-center pointer-events-none">
              <motion.div
                  initial={{ y: '100%' }}
                  animate={{ y: 0 }}
                  exit={{ y: '100%' }}
                  transition={{ type: 'spring', damping: 25, stiffness: 300 }}
                  drag="y"
                  dragConstraints={{ top: 0 }}
                  dragElastic={0.05} // Resistance when pulling up
                  onDragEnd={onDragEnd}
                  className={cn(
                    "w-full bg-white dark:bg-slate-950 pointer-events-auto",
                    "rounded-t-3xl shadow-[0_-8px_30px_rgba(0,0,0,0.12)]",
                    "border-t border-white/10 dark:border-slate-800",
                    "flex flex-col max-h-[92vh]" // iOS feeling: leaves a bit of space at top
                  )}
                >
                {/* Drag Handle Area */}
                <div className="shrink-0 pt-4 pb-2 px-6 flex flex-col items-center cursor-grab active:cursor-grabbing touch-none">
                  <div className="w-10 h-1.5 bg-slate-300 dark:bg-slate-700/50 rounded-full mb-4" />

                  {/* Header without X button on mobile (cleaner, gesture based) */}
                  <div className="w-full text-center pb-2">
                    <h2 className="text-lg font-bold text-slate-900 dark:text-slate-100">
                        {title}
                    </h2>
                  </div>
                </div>

                {/* Content Area */}
                <div className="p-6 pt-0 overflow-y-auto overscroll-contain pb-safe-area-bottom">
                  {children}
                </div>
              </motion.div>
            </div>
          )}
        </>
      )}
    </AnimatePresence>,
    document.body
  );
};
