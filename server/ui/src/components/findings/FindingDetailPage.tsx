import { useParams, useNavigate } from "react-router-dom";
import { useEffect } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Stack, Sidebar } from "@/components/ui/layout";
import { useFindingDetail } from "@/hooks/useFindingDetail";
import { useFindingsNavigation } from "@/hooks/useFindingsNavigation";
import { buildMarkdownReport } from "@/lib/findings/copy-report";
import { FindingHeader } from "./FindingHeader";
import { AttackPathDiagram } from "./AttackPathDiagram";
import { HopEvidenceTimeline } from "./HopEvidenceTimeline";
import { FindingImpact } from "./FindingImpact";
import { FindingRemediation } from "./FindingRemediation";
import { FindingReferences } from "./FindingReferences";

export function FindingDetailPage() {
  const { findingId } = useParams<{ findingId: string }>();
  const navigate = useNavigate();
  const { data: detail, isLoading, error } = useFindingDetail(findingId);
  const { prevId, nextId } = useFindingsNavigation(findingId);

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;
      if (e.key === "ArrowLeft" && prevId) navigate(`/findings/${prevId}`);
      if (e.key === "ArrowRight" && nextId) navigate(`/findings/${nextId}`);
      if (e.key === "Escape") navigate("/findings");
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [prevId, nextId, navigate]);

  if (isLoading) {
    return (
      <div className="dashboard-bg min-h-full p-3 sm:p-4 lg:p-5">
        <div className="mx-auto max-w-[88rem]">
          <Stack gap="0.75rem">
            <Skeleton className="h-36 w-full rounded-md" />
            <Skeleton className="h-44 w-full rounded-md" />
            <Sidebar
              sidePosition="right"
              sideWidth="22rem"
              contentMin="58%"
              side={
                <Stack gap="0.75rem">
                  <Skeleton className="h-32 w-full rounded-md" />
                  <Skeleton className="h-48 w-full rounded-md" />
                </Stack>
              }
              main={<Skeleton className="h-64 w-full rounded-md" />}
            />
          </Stack>
        </div>
      </div>
    );
  }

  if (error || !detail) {
    return (
      <div className="dashboard-bg flex min-h-full flex-col items-center justify-center gap-4 p-6 text-center">
        <div className="font-mono text-sm font-semibold uppercase tracking-[0.12em] text-foreground">
          Finding not found
        </div>
        <p className="max-w-md text-sm text-muted-foreground">
          This finding may have been resolved in a recent scan, or the scan data may have been cleared.
        </p>
        <button
          onClick={() => navigate("/findings")}
          className="font-mono text-xs uppercase tracking-[0.08em] text-primary transition-colors hover:text-primary/80"
        >
          &#9656; Back to Findings
        </button>
      </div>
    );
  }

  const f = detail.finding;

  function handleCopyReport() {
    const md = buildMarkdownReport(f, detail!.attack_path, detail!.remediation);
    navigator.clipboard.writeText(md);
  }

  return (
    <div className="dashboard-bg min-h-full p-3 sm:p-4 lg:p-5">
      <div className="mx-auto max-w-[88rem]">
        <Stack gap="0.75rem">
          <FindingHeader
            detail={detail}
            prevId={prevId}
            nextId={nextId}
            onCopyReport={handleCopyReport}
          />

          <AttackPathDiagram
            path={detail.attack_path}
            severity={f.severity}
            sourceId={f.source_id}
            sourceName={f.source_name}
            sourceKind={f.source_kind}
            targetId={f.target_id}
            targetName={f.target_name}
            targetKind={f.target_kind}
          />

          <Sidebar
            sidePosition="right"
            sideWidth="22rem"
            contentMin="58%"
            side={
              <Stack gap="0.75rem">
                <FindingImpact impact={detail.impact} path={detail.attack_path} />
                <FindingRemediation steps={detail.remediation} />
                <FindingReferences finding={f} />
              </Stack>
            }
            main={<HopEvidenceTimeline path={detail.attack_path} />}
          />
        </Stack>
      </div>
    </div>
  );
}
