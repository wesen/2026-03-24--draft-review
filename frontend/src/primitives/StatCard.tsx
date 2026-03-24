import "./StatCard.css";

interface StatCardProps {
  icon: string;
  value: string | number;
  label: string;
}

export function StatCard({ icon, value, label }: StatCardProps) {
  return (
    <div className="dr-stat-card">
      <div className="dr-stat-card__value">
        {icon} {value}
      </div>
      <div className="dr-stat-card__label">{label}</div>
    </div>
  );
}
