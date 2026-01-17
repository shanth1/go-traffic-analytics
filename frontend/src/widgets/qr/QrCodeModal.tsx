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

// --- Styles Definition ---
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
    bg: '#0f172a', // slate-950
    fg: '#ffffff',
    wrapperClass: 'bg-slate-950',
  },
  {
    id: 'brand',
    name: 'Brand Indigo',
    bg: '#e0e7ff', // indigo-100
    fg: '#4f46e5', // indigo-600
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

  const nextStyle = () => {
    setStyleIndex((prev) => (prev + 1) % QR_STYLES.length);
  };

  const prevStyle = () => {
    setStyleIndex((prev) => (prev - 1 + QR_STYLES.length) % QR_STYLES.length);
  };

  const handleDownload = () => {
    const svg = svgRef.current?.querySelector('svg');
    if (!svg) return;

    // Сериализация SVG в строку
    const svgData = new XMLSerializer().serializeToString(svg);
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const img = new Image();

    // Размер итоговой картинки (1024x1024 для хорошего качества)
    const size = 1024;
    canvas.width = size;
    canvas.height = size;

    // Создаем Blob URL
    const svgBlob = new Blob([svgData], {
      type: 'image/svg+xml;charset=utf-8',
    });
    const url = URL.createObjectURL(svgBlob);

    img.onload = () => {
      if (ctx) {
        // Рисуем фон (обязательно для PNG, чтобы не был прозрачным)
        ctx.fillStyle = currentStyle.bg;
        ctx.fillRect(0, 0, size, size);

        // Рисуем QR код
        ctx.drawImage(img, 0, 0, size, size);

        // Скачивание
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

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      <div className="relative bg-white dark:bg-slate-900 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden flex flex-col items-center animate-in fade-in zoom-in-95 duration-200">
        {/* Заголовок */}
        <div className="w-full flex justify-between items-center p-4 border-b border-slate-100 dark:border-slate-800">
          <h3 className="font-semibold text-lg flex items-center gap-2 text-slate-900 dark:text-slate-100">
            <QrCodeIcon size={18} className="text-indigo-600" />
            QR Code
          </h3>
          <button
            onClick={onClose}
            className="p-1 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full text-slate-500"
          >
            <X size={20} />
          </button>
        </div>

        {/* Область предпросмотра */}
        <div className="p-8 w-full flex flex-col items-center gap-6">
          <div className="flex items-center gap-4 w-full justify-between">
            <Button
              variant="ghost"
              size="icon"
              onClick={prevStyle}
              className="rounded-full h-10 w-10 shrink-0"
            >
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
                includeMargin={false}
              />
            </div>

            <Button
              variant="ghost"
              size="icon"
              onClick={nextStyle}
              className="rounded-full h-10 w-10 shrink-0"
            >
              <ChevronRight />
            </Button>
          </div>

          <div className="text-center">
            <p className="text-sm font-medium text-slate-900 dark:text-slate-100">
              {currentStyle.name}
            </p>
            <p className="text-xs text-slate-500 mt-1 truncate max-w-[200px] mx-auto">
              {linkUrl}
            </p>
          </div>
        </div>

        {/* Футер */}
        <div className="w-full p-4 border-t border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-950/50">
          <Button
            onClick={handleDownload}
            className="w-full gap-2 bg-indigo-600 hover:bg-indigo-700 text-white shadow-md shadow-indigo-200 dark:shadow-none"
          >
            <Download size={16} /> Download PNG
          </Button>
        </div>
      </div>
    </div>
  );
};
