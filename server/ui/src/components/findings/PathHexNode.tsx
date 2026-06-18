import { memo } from "react";
import { getHexConfig } from "@/lib/explorer/hex-config";
import { getPropertyChips } from "@/lib/findings/property-chips";
import { cn } from "@/lib/utils";
import { SEVERITY, severityColor } from "@/theme/tokens";
import type { AttackPathNode } from "@/api/types";

interface PathHexNodeProps {
  node: AttackPathNode;
  isFirst: boolean;
  isLast: boolean;
  severity?: string;
}

const SCALE = 64 / 84;
const HEX_W = 64;
const HEX_H = Math.round(96 * SCALE);
const SCALED_POINTS = [
  [42, 4], [78, 22], [78, 74], [42, 92], [6, 74], [6, 22],
]
  .map(([x, y]) => `${Math.round(x! * SCALE)},${Math.round(y! * SCALE)}`)
  .join(" ");

function PathHexNodeComponent({ node, isFirst, isLast, severity }: PathHexNodeProps) {
  const kind = node.kinds[0] ?? "";
  const config = getHexConfig(kind);
  const Icon = config.icon;
  const name = (node.properties?.name as string) || node.id.slice(0, 12);
  const chips = getPropertyChips(kind, node.properties ?? {});

  const sevColor = severity ? severityColor(severity) : undefined;
  const accent = isLast && sevColor ? sevColor : config.strokeColor;
  const categoryLabel = isFirst ? `ENTRY \u00b7 ${config.kindTag}` : config.kindTag;
  const critTarget = isLast && severity === "critical";

  return (
    <div
      className={cn(
        "relative flex w-[150px] flex-shrink-0 flex-col items-center rounded-[3px] border border-border/70 bg-black/30 px-3 py-3",
      )}
      style={{
        boxShadow: critTarget
          ? `inset 2px 0 0 0 ${accent}, 0 0 18px -6px ${sevColor}`
          : `inset 2px 0 0 0 ${accent}`,
      }}
    >
      <span aria-hidden className="pointer-events-none absolute inset-x-0 top-0 h-px bg-white/[0.05]" />

      <div
        className="mb-2 w-full truncate text-center font-mono text-[9px] font-semibold uppercase tracking-[0.14em]"
        style={{ color: accent }}
      >
        {categoryLabel}
      </div>

      <div className="relative mb-2" style={{ width: HEX_W, height: HEX_H }}>
        <svg
          width={HEX_W}
          height={HEX_H}
          viewBox={`0 0 ${HEX_W} ${HEX_H}`}
          className="absolute inset-0"
        >
          <polygon
            points={SCALED_POINTS}
            fill={config.fillColor}
            stroke={config.strokeColor}
            strokeWidth={2}
            strokeLinejoin="round"
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <Icon className="h-5 w-5" style={{ color: config.strokeColor }} strokeWidth={2} />
        </div>
      </div>

      <div className="w-full truncate text-center text-[11px] font-semibold text-foreground" title={name}>
        {name}
      </div>

      {chips.length > 0 && (
        <div className="mt-1.5 flex flex-wrap justify-center gap-1">
          {chips.map((chip) => {
            const tone =
              chip === "critical" || chip === "exposed"
                ? SEVERITY.critical
                : chip === "high"
                  ? SEVERITY.high
                  : null;
            return (
              <span
                key={chip}
                className="rounded-[2px] px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-[0.04em]"
                style={
                  tone
                    ? { backgroundColor: tone.bg, color: tone.text, boxShadow: `inset 0 0 0 1px ${tone.border}` }
                    : { backgroundColor: "rgb(255 255 255 / 0.05)", color: "rgb(var(--mauve-11-raw))" }
                }
              >
                {chip}
              </span>
            );
          })}
        </div>
      )}
    </div>
  );
}

export const PathHexNode = memo(PathHexNodeComponent);
