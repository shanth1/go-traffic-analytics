import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react';
import { Button } from '@/shared/ui/button';
import { cn } from '@/shared/lib/utils';

interface PaginationProps {
  total: number;
  limit: number;
  offset: number;
  onChange: (newOffset: number) => void;
  className?: string;
}

export const Pagination = ({
  total,
  limit,
  offset,
  onChange,
  className,
}: PaginationProps) => {
  const currentPage = Math.floor(offset / limit) + 1;
  const totalPages = Math.ceil(total / limit);

  if (totalPages <= 1) return null;

  const handlePageChange = (page: number) => {
    const newOffset = (page - 1) * limit;
    onChange(newOffset);
  };

  const renderPageNumbers = () => {
    const pages = [];
    const maxVisiblePages = 5;

    if (totalPages <= maxVisiblePages) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always show first, last, current, and neighbors
      if (currentPage <= 3) {
        pages.push(1, 2, 3, 4, '...', totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(
          1,
          '...',
          totalPages - 3,
          totalPages - 2,
          totalPages - 1,
          totalPages
        );
      } else {
        pages.push(
          1,
          '...',
          currentPage - 1,
          currentPage,
          currentPage + 1,
          '...',
          totalPages
        );
      }
    }

    return pages.map((page, idx) => {
      if (page === '...') {
        return (
          <div
            key={`ellipsis-${idx}`}
            className="flex items-center justify-center w-9 h-9"
          >
            <MoreHorizontal className="h-4 w-4 text-slate-400" />
          </div>
        );
      }

      const p = page as number;
      const isActive = p === currentPage;

      return (
        <Button
          key={p}
          variant={isActive ? 'default' : 'outline'}
          size="icon"
          className={cn(
            'w-9 h-9 transition-all',
            isActive
              ? 'pointer-events-none'
              : 'hover:bg-slate-100 dark:hover:bg-slate-800'
          )}
          onClick={() => handlePageChange(p)}
        >
          {p}
        </Button>
      );
    });
  };

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Button
        variant="outline"
        size="icon"
        className="w-9 h-9"
        onClick={() => handlePageChange(currentPage - 1)}
        disabled={currentPage === 1}
      >
        <ChevronLeft className="h-4 w-4" />
      </Button>

      <div className="flex items-center gap-1">{renderPageNumbers()}</div>

      <Button
        variant="outline"
        size="icon"
        className="w-9 h-9"
        onClick={() => handlePageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
      >
        <ChevronRight className="h-4 w-4" />
      </Button>

      <div className="ml-4 text-xs text-slate-500 hidden sm:block">
        Showing {offset + 1}-{Math.min(offset + limit, total)} of {total}
      </div>
    </div>
  );
};
