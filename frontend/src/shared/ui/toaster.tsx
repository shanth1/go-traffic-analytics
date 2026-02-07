import { AnimatePresence, motion } from 'framer-motion';
import type { PanInfo } from 'framer-motion';
import { X, CheckCircle2, AlertTriangle, AlertCircle } from 'lucide-react';
import { useToastStore } from '@/entities/notification/store';
import { cn } from '@/shared/lib/utils';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';
import type { Toast } from '@/entities/notification/store';

const icons = {
  info: <CheckCircle2 size={24} />,
  warn: <AlertTriangle size={24} />,
  error: <AlertCircle size={24} />,
  success: <CheckCircle2 size={24} />,
};

const styles = {
  info: 'border-chart-2/20 bg-chart-2/10 text-foreground',
  warn: 'border-chart-3/50 bg-chart-3/10 text-foreground',
  error: 'border-destructive/40 bg-destructive/10 text-foreground',
  success: 'border-primary/20 bg-card/95 text-foreground',
};

const iconStyles = {
  info: 'text-chart-2',
  warn: 'text-chart-3',
  error: 'text-destructive',
  success: 'text-primary',
};

// Animation Variants for "Snappy" feel
const toastVariants = {
  initial: { opacity: 0, y: 20, scale: 0.95 },
  animate: { opacity: 1, y: 0, scale: 1 },
  exit: { opacity: 0, scale: 0.95, transition: { duration: 0.15 } },
};

const ToastItem = ({
  t,
  onRemove,
}: {
  t: Toast;
  onRemove: (id: string) => void;
}) => {
  const type = t.type in styles ? t.type : 'info';

  const handleDragEnd = (_: unknown, info: PanInfo) => {
    // Dismiss if dragged horizontally more than 50px
    if (Math.abs(info.offset.x) > 50) {
      onRemove(t.id);
    }
  };

  return (
    <motion.div
      layout
      variants={toastVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      drag="x"
      dragConstraints={{ left: 0, right: 0 }} // Snap back if not dismissed
      dragElastic={0.7}
      onDragEnd={handleDragEnd}
      className={cn(
        'pointer-events-auto relative flex w-full items-start gap-4 rounded-xl border p-5 shadow-2xl backdrop-blur-2xl cursor-grab active:cursor-grabbing',
        'md:p-4 md:gap-3',
        styles[type as keyof typeof styles]
      )}
    >
      <div
        className={cn(
          'shrink-0 mt-0.5',
          iconStyles[type as keyof typeof iconStyles]
        )}
      >
        {icons[type as keyof typeof icons]}
      </div>

      <div className="flex-1 grid gap-1.5 select-none">
        <h4 className="font-bold text-base leading-none md:text-sm">
          {t.title}
        </h4>
        {t.message && (
          <p className="text-sm text-muted-foreground leading-relaxed md:text-xs">
            {t.message}
          </p>
        )}
      </div>

      <button
        onClick={() => onRemove(t.id)}
        className="shrink-0 -mr-2 -mt-2 p-2 rounded-lg text-muted-foreground opacity-70 transition-opacity hover:bg-secondary hover:text-foreground hover:opacity-100 focus:opacity-100 focus:outline-none"
      >
        <X size={20} className="md:w-4 md:h-4" />
      </button>
    </motion.div>
  );
};

export const Toaster = () => {
  const { toasts, removeToast } = useToastStore();
  const isDesktop = useIsDesktop();

  // LIMIT LOGIC:
  const LIMIT = isDesktop ? 5 : 2;

  // We take the LAST 'LIMIT' items (newest ones)
  const visibleToasts = toasts.slice(-LIMIT);

  return (
    <div
      className={cn(
        'pointer-events-none fixed z-9999 flex flex-col gap-3',
        'top-0 left-0 right-0 p-4',
        'md:top-auto md:bottom-0 md:left-auto md:right-0 md:p-6 md:max-w-sm md:w-full'
      )}
    >
      <AnimatePresence mode="popLayout" initial={false}>
        {visibleToasts.map((t) => (
          <ToastItem key={t.id} t={t} onRemove={removeToast} />
        ))}
      </AnimatePresence>
    </div>
  );
};
