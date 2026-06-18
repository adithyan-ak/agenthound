import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { fetchRuleDetail } from "@/api/rules";
import type { MatcherSpec, TestCase } from "@/api/rules";

interface RuleDetailProps {
  ruleId: string;
}

const SIGNAL_OK = "#3FB950";

function Label({ children }: { children: ReactNode }) {
  return (
    <div className="font-mono text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
      {children}
    </div>
  );
}

function Chips({ items }: { items?: string[] }) {
  if (!items || items.length === 0) return null;
  return (
    <div className="flex flex-wrap gap-1.5">
      {items.map((kw) => (
        <span
          key={kw}
          className="rounded-[2px] border border-border bg-black/40 px-1.5 py-0.5 font-mono text-[11px] text-foreground/80"
        >
          {kw}
        </span>
      ))}
    </div>
  );
}

function MatcherDisplay({ matcher }: { matcher: MatcherSpec }) {
  switch (matcher.Type) {
    case "regex":
      return (
        <div className="space-y-2">
          <Label>Pattern {matcher.CaseInsensitive && "(case-insensitive)"}</Label>
          <code className="block break-all rounded-[3px] border border-border bg-black/50 px-3 py-2 font-mono text-xs text-foreground/90">
            {matcher.Pattern}
          </code>
        </div>
      );

    case "keyword":
      return (
        <div className="space-y-2">
          <Label>
            Keywords
            {matcher.CaseInsensitive && " (case-insensitive)"}
            {matcher.MatchMode && ` \u00b7 match ${matcher.MatchMode}`}
          </Label>
          <Chips items={matcher.Keywords} />
        </div>
      );

    case "prefix":
      return (
        <div className="space-y-2">
          <Label>Prefixes {matcher.CaseInsensitive && "(case-insensitive)"}</Label>
          <Chips items={matcher.Prefixes} />
        </div>
      );

    case "entropy":
      return (
        <div className="space-y-2">
          <Label>Entropy Check</Label>
          <div className="flex flex-wrap gap-x-4 gap-y-1 font-mono text-xs">
            <span>
              <span className="text-muted-foreground">charset</span>{" "}
              <span className="text-foreground">{matcher.Charset}</span>
            </span>
            <span>
              <span className="text-muted-foreground">threshold</span>{" "}
              <span className="text-foreground">{matcher.Threshold}</span>
            </span>
            {matcher.MinLength != null && matcher.MinLength > 0 && (
              <span>
                <span className="text-muted-foreground">min length</span>{" "}
                <span className="text-foreground">{matcher.MinLength}</span>
              </span>
            )}
          </div>
        </div>
      );

    case "compound":
      return (
        <div className="space-y-3">
          <Label>Compound ({matcher.Operator?.toUpperCase()})</Label>
          <div className="space-y-2 border-l border-border pl-3">
            {matcher.Matchers?.map((sub, i) => (
              <MatcherDisplay key={i} matcher={sub} />
            ))}
          </div>
        </div>
      );

    default:
      return (
        <div className="font-mono text-xs text-muted-foreground">
          Unknown matcher type: {matcher.Type}
        </div>
      );
  }
}

function TestCasesTable({ tests }: { tests: TestCase[] }) {
  if (tests.length === 0) return null;

  return (
    <div className="space-y-2">
      <Label>Test Cases ({tests.length})</Label>
      <div className="overflow-x-auto rounded-[3px] border border-border/70">
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-border bg-black/30">
              <th className="px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.1em] text-muted-foreground">
                Input
              </th>
              <th className="w-[90px] px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.1em] text-muted-foreground">
                Expected
              </th>
              <th className="px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.1em] text-muted-foreground">
                Description
              </th>
            </tr>
          </thead>
          <tbody>
            {tests.map((tc, i) => (
              <tr key={i} className="border-b border-border/50 last:border-0">
                <td className="px-3 py-1.5 align-top">
                  <code className="break-all font-mono text-[11px] text-foreground/90">{tc.Input}</code>
                </td>
                <td className="px-3 py-1.5 align-top">
                  {tc.ShouldMatch ? (
                    <span
                      className="inline-flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-[0.06em]"
                      style={{ color: SIGNAL_OK }}
                    >
                      <span className="h-1.5 w-1.5 rounded-[1px]" style={{ backgroundColor: SIGNAL_OK }} />
                      match
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-[0.06em] text-muted-foreground">
                      <span className="h-1.5 w-1.5 rounded-[1px] bg-mauve-8" />
                      no match
                    </span>
                  )}
                </td>
                <td className="px-3 py-1.5 align-top text-xs text-muted-foreground">{tc.Description}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function RuleDetail({ ruleId }: RuleDetailProps) {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["rule-detail", ruleId],
    queryFn: () => fetchRuleDetail(ruleId),
    staleTime: 30_000,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center gap-2 py-4 font-mono text-xs uppercase tracking-[0.1em] text-muted-foreground">
        <Loader2 className="h-4 w-4 animate-spin text-primary" />
        Loading rule details…
      </div>
    );
  }

  if (isError) {
    return (
      <div
        className="rounded-[3px] border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
        style={{ boxShadow: "inset 2px 0 0 0 rgb(var(--tomato-9-raw))" }}
      >
        {error instanceof Error ? error.message : "Failed to load rule"}
      </div>
    );
  }

  if (!data) return null;

  return (
    <div className="space-y-3.5">
      <p className="text-sm leading-relaxed text-muted-foreground">{data.description}</p>
      <div className="h-px bg-border/60" />
      <MatcherDisplay matcher={data.matcher} />
      {data.tests.length > 0 && (
        <>
          <div className="h-px bg-border/60" />
          <TestCasesTable tests={data.tests} />
        </>
      )}
    </div>
  );
}
