import { useParams } from 'react-router-dom';

export const AnalyticsPage = () => {
  const { id } = useParams();
  return (
    <div className="p-4 border-2 border-dashed border-slate-300 rounded-xl h-96 flex items-center justify-center bg-slate-50">
      <div className="text-center">
        <h2 className="text-xl font-semibold text-slate-700">
          Аналитика ссылки: {id}
        </h2>
        <p className="text-slate-500">
          Будет реализовано в Части 3 (Графики: Streamgraph, Radar, Sankey)
        </p>
      </div>
    </div>
  );
};
