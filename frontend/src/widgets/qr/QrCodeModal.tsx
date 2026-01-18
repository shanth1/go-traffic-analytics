import { useState, useRef } from 'react';
import { QRCodeSVG } from 'qrcode.react';
import {
  Download,
  ChevronLeft,
  ChevronRight,
  X,
  QrCode as QrCodeIcon,
} from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { getShortLink } from '@/shared/config';
import { toast } from '@/entities/notification/store';

const QR_STYLES = [
  {
    id: 'classic',
    name: 'Classic Black',
    bg: '#ffffff',
    fg: '#000000',
    wrapperClass: 'bg-white',
  },
  {
    id: 'inverted',
    name: 'Dark Mode',
    bg: '#0f172a',
    fg: '#ffffff',
    wrapperClass: 'bg-slate-950',
  },
  {
    id: 'brand',
    name: 'Brand Primary',
    // Hardcoded colors for QR ensuring scannability/contrast regardless of theme
    bg: '#e0e7ff',
    fg: '#4f46e5',
    wrapperClass: 'bg-indigo-100',
  },
];

interface QrCodeModalProps {
  isOpen: boolean;
  onClose: () => void;
  slug: string;
}

export const QrCodeModal = ({ isOpen, onClose, slug }: QrCodeModalProps) => {
  const [styleIndex, setStyleIndex] = useState(0);
  const svgRef = useRef<HTMLDivElement>(null);
  const linkUrl = getShortLink(slug);
  const currentStyle = QR_STYLES[styleIndex];

  const nextStyle = () => setStyleIndex((prev) => (prev + 1) % QR_STYLES.length);
  const prevStyle = () => setStyleIndex((prev) => (prev - 1 + QR_STYLES.length) % QR_STYLES.length);

  const handleDownload = () => {
    const svg = svgRef.current?.querySelector('svg');
    if (!svg) return;
    const svgData = new XMLSerializer().serializeToString(svg);
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const img = new Image();
    const size = 1024;
    canvas.width = size;
    canvas.height = size;
    const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' });
    const url = URL.createObjectURL(svgBlob);

    img.onload = () => {
      if (ctx) {
        ctx.fillStyle = currentStyle.bg;
        ctx.fillRect(0, 0, size, size);
        ctx.drawImage(img, 0, 0, size, size);
        const pngUrl = canvas.toDataURL('image/png');
        const downloadLink = document.createElement('a');
        downloadLink.href = pngUrl;
        downloadLink.download = `qr-${slug}-${currentStyle.id}.png`;
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);

        toast.info("Saved", "QR Code successfully downloaded.");
      }
      URL.revokeObjectURL(url);
    };
    img.src = url;
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="fixed inset-0 bg-background/80 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      <div className="relative bg-card rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden flex flex-col items-center animate-in fade-in zoom-in-95 duration-200 border border-border">
        <div className="w-full flex justify-between items-center p-4 border-b border-border">
          <h3 className="font-semibold text-lg flex items-center gap-2 text-foreground">
            <QrCodeIcon size={18} className="text-primary" />
            QR Code
          </h3>
          <button
            onClick={onClose}
            className="p-1 hover:bg-accent rounded-full text-muted-foreground"
          >
            <X size={20} />
          </button>
        </div>

        <div className="p-8 w-full flex flex-col items-center gap-6">
          <div className="flex items-center gap-4 w-full justify-between">
            <Button variant="ghost" size="icon" onClick={prevStyle} className="rounded-full h-10 w-10 shrink-0">
              <ChevronLeft />
            </Button>

            <div
              ref={svgRef}
              className={`p-4 rounded-xl shadow-inner transition-colors duration-300 ${currentStyle.wrapperClass}`}
            >
              <QRCodeSVG
                value={linkUrl}
                size={180}
                bgColor={currentStyle.bg}
                fgColor={currentStyle.fg}
                level="M"
              />
            </div>

            <Button variant="ghost" size="icon" onClick={nextStyle} className="rounded-full h-10 w-10 shrink-0">
              <ChevronRight />
            </Button>
          </div>

          <div className="text-center">
            <p className="text-sm font-medium text-foreground">{currentStyle.name}</p>
            <p className="text-xs text-muted-foreground mt-1 truncate max-w-[200px] mx-auto">
              {linkUrl}
            </p>
          </div>
        </div>

        <div className="w-full p-4 border-t border-border bg-muted/30">
          <Button onClick={handleDownload} className="w-full gap-2 shadow-sm">
            <Download size={16} /> Download PNG
          </Button>
        </div>
      </div>
    </div>
  );
};
