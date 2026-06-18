import { useParams, useNavigate } from "react-router-dom";
import { useEffect } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Stack, Sidebar, Center } from "@/components/ui/layout";
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
      <Stack className="p-6">
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-48 w-full" />
        <Sidebar
          sidePosition="right"
          sideWidth="22rem"
          contentMin="58%"
          side={
            <Stack gap="1rem">
              <Skeleton className="h-32 w-full" />
              <Skeleton className="h-48 w-full" />
            </Stack>
          }
          main={<Skeleton className="h-64 w-full" />}
        />
      </Stack>
    );
  }

  if (error || !detail) {
    return (
      <Center
        measure="32rem"
        className="flex flex-col items-center justify-center h-full gap-4 p-6 text-center"
      >
        <div className="text-lg font-semibold text-foreground">Finding not found</div>
        <p className="text-sm text-muted-foreground">
          This finding may have been resolved in a recent scan, or the scan data may have been cleared.
        </p>
        <button
          onClick={() => navigate("/findings")}
          className="text-sm text-primary hover:underline"
        >
          Back to Findings
        </button>
      </Center>
    );
  }

  const f = detail.finding;

  function handleCopyReport() {
    const md = buildMarkdownReport(f, detail!.attack_path, detail!.remediation);
    navigator.clipboard.writeText(md);
  }

  return (
    <Center measure="84rem" className="p-6">
      <Stack gap="1.5rem">
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
            <Stack gap="1rem">
              <FindingImpact impact={detail.impact} path={detail.attack_path} />
              <FindingRemediation steps={detail.remediation} />
              <FindingReferences finding={f} />
            </Stack>
          }
          main={<HopEvidenceTimeline path={detail.attack_path} />}
        />
      </Stack>
    </Center>
  );
}
