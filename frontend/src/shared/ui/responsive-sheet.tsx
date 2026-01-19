import React, { useEffect } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import type { PanInfo } from 'framer-motion';
import { X } from 'lucide-react';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';
import { useScrollLock } from '@/shared/lib/hooks/useScrollLock';

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
    className="fixed inset-0 bg-black/60 backdrop-blur-sm z-90"
  />
);

export const ResponsiveSheet: React.FC<ResponsiveSheetProps> = ({
  isOpen,
  onClose,
  title,
  children,
}) => {
  const isDesktop = useIsDesktop();
  useScrollLock(isOpen);

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  const onDragEnd = (_: unknown, info: PanInfo) => {
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
            /* --- Desktop --- */
            <div className="fixed inset-0 z-100 flex items-center justify-center pointer-events-none p-4">
              <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 10 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 10 }}
                transition={{ duration: 0.2, ease: 'easeOut' }}
                className="w-full max-w-lg bg-card rounded-2xl shadow-2xl border border-border pointer-events-auto overflow-hidden flex flex-col max-h-[90vh]"
              >
                <div className="flex items-center justify-between p-6 border-b border-border shrink-0">
                  <h2 className="text-xl font-semibold text-foreground">
                    {title}
                  </h2>
                  <button
                    onClick={onClose}
                    className="p-2 hover:bg-secondary rounded-full transition-colors text-muted-foreground hover:text-foreground"
                  >
                    <X size={20} />
                  </button>
                </div>
                <div className="p-6 text-foreground overflow-y-auto">
                  {children}
                </div>
              </motion.div>
            </div>
          ) : (
            /* --- Mobile Drawer --- */
            <motion.div
              className="fixed bottom-0 left-0 right-0 z-100 outline-none flex flex-col"
              style={{
                maxHeight: '92vh',
              }}
              initial={{ y: '100%' }}
              animate={{ y: 0 }}
              exit={{ y: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 300 }}
              drag="y"
              dragConstraints={{ top: 0 }}
              dragElastic={0.1}
              onDragEnd={onDragEnd}
            >
              {/* FIXED CLASSNAME: removed backtick and used arbitrary value for strong rounding */}
              <div className="bg-card rounded-t-4xl shadow-[0_-8px_30px_rgba(0,0,0,0.3)] border-t border-border flex flex-col pb-[50vh] -mb-[50vh]">
                {/* Handle Bar */}
                <div className="pt-4 pb-2 flex justify-center cursor-grab active:cursor-grabbing touch-none shrink-0">
                  <div className="w-12 h-1.5 bg-muted-foreground/20 rounded-full" />
                </div>

                {/* Header */}
                <div className="px-6 pb-4 text-center shrink-0">
                  <h2 className="text-lg font-bold text-foreground">{title}</h2>
                </div>

                {/* Content */}
                <div
                  className="px-6 pb-10 overflow-y-auto overscroll-contain max-h-[80vh]"
                  onPointerDown={(e) => e.stopPropagation()}
                >
                  {children}
                </div>
              </div>
            </motion.div>
          )}
        </>
      )}
    </AnimatePresence>
  );
};
