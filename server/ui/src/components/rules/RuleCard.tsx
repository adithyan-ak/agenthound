import { ChevronDown, ChevronUp } from "lucide-react";
import type { RuleInfo } from "@/api/rules";
import { severityColor } from "@/theme/tokens";
import { RuleDetail } from "./RuleDetail";

const COLLECTOR_COLOR: Record<string, string> = {
  mcp: "#10B981",
  a2a: "#A855F7",
  config: "#D97706",
  all: "#7A828E",
};

interface RuleCardProps {
  rule: RuleInfo;
  isExpanded: boolean;
  onToggle: () => void;
}

export function RuleCard({ rule, isExpanded, onToggle }: RuleCardProps) {
  const sevColor = severityColor(rule.severity);
  const collectorColor = COLLECTOR_COLOR[rule.collector] ?? "#7A828E";
  const isCustom = rule.source === "custom";

  return (
    <div className="card-elevated relative overflow-hidden rounded-md">
      <span
        aria-hidden
        className="pointer-events-none absolute inset-y-0 left-0 w-0.5"
        style={{ backgroundColor: sevColor, opacity: 0.85 }}
      />
      <button
        onClick={onToggle}
        className="flex w-full items-center gap-3 px-3.5 py-3 text-left transition-colors hover:bg-white/[0.02]"
      >
        <div className="min-w-0 flex-1">
          <div className="mb-0.5 flex flex-wrap items-center gap-2">
            <span className="text-sm font-medium text-foreground">{rule.name}</span>
            <span
              className="inline-flex items-center gap-1 rounded-[2px] px-1.5 py-0.5 font-mono text-[10px] font-semibold uppercase tracking-[0.08em]"
              style={{ backgroundColor: `${sevColor}1A`, color: sevColor, boxShadow: `inset 0 0 0 1px ${sevColor}55` }}
            >
              <span className="h-1.5 w-1.5 rounded-[1px]" style={{ backgroundColor: sevColor }} />
              {rule.severity}
            </span>
            <span className="inline-flex items-center gap-1.5 rounded-[2px] border border-border bg-black/40 px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-[0.06em] text-muted-foreground">
              <span className="h-1.5 w-1.5 rounded-[1px]" style={{ backgroundColor: collectorColor }} />
              {rule.collector}
            </span>
            <span
              className={`rounded-[2px] border px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-[0.06em] ${
                isCustom
                  ? "border-primary/40 bg-primary/10 text-primary/90"
                  : "border-border bg-black/40 text-muted-foreground"
              }`}
            >
              {rule.source}
            </span>
          </div>
          <p className="truncate text-xs text-muted-foreground">{rule.description}</p>
          <div className="mt-1 flex flex-wrap gap-1">
            {rule.targets.map((t) => (
              <span
                key={t}
                className="rounded-[2px] border border-border bg-black/40 px-1 py-0 font-mono text-[9px] uppercase tracking-[0.04em] text-muted-foreground"
              >
                {t}
              </span>
            ))}
            {rule.owasp?.map((tag) => (
              <span
                key={tag}
                className="rounded-[2px] border border-border bg-black/40 px-1 py-0 font-mono text-[9px] uppercase tracking-[0.06em] text-primary/80"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
        <div className="flex flex-shrink-0 items-center gap-2.5">
          {rule.test_count > 0 && (
            <span className="font-mono text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
              {rule.test_count} tests
            </span>
          )}
          {isExpanded ? (
            <ChevronUp className="h-4 w-4 text-primary" />
          ) : (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          )}
        </div>
      </button>

      {isExpanded && (
        <>
          <div className="h-px bg-border/70" />
          <div className="bg-black/20 px-3.5 py-3">
            <RuleDetail ruleId={rule.id} />
          </div>
        </>
      )}
    </div>
  );
}
