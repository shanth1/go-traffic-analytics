import { NavLink } from 'react-router-dom';
import { Home, BarChart2, User, PlusCircle } from 'lucide-react';
import { cn } from '@/shared/lib/utils';

const items = [
  { to: '/campaigns', icon: Home, label: 'Campaigns' },
  { to: '/create', icon: PlusCircle, label: 'Create', isSpecial: true },
  { to: '/analytics', icon: BarChart2, label: 'Analytics' },
  { to: '/profile', icon: User, label: 'Profile' },
];

export const BottomMenu = () => {
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-surface/90 backdrop-blur-md border-t border-slate-800 pb-safe pt-2 px-6 h-20 z-50">
      <ul className="flex justify-between items-center max-w-md mx-auto">
        {items.map((item) => (
          <li key={item.label}>
            <NavLink
              to={item.to}
              className={({ isActive }) =>
                cn(
                  'flex flex-col items-center gap-1 transition-colors',
                  isActive
                    ? 'text-primary'
                    : 'text-slate-500 hover:text-slate-300',
                  item.isSpecial &&
                    'text-white bg-primary p-3 rounded-full -mt-6 shadow-lg shadow-blue-500/20 border-4 border-background'
                )
              }
            >
              <item.icon size={item.isSpecial ? 24 : 20} />
              {!item.isSpecial && (
                <span className="text-[10px]">{item.label}</span>
              )}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  );
};
