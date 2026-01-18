import React from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import type { PanInfo } from 'framer-motion';
import { X } from 'lucide-react';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';

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
    className="fixed inset-0 bg-black/60 backdrop-blur-sm z-99"
  />
);

export const ResponsiveSheet: React.FC<ResponsiveSheetProps> = ({
  isOpen,
  onClose,
  title,
  children,
}) => {
  const isDesktop = useIsDesktop();

  const onDragEnd = (_: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => {
    if (info.offset.y > 100 || info.velocity.y > 500) {
      onClose();
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <Backdrop onClose={onClose} />

          {isDesktop ? (
            /* --- Desktop Modal --- */
            <div className="fixed inset-0 z-100 flex items-center justify-center pointer-events-none">
              <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 10 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 10 }}
                transition={{ duration: 0.2 }}
                className="w-full max-w-lg bg-card rounded-2xl shadow-2xl border border-border pointer-events-auto overflow-hidden"
              >
                <div className="flex items-center justify-between p-6 border-b border-border">
                  <h2 className="text-xl font-semibold text-foreground">{title}</h2>
                  <button
                    onClick={onClose}
                    className="p-2 hover:bg-secondary rounded-full transition-colors text-muted-foreground hover:text-foreground"
                  >
                    <X size={20} />
                  </button>
                </div>
                <div className="p-6 text-foreground">{children}</div>
              </motion.div>
            </div>
          ) : (
            /* --- Mobile Drawer (Bottom Sheet) --- */
            <motion.div
              className="fixed bottom-0 left-0 right-0 z-100 bg-card rounded-t-3xl shadow-[0_-8px_30px_rgba(0,0,0,0.2)] border-t border-border outline-none"

              // Animation config
              initial={{ y: '100%' }}
              animate={{ y: 0 }}
              exit={{ y: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 300 }}

              // Drag config
              drag="y"
              dragConstraints={{ top: 0 }}
              dragElastic={0.2}
              onDragEnd={onDragEnd}
            >
              {/* Handle Bar */}
              <div className="pt-4 pb-2 flex justify-center cursor-grab active:cursor-grabbing touch-none">
                <div className="w-12 h-1.5 bg-muted rounded-full opacity-50" />
              </div>

              {/* Header  */}
              <div className="px-6 pb-4 text-center">
                <h2 className="text-lg font-bold text-foreground">{title}</h2>
              </div>

              {/* Content Container */}
              <div className="px-6 pb-10 max-h-[80vh] overflow-y-auto overscroll-contain">
                {children}
              </div>
            </motion.div>
          )}
        </>
      )}
    </AnimatePresence>
  );
};
