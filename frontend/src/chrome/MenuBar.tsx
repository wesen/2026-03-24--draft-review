import { useState, type ReactNode } from "react";
import "./MenuBar.css";

export interface MenuItem {
  label: string;
  action?: () => void;
  shortcut?: string;
  divider?: boolean;
}

interface MenuBarProps {
  menus: { label: string; items: MenuItem[] }[];
  rightStatus?: ReactNode;
}

function MenuDropdown({ items, onClose }: { items: MenuItem[]; onClose: () => void }) {
  return (
    <div className="dr-menu-dropdown">
      {items.map((item, i) =>
        item.divider ? (
          <div key={i} className="dr-menu-dropdown__divider" />
        ) : (
          <div
            key={i}
            className="dr-menu-dropdown__item"
            onClick={() => {
              item.action?.();
              onClose();
            }}
          >
            <span>{item.label}</span>
            {item.shortcut && (
              <span className="dr-menu-dropdown__shortcut">{item.shortcut}</span>
            )}
          </div>
        )
      )}
    </div>
  );
}

export function MenuBar({ menus, rightStatus }: MenuBarProps) {
  const [openMenu, setOpenMenu] = useState<string | null>(null);

  return (
    <div className="dr-menubar">
      <div className="dr-menubar__logo">{"\uF8FF"}</div>
      {menus.map((menu) => (
        <div
          key={menu.label}
          className="dr-menubar__item-wrapper"
          onMouseEnter={() => openMenu && setOpenMenu(menu.label)}
          onClick={() =>
            setOpenMenu(openMenu === menu.label ? null : menu.label)
          }
        >
          <div
            className={`dr-menubar__item ${
              openMenu === menu.label ? "dr-menubar__item--active" : ""
            }`}
          >
            {menu.label}
          </div>
          {openMenu === menu.label && (
            <MenuDropdown
              items={menu.items}
              onClose={() => setOpenMenu(null)}
            />
          )}
        </div>
      ))}
      <div className="dr-menubar__spacer" />
      {rightStatus && (
        <div className="dr-menubar__status">{rightStatus}</div>
      )}
      {openMenu && (
        <div
          className="dr-menubar__backdrop"
          onClick={() => setOpenMenu(null)}
        />
      )}
    </div>
  );
}
