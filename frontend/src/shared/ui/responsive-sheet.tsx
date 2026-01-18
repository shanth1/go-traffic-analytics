import React from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { X } from 'lucide-react';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';

interface ResponsiveSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  trigger?: React.ReactNode;
}

const Backdrop = ({ onClose }: { onClose: () => void }) => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
    onClick={onClose}
    className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50"
  />
);

export const ResponsiveSheet: React.FC<ResponsiveSheetProps> = ({
  isOpen,
  onClose,
  title,
  children,
}) => {
  const isDesktop = useIsDesktop();

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <Backdrop onClose={onClose} />

          {isDesktop ? (
            /* --- Desktop Modal --- */
            <div className="fixed inset-0 z-50 flex items-center justify-center pointer-events-none">
              <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 10 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 10 }}
                className="w-full max-w-lg bg-background rounded-xl shadow-2xl border border-border pointer-events-auto overflow-hidden"
              >
                <div className="flex items-center justify-between p-6 border-b border-border">
                  <h2 className="text-xl font-semibold text-foreground">{title}</h2>
                  <button
                    onClick={onClose}
                    className="p-2 hover:bg-accent rounded-full transition-colors text-muted-foreground hover:text-accent-foreground"
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
              initial={{ y: '100%' }}
              animate={{ y: 0 }}
              exit={{ y: '100%' }}
              transition={{ type: 'spring', damping: 25, stiffness: 300 }}
              className="fixed bottom-0 left-0 right-0 z-50 bg-background rounded-t-[20px] shadow-[0_-5px_20px_rgba(0,0,0,0.1)] border-t border-border max-h-[85vh] overflow-y-auto"
            >
              <div className="sticky top-0 bg-background z-10 pt-4 pb-2 px-6 flex flex-col items-center">
                <div className="w-12 h-1.5 bg-muted rounded-full mb-4" />
                <div className="flex items-center justify-between w-full">
                  <h2 className="text-lg font-semibold text-foreground">{title}</h2>
                  <button
                    onClick={onClose}
                    className="p-2 -mr-2 text-muted-foreground"
                  >
                    <X size={20} />
                  </button>
                </div>
              </div>
              <div className="p-6 pt-2 pb-10 safe-area-bottom text-foreground">{children}</div>
            </motion.div>
          )}
        </>
      )}
    </AnimatePresence>
  );
};
