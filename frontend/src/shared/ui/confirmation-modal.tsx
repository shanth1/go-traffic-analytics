import React from 'react';
import { ResponsiveSheet } from './responsive-sheet';
import { Button } from './button';
import { AlertTriangleIcon } from 'lucide-react';
import { useIsDesktop } from '@/shared/lib/hooks/useMediaQuery';

interface ConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  isLoading?: boolean;
}

export const ConfirmationModal: React.FC<ConfirmationModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  description,
  confirmLabel = 'Delete',
  cancelLabel = 'Cancel',
  isLoading = false,
}) => {
  const isDesktop = useIsDesktop();

  return (
    <ResponsiveSheet isOpen={isOpen} onClose={onClose} title="Confirm Action">
      <div className="flex flex-col gap-6 pt-2 pb-2">
        {/* Warning Icon & Text */}
        <div className="flex flex-col items-center text-center gap-4">
          <div className="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center text-destructive">
            <AlertTriangleIcon size={32} />
          </div>
          <div>
            <h3 className="text-xl font-bold text-foreground mb-2">{title}</h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              {description}
            </p>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col-reverse sm:flex-row gap-3 mt-2">
          {/* Cancel button: Visible ONLY on Desktop */}
          {isDesktop && (
            <Button
              variant="outline"
              className="flex-1 h-12 rounded-xl text-base"
              onClick={onClose}
              disabled={isLoading}
            >
              {cancelLabel}
            </Button>
          )}

          <Button
            variant="destructive"
            className="flex-1 h-12 rounded-xl text-base font-bold shadow-lg shadow-destructive/20 w-full"
            onClick={onConfirm}
            disabled={isLoading}
          >
            {isLoading ? 'Deleting...' : confirmLabel}
          </Button>
        </div>
      </div>
    </ResponsiveSheet>
  );
};
