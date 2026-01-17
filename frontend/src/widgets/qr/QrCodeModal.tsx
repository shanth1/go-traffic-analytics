import { useState, useRef } from 'react';
import { QRCodeSVG } from 'qrcode.react';
import {
  Download,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { getShortLink } from '@/shared/config';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';

// --- Styles Definition ---
const QR_STYLES = [
  {
    id: 'classic',
    name: 'Classic Black',
    bg: '#ffffff',
    fg: '#000000',
    wrapperClass: 'bg-white border-slate-200',
  },
  {
    id: 'inverted',
    name: 'Dark Mode',
    bg: '#0f172a', // slate-950
    fg: '#ffffff',
    wrapperClass: 'bg-slate-950 border-slate-800',
  },
  {
    id: 'brand',
    name: 'Brand Indigo',
    bg: '#e0e7ff', // indigo-100
    fg: '#4f46e5', // indigo-600
    wrapperClass: 'bg-indigo-100 border-indigo-200',
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

  const nextStyle = () => {
    setStyleIndex((prev) => (prev + 1) % QR_STYLES.length);
  };

  const prevStyle = () => {
    setStyleIndex((prev) => (prev - 1 + QR_STYLES.length) % QR_STYLES.length);
  };

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

    const svgBlob = new Blob([svgData], {
      type: 'image/svg+xml;charset=utf-8',
    });
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
      }
      URL.revokeObjectURL(url);
    };
    img.src = url;
  };

  return (
    <ResponsiveSheet
      isOpen={isOpen}
      onClose={onClose}
      title="QR Code"
    >
      <div className="flex flex-col items-center w-full max-w-sm mx-auto gap-6 animate-in fade-in duration-300">

        {/* Style Selector & Preview */}
        <div className="flex items-center justify-between w-full gap-4">
          <Button
            variant="outline"
            size="icon"
            onClick={prevStyle}
            className="rounded-full h-12 w-12 shrink-0 border-slate-200 dark:border-slate-700"
          >
            <ChevronLeft size={20} />
          </Button>

          <div
            ref={svgRef}
            className={`p-4 rounded-2xl shadow-sm border transition-all duration-300 ${currentStyle.wrapperClass}`}
          >
            <QRCodeSVG
              value={linkUrl}
              size={160}
              bgColor={currentStyle.bg}
              fgColor={currentStyle.fg}
              level="M"
              includeMargin={false}
            />
          </div>

          <Button
            variant="outline"
            size="icon"
            onClick={nextStyle}
            className="rounded-full h-12 w-12 shrink-0 border-slate-200 dark:border-slate-700"
          >
            <ChevronRight size={20} />
          </Button>
        </div>

        {/* Info Text */}
        <div className="text-center space-y-1">
          <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
            {currentStyle.name}
          </p>
          <p className="text-xs text-slate-500 truncate max-w-[250px] mx-auto select-all">
            {linkUrl}
          </p>
        </div>

        {/* Action Button */}
        <Button
          onClick={handleDownload}
          className="w-full h-12 gap-2 text-base font-medium shadow-lg shadow-indigo-500/20 active:scale-[0.98] transition-all"
        >
          <Download size={18} />
          Download PNG
        </Button>
      </div>
    </ResponsiveSheet>
  );
};
