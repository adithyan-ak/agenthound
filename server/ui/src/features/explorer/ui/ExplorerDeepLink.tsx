import { useEffect, useRef } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Crosshair, ExternalLink, X } from "lucide-react";
import { useFindingDetail } from "@entities/finding";
import { useExplorerStore } from "@features/explorer/model/store";
import { lensForEdgeKind } from "@features/explorer/model/lens-config";
import { cn } from "@shared/lib/utils";
import { SEVERITY, SEVERITY_BY_KEY } from "@shared/theme/tokens";

/**
 * Reads `?finding=<id>` from the URL and focuses the explorer on that finding's
 * attack path — selects the framing lens, highlights the path nodes/edges, and
 * pans the camera onto them. This is the context-carrying half of the
 * dossier ↔ map handoff (the dossier's "View in Explorer" links here). Renders
 * a dismissible focus banner while active. The param is preserved so the view
 * is shareable/bookmarkable; the focus is applied once per finding id.
 */
export function ExplorerDeepLink() {
  const [params, setParams] = useSearchParams();
  const navigate = useNavigate();
  const findingId = params.get("finding");

  const { data: detail } = useFindingDetail(findingId ?? undefined);

  const setActiveLens = useExplorerStore((s) => s.setActiveLens);
  const setHighlight = useExplorerStore((s) => s.setHighlight);
  const setPendingFocus = useExplorerStore((s) => s.setPendingFocus);
  const clearHighlight = useExplorerStore((s) => s.clearHighlight);

  const appliedRef = useRef<string | null>(null);

  useEffect(() => {
    if (!findingId || !detail) return;
    if (appliedRef.current === findingId) return;
    appliedRef.current = findingId;

    const f = detail.finding;
    const path = detail.attack_path;
    const nodeIds =
      path && path.nodes.length > 0
        ? path.nodes.map((n) => n.id)
        : [f.source_id, f.target_id];
    const edgeIds = path
      ? path.edges.map((e) => `${e.source}|${e.target}|${e.kind}`)
      : [];

    setActiveLens(lensForEdgeKind(f.edge_kind));
    setHighlight({ nodeIds, edgeIds, title: f.title });
    setPendingFocus({ nodeIds, title: f.title });
  }, [findingId, detail, setActiveLens, setHighlight, setPendingFocus]);

  if (!findingId || !detail) return null;

  const f = detail.finding;
  const sev = SEVERITY_BY_KEY[f.severity] ?? SEVERITY.low;

  function clear() {
    clearHighlight();
    appliedRef.current = null;
    const next = new URLSearchParams(params);
    next.delete("finding");
    setParams(next, { replace: true });
  }

  return (
    <div className="pointer-events-auto absolute left-1/2 top-16 z-30 -translate-x-1/2">
      <div
        className={cn(
          "flex items-center gap-3 overflow-hidden rounded-md border border-border bg-card/95 px-3 py-2 backdrop-blur-md elev-2",
        )}
        style={{ boxShadow: `inset 2px 0 0 0 ${sev.solid}` }}
      >
        <span className="flex items-center gap-1.5">
          <Crosshair className="h-3.5 w-3.5" style={{ color: sev.solid }} strokeWidth={2.25} />
          <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
            Focused on finding
          </span>
        </span>
        <span className="max-w-[320px] truncate text-xs font-medium text-foreground">
          {f.title}
        </span>
        <button
          onClick={() => navigate(`/findings/${f.id}`)}
          className="inline-flex items-center gap-1 rounded-[3px] border border-border bg-black/30 px-2 py-1 font-mono text-[10px] uppercase tracking-[0.06em] text-foreground/80 transition-colors hover:border-primary/50 hover:text-primary"
        >
          <ExternalLink className="h-3 w-3" /> Open finding
        </button>
        <button
          onClick={clear}
          className="flex h-6 w-6 items-center justify-center rounded-[3px] text-muted-foreground transition-colors hover:bg-white/[0.06] hover:text-foreground"
          aria-label="Clear focus"
        >
          <X className="h-3.5 w-3.5" />
        </button>
      </div>
    </div>
  );
}
