import { Link } from "react-router-dom";
import { ScanSearch, ArrowRight } from "lucide-react";
import { useGraphStats } from "@/hooks/useGraph";
import { DashboardHeader } from "./DashboardHeader";
import { StatCards } from "./StatCards";
import { ExposureGauge } from "./ExposureGauge";
import { SeverityRings } from "./SeverityRings";
import { CategoryBreakdown } from "./CategoryBreakdown";
import { AuthCoverage } from "./AuthCoverage";
import { InventoryTrend } from "./InventoryTrend";
import { TopRiskyEntities } from "./TopRiskyEntities";
import { TopFindings } from "./TopFindings";
import { CrossProtocol } from "./CrossProtocol";
import { Chokepoints } from "./Chokepoints";
import { RecentScans } from "./RecentScans";

function EmptyState() {
  return (
    <div className="card-elevated mt-6 flex flex-col items-center justify-center gap-4 rounded-xl px-6 py-16 text-center">
      <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/10 ring-1 ring-inset ring-primary/20">
        <ScanSearch className="h-7 w-7 text-primary" />
      </div>
      <div className="space-y-1.5">
        <h2 className="text-lg font-semibold text-foreground">No attack surface mapped yet</h2>
        <p className="max-w-md text-sm text-muted-foreground">
          Run a scan to discover your agent, MCP, and A2A infrastructure. Once ingested, your
          exposure index, findings, and attack paths will appear here.
        </p>
      </div>
      <code className="rounded-md border border-white/[0.08] bg-black/40 px-3 py-1.5 font-mono text-sm text-primary">
        agenthound scan
      </code>
      <Link
        to="/scans"
        className="inline-flex items-center gap-1.5 text-sm font-medium text-primary transition-colors hover:text-primary/80"
      >
        Go to Scans <ArrowRight className="h-3.5 w-3.5" />
      </Link>
    </div>
  );
}

const ROW = "animate-fade-up";

export function Dashboard() {
  const { data: stats, isLoading } = useGraphStats();
  const isEmpty = !isLoading && (stats?.total_nodes ?? 0) === 0;

  return (
    <div className="dashboard-bg min-h-full p-4 sm:p-6">
      <div className="mx-auto max-w-[1600px] space-y-5">
        <DashboardHeader />

        {isEmpty ? (
          <EmptyState />
        ) : (
          <>
            <div className={ROW} style={{ animationDelay: "40ms" }}>
              <StatCards />
            </div>

            <div className={`grid gap-4 lg:grid-cols-3 ${ROW}`} style={{ animationDelay: "100ms" }}>
              <div className="lg:col-span-1">
                <ExposureGauge />
              </div>
              <div className="grid gap-4 lg:col-span-2">
                <SeverityRings />
                <InventoryTrend />
              </div>
            </div>

            <div className={`grid gap-4 lg:grid-cols-3 ${ROW}`} style={{ animationDelay: "160ms" }}>
              <div className="lg:col-span-2">
                <CategoryBreakdown />
              </div>
              <div className="lg:col-span-1">
                <AuthCoverage />
              </div>
            </div>

            <div className={`grid gap-4 lg:grid-cols-2 ${ROW}`} style={{ animationDelay: "220ms" }}>
              <TopRiskyEntities />
              <TopFindings />
            </div>

            <div className={`grid gap-4 lg:grid-cols-2 ${ROW}`} style={{ animationDelay: "280ms" }}>
              <CrossProtocol />
              <Chokepoints />
            </div>

            <div className={ROW} style={{ animationDelay: "340ms" }}>
              <RecentScans />
            </div>
          </>
        )}
      </div>
    </div>
  );
}
