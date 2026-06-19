import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowRight } from "lucide-react";
import { useFindings, SEVERITY_RANK } from "@entities/finding";
import { edgeLabel } from "@entities/edge";
import { SEVERITY, SEVERITY_BY_KEY } from "@shared/theme/tokens";

/**
 * Lists the findings that touch this node (as source or target) with a
 * deep-link into the dossier — the node-level half of the map → findings
 * handoff. Uses the shared findings cache, so it is free once the page has
 * loaded.
 */
export function FindingsTab({ nodeId }: { nodeId: string }) {
  const navigate = useNavigate();
  const { data: findings, isLoading } = useFindings();

  const related = useMemo(() => {
    return (findings ?? [])
      .filter((f) => f.source_id === nodeId || f.target_id === nodeId)
      .sort(
        (a, b) =>
          (SEVERITY_RANK[a.severity] ?? 4) - (SEVERITY_RANK[b.severity] ?? 4) ||
          b.confidence - a.confidence,
      );
  }, [findings, nodeId]);

  if (isLoading) {
    return (
      <div className="font-mono text-xs uppercase tracking-[0.1em] text-muted-foreground">
        Loading findings…
      </div>
    );
  }

  if (related.length === 0) {
    return (
      <div className="font-mono text-xs uppercase tracking-[0.1em] text-muted-foreground">
        No findings reference this node.
      </div>
    );
  }

  return (
    <div className="space-y-1.5">
      <div className="mb-2 font-mono text-[10px] uppercase tracking-[0.12em] text-muted-foreground">
        {related.length} finding{related.length === 1 ? "" : "s"} on this node
      </div>
      {related.map((f) => {
        const sev = SEVERITY_BY_KEY[f.severity] ?? SEVERITY.low;
        const direction = f.source_id === nodeId ? "outbound" : "inbound";
        return (
          <button
            key={f.id}
            onClick={() => navigate(`/findings/${f.id}`)}
            className="flex w-full items-center gap-2.5 rounded-[3px] border border-border bg-black/30 px-3 py-2 text-left transition-colors hover:border-mauve-7 hover:bg-white/[0.03]"
          >
            <span
              className="h-2 w-2 flex-shrink-0 rounded-[1px]"
              style={{ backgroundColor: sev.solid, boxShadow: `0 0 6px -1px ${sev.solid}` }}
            />
            <div className="flex min-w-0 flex-1 flex-col">
              <span className="truncate text-xs font-medium text-foreground">{f.title}</span>
              <span className="truncate font-mono text-[10px] text-muted-foreground">
                {edgeLabel(f.edge_kind)} · {direction} ·{" "}
                {Math.round(f.confidence * 100)}% conf
              </span>
            </div>
            <span
              className="flex-shrink-0 font-mono text-[9px] font-bold uppercase tracking-[0.08em]"
              style={{ color: sev.solid }}
            >
              {f.severity}
            </span>
            <ArrowRight className="h-3.5 w-3.5 flex-shrink-0 text-primary/50" />
          </button>
        );
      })}
    </div>
  );
}
