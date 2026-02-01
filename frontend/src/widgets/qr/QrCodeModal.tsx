import { useState, useRef } from 'react';
import { QRCodeSVG } from 'qrcode.react';
import { Download, ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { ResponsiveSheet } from '@/shared/ui/responsive-sheet';
import { getShortLink } from '@/shared/config';

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
    bg: '#000000',
    fg: '#ffffff',
    wrapperClass: 'bg-black border-slate-800',
  },
  {
    id: 'brand',
    name: 'Brand Primary',
    bg: '#e0e7ff',
    fg: '#4f46e5',
    wrapperClass: 'bg-indigo-100 border-indigo-200',
  },
  {
    id: 'brand',
    name: 'Brand Primary',
    bg: '#e0e7ff',
    fg: '#4f46e5',
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

  const nextStyle = () =>
    setStyleIndex((prev) => (prev + 1) % QR_STYLES.length);
  const prevStyle = () =>
    setStyleIndex((prev) => (prev - 1 + QR_STYLES.length) % QR_STYLES.length);

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
        downloadLink.click();
      }
      URL.revokeObjectURL(url);
    };
    img.src = url;
  };

  return (
    <ResponsiveSheet isOpen={isOpen} onClose={onClose} title="QR Code">
      <div className="flex flex-col items-center gap-6 py-2">
        <div className="flex items-center gap-4 w-full justify-between">
          <Button
            variant="ghost"
            size="icon"
            onClick={prevStyle}
            className="rounded-full h-10 w-10 shrink-0 border border-slate-100 dark:border-slate-800"
          >
            <ChevronLeft size={20} />
          </Button>

          <div
            ref={svgRef}
            className={`p-6 rounded-2xl shadow-xl transition-all duration-300 border ${currentStyle.wrapperClass}`}
          >
            <QRCodeSVG
              value={linkUrl}
              size={180}
              bgColor={currentStyle.bg}
              fgColor={currentStyle.fg}
              level="M"
            />
          </div>

          <Button
            variant="ghost"
            size="icon"
            onClick={nextStyle}
            className="rounded-full h-10 w-10 shrink-0 border border-slate-100 dark:border-slate-800"
          >
            <ChevronRight size={20} />
          </Button>
        </div>

        <div className="text-center space-y-1">
          <p className="text-sm font-bold text-slate-900 dark:text-slate-100 flex items-center justify-center gap-2">
            <span
              className="w-3 h-3 rounded-full border border-slate-200"
              style={{ backgroundColor: currentStyle.fg }}
            />
            {currentStyle.name}
          </p>
          <p className="text-xs text-slate-500 truncate max-w-[250px]">
            {linkUrl}
          </p>
        </div>

        <div className="w-full pt-4">
          <Button
            onClick={handleDownload}
            className="w-full gap-2 h-12 text-base bg-indigo-600 hover:bg-indigo-700 text-white shadow-lg shadow-indigo-200 dark:shadow-none"
          >
            <Download size={18} /> Download PNG
          </Button>
        </div>
      </div>
    </ResponsiveSheet>
  );
};
