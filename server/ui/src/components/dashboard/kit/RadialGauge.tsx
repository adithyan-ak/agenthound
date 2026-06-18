import { useEffect, useId, useState } from "react";
import { cn } from "@/lib/utils";

interface RadialGaugeProps {
  /** 0-100. */
  value: number;
  /** Color for the readout number + needle. Defaults to the gradient end. */
  valueColor?: string;
  /** Large readout shown inside the arc. Defaults to the rounded value. */
  readout?: string;
  /** Small caption under the readout. */
  caption?: string;
  size?: number;
  className?: string;
}

const TRACK = "rgba(255,255,255,0.07)";

/**
 * Semicircular speedometer gauge drawn in pure SVG (no chart dependency).
 * The arc reveals a green->amber->red risk scale up to `value`, and a needle
 * points to the same position. Animates on mount via a CSS transition.
 */
export function RadialGauge({
  value,
  valueColor,
  readout,
  caption,
  size = 208,
  className,
}: RadialGaugeProps) {
  const gradientId = useId();
  const glowId = useId();
  const clamped = Math.max(0, Math.min(100, value));

  const [animated, setAnimated] = useState(0);
  useEffect(() => {
    const frame = requestAnimationFrame(() => setAnimated(clamped));
    return () => cancelAnimationFrame(frame);
  }, [clamped]);

  const strokeWidth = 16;
  const pad = strokeWidth / 2 + 6;
  const cx = size / 2;
  const cy = size / 2;
  const r = size / 2 - pad;
  const height = cy + 34;

  const arcPath = `M ${cx - r} ${cy} A ${r} ${r} 0 0 1 ${cx + r} ${cy}`;

  // Needle: value 0 -> 180deg (left), value 100 -> 0deg (right), over the top.
  const needleR = r - strokeWidth / 2 - 2;
  const angle = (180 - 1.8 * animated) * (Math.PI / 180);
  const tipX = cx + needleR * Math.cos(angle);
  const tipY = cy - needleR * Math.sin(angle);

  const display = readout ?? String(Math.round(clamped));

  return (
    <div className={cn("relative", className)}>
      <svg
        width="100%"
        viewBox={`0 0 ${size} ${height}`}
        className="overflow-visible"
        role="img"
        aria-label={caption ? `${caption}: ${display}` : display}
      >
        <defs>
          <linearGradient id={gradientId} x1="0" y1="0" x2="1" y2="0">
            <stop offset="0%" stopColor="#22C55E" />
            <stop offset="45%" stopColor="#EAB308" />
            <stop offset="75%" stopColor="#F97316" />
            <stop offset="100%" stopColor="#EF4444" />
          </linearGradient>
          <filter id={glowId} x="-30%" y="-30%" width="160%" height="160%">
            <feGaussianBlur stdDeviation="4" result="blur" />
            <feMerge>
              <feMergeNode in="blur" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>

        <path
          d={arcPath}
          fill="none"
          stroke={TRACK}
          strokeWidth={strokeWidth}
          strokeLinecap="round"
        />
        <path
          d={arcPath}
          fill="none"
          stroke={`url(#${gradientId})`}
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          pathLength={100}
          strokeDasharray={`${animated} 100`}
          filter={`url(#${glowId})`}
          style={{ transition: "stroke-dasharray 900ms cubic-bezier(0.22,1,0.36,1)" }}
        />

        <line
          x1={cx}
          y1={cy}
          x2={tipX}
          y2={tipY}
          stroke={valueColor ?? "#EDF0F3"}
          strokeWidth={3}
          strokeLinecap="round"
          style={{ transition: "all 900ms cubic-bezier(0.22,1,0.36,1)" }}
        />
        <circle cx={cx} cy={cy} r={6} fill={valueColor ?? "#EDF0F3"} />
        <circle cx={cx} cy={cy} r={11} fill="none" stroke={TRACK} strokeWidth={2} />

        <text
          x={cx}
          y={cy - 14}
          textAnchor="middle"
          className="font-mono font-bold"
          style={{ fontSize: 38, fill: valueColor ?? "#EDF0F3" }}
        >
          {display}
        </text>
        {caption && (
          <text
            x={cx}
            y={cy + 22}
            textAnchor="middle"
            style={{
              fontSize: 11,
              letterSpacing: "0.12em",
              textTransform: "uppercase",
              fill: "#64788F",
            }}
          >
            {caption}
          </text>
        )}
      </svg>
    </div>
  );
}
