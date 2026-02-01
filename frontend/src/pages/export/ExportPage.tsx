import { useState } from 'react';
import { FileDownIcon, LockIcon, SparklesIcon } from 'lucide-react';
import { useAuthStore } from '@/entities/session/store';
import { api } from '@/shared/api/base';
import { Button } from '@/shared/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/ui/card';
import { DateRangePicker } from '@/features/analytics-filters/DateRangePicker';
import { useAnalyticsFilter } from '@/entities/analytics/model/filters';
import { toRFC3339 } from '@/shared/lib/date';
import { TG_SUPPORT_URL } from '@/shared/config';
import { toast } from '@/entities/notification/store';

export const ExportPage = () => {
  const { user } = useAuthStore();
  const { startDate, endDate } = useAnalyticsFilter();
  const [isExporting, setIsExporting] = useState(false);

  const canExport =
    user?.role === 'admin' ||
    user?.plan_id === 'pro' ||
    user?.plan_id === 'enterprise';

  const handleExport = async () => {
    setIsExporting(true);
    try {
      const params = {
        from: toRFC3339(startDate),
        to: toRFC3339(endDate),
      };

      const response = await api.get('/analytics/export', {
        params,
        responseType: 'blob',
      });

      const blob = new Blob([response.data], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;

      const fileName = `export_${params.from.split('T')[0]}_to_${params.to.split('T')[0]}.xlsx`;
      link.setAttribute('download', fileName);

      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      toast.info('Success', 'Your data export is ready.');
    } catch (error) {
      console.error('Export failed', error);
    } finally {
      setIsExporting(false);
    }
  };

  if (!canExport) {
    return (
      <div className="flex items-center justify-center min-h-[60vh] animate-in fade-in zoom-in-95 duration-300">
        <Card className="max-w-md text-center p-8 border-dashed border-2 bg-card/50 backdrop-blur-sm">
          <div className="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-6">
            <LockIcon className="text-primary" size={40} />
          </div>
          <CardTitle className="text-2xl mb-3">Feature Locked</CardTitle>
          <CardDescription className="text-base mb-8">
            Raw data export to Excel is available only for{' '}
            <strong>Professional</strong> plans. Upgrade to access deep
            analytics and reporting.
          </CardDescription>
          <Button
            className="w-full h-12 gap-2 font-bold"
            onClick={() => window.open(TG_SUPPORT_URL, '_blank')}
          >
            <SparklesIcon size={18} /> Upgrade Plan
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500 min-h-[calc(100vh-160px)] flex flex-col">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-border pb-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">
            Data Export
          </h1>
          <p className="text-muted-foreground mt-1">
            Generate and download reports based on your filters.
          </p>
        </div>
        <DateRangePicker />
      </div>

      <div className="flex-1 flex items-center justify-center">
        <Card className="w-full max-w-md bg-card/30 backdrop-blur-md border-border shadow-2xl overflow-hidden">
          <CardHeader className="text-center pb-2">
            <div className="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-4 text-primary">
              <FileDownIcon size={32} />
            </div>
            <CardTitle>Excel Spreadsheet</CardTitle>
            <CardDescription>
              Exporting data for the period: <br />
              <span className="font-bold text-foreground">
                {startDate.toLocaleDateString()} —{' '}
                {endDate.toLocaleDateString()}
              </span>
            </CardDescription>
          </CardHeader>

          <CardContent className="p-6 pt-4 space-y-4">
            <Button
              size="lg"
              className="w-full h-16 text-lg font-bold gap-3 shadow-xl shadow-primary/20 hover:scale-[1.01] transition-transform"
              onClick={handleExport}
              disabled={isExporting}
            >
              {isExporting ? (
                'Generating File...'
              ) : (
                <>
                  <FileDownIcon size={24} />
                  Download XLSX
                </>
              )}
            </Button>
            <p className="text-[10px] text-center uppercase tracking-widest text-muted-foreground font-semibold">
              Raw Data
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
