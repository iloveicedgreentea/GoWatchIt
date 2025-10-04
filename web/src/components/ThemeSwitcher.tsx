import { Palette } from 'lucide-react';
import { useTheme } from '../theme/provider';
import { AVAILABLE_THEMES } from '../theme/config';

export function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();

  return (
    <div className="dropdown dropdown-end">
      <button
        tabIndex={0}
        className="btn btn-ghost btn-circle hover:bg-base-300 tooltip tooltip-left"
        aria-label="Change theme"
        data-tip="Change theme"
      >
        <Palette className="h-5 w-5" />
      </button>
      <ul
        tabIndex={0}
        className="dropdown-content menu bg-base-200 border border-base-300 rounded-box z-[1] w-64 p-2 shadow-2xl mt-3"
      >
        <li className="menu-title">
          <span className="text-sm font-semibold">Select Theme</span>
        </li>
        {AVAILABLE_THEMES.map((themeOption) => (
          <li key={themeOption.name}>
            <button
              onClick={() => setTheme(themeOption.name)}
              className={`${theme === themeOption.name ? 'active bg-primary/20' : 'hover:bg-base-300'}`}
            >
              <div className="flex flex-col items-start">
                <span className="font-medium">{themeOption.displayName}</span>
                <span className="text-xs opacity-60">{themeOption.description}</span>
              </div>
              {theme === themeOption.name && (
                <span className="ml-auto text-primary">✓</span>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
