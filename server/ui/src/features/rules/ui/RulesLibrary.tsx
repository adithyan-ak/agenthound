import { useState, useMemo } from "react";
import { ShieldCheck } from "lucide-react";
import { useRules } from "@entities/rule";
import type { RuleInfo } from "@entities/rule";
import { Skeleton } from "@shared/ui/primitives/skeleton";
import { RuleCard } from "./RuleCard";
import { RuleFilters } from "./RuleFilters";

const TAG_ORDER = [
  "injection",
  "credential",
  "supply-chain",
  "sensitivity",
  "capability",
  "impersonation",
  "instruction-poisoning",
];

function groupByTag(rules: RuleInfo[]): Map<string, RuleInfo[]> {
  const grouped = new Map<string, RuleInfo[]>();
  for (const rule of rules) {
    const tag = rule.tags[0] ?? "other";
    const list = grouped.get(tag) ?? [];
    list.push(rule);
    grouped.set(tag, list);
  }
  return grouped;
}

export function RulesLibrary() {
  const [severity, setSeverity] = useState("all");
  const [collector, setCollector] = useState("all");
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const filterParams = useMemo(() => {
    const params: { severity?: string; collector?: string } = {};
    if (severity !== "all") params.severity = severity;
    if (collector !== "all") params.collector = collector;
    return params;
  }, [severity, collector]);

  const { data, isLoading } = useRules(filterParams);

  const grouped = useMemo(() => groupByTag(data?.rules ?? []), [data?.rules]);

  const sortedTags = useMemo(() => {
    const tags = TAG_ORDER.filter((t) => grouped.has(t));
    for (const t of grouped.keys()) {
      if (!tags.includes(t)) tags.push(t);
    }
    return tags;
  }, [grouped]);

  function handleClear() {
    setSeverity("all");
    setCollector("all");
  }

  return (
    <div className="dashboard-bg min-h-full p-3 sm:p-4 lg:p-5">
      <div className="mx-auto max-w-[1100px] space-y-4">
        <header className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          <div className="min-w-0">
            <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
              Detection Engine <span className="text-primary/60">//</span> Ruleset
            </p>
            <h1 className="mt-1.5 flex items-center gap-2.5 font-mono text-2xl font-bold uppercase tracking-[0.04em] text-foreground sm:text-[26px]">
              <span className="flex h-7 w-7 items-center justify-center rounded-[3px] bg-primary/10 ring-1 ring-inset ring-primary/30">
                <ShieldCheck className="h-4 w-4 text-primary" />
              </span>
              <span className="text-primary">▸</span>
              Detection Rules
              {data && (
                <span className="font-mono text-base font-semibold tabular-nums text-muted-foreground">
                  {String(data.total).padStart(2, "0")}
                </span>
              )}
              <span className="blink-caret text-primary" aria-hidden>
                _
              </span>
            </h1>
            <p className="mt-1.5 text-sm text-muted-foreground">
              The matchers AgentHound runs against collected metadata to raise findings.
            </p>
          </div>
          <RuleFilters
            severity={severity}
            collector={collector}
            onSeverityChange={setSeverity}
            onCollectorChange={setCollector}
            onClear={handleClear}
          />
        </header>

        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="space-y-2">
                <Skeleton className="h-3 w-32 rounded-[2px]" />
                <Skeleton className="h-14 w-full rounded-[3px]" />
                <Skeleton className="h-14 w-full rounded-[3px]" />
              </div>
            ))}
          </div>
        ) : data?.rules.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 py-16 text-center">
            <span className="flex h-12 w-12 items-center justify-center rounded-[4px] bg-primary/10 ring-1 ring-inset ring-primary/20">
              <ShieldCheck className="h-6 w-6 text-primary" />
            </span>
            <p className="font-mono text-xs uppercase tracking-[0.12em] text-muted-foreground">
              No rules match the selected filters
            </p>
          </div>
        ) : (
          <div className="space-y-6">
            {sortedTags.map((tag) => {
              const rules = grouped.get(tag) ?? [];
              return (
                <div key={tag}>
                  <div className="mb-2.5 flex items-center gap-2">
                    <span aria-hidden className="h-px w-6 bg-primary/50" />
                    <h3 className="font-mono text-console uppercase tracking-[0.18em] text-muted-foreground">
                      {tag}
                    </h3>
                    <span aria-hidden className="h-px flex-1 bg-border/60" />
                    <span className="font-mono text-[10px] tabular-nums text-muted-foreground/70">
                      {String(rules.length).padStart(2, "0")}
                    </span>
                  </div>
                  <div className="space-y-2">
                    {rules.map((rule) => (
                      <RuleCard
                        key={rule.id}
                        rule={rule}
                        isExpanded={expandedId === rule.id}
                        onToggle={() => setExpandedId(expandedId === rule.id ? null : rule.id)}
                      />
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
