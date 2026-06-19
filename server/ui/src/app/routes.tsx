import { lazy, Suspense } from "react";
import { Routes, Route } from "react-router-dom";
import { AppLayout } from "./layout";
import { Dashboard } from "@features/dashboard";
import { ScanManager } from "@features/scans";
import { QueryLibrary } from "@features/queries";
import { RulesLibrary } from "@features/rules";

const ExplorerPage = lazy(() =>
  import("@features/explorer").then((m) => ({
    default: m.ExplorerPage,
  })),
);

const FindingsListPage = lazy(() =>
  import("@features/findings/ui/FindingsListPage").then((m) => ({
    default: m.FindingsListPage,
  })),
);

const FindingDetailPage = lazy(() =>
  import("@features/findings/ui/FindingDetailPage").then((m) => ({
    default: m.FindingDetailPage,
  })),
);

function ExplorerFallback() {
  return (
    <div className="flex h-full items-center justify-center">
      <p className="text-sm text-muted-foreground">Loading Explorer…</p>
    </div>
  );
}

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<Dashboard />} />
        <Route
          path="/explorer"
          element={
            <Suspense fallback={<ExplorerFallback />}>
              <ExplorerPage />
            </Suspense>
          }
        />
        <Route
          path="/findings"
          element={
            <Suspense fallback={<div className="flex h-full items-center justify-center"><p className="text-sm text-muted-foreground">Loading Findings…</p></div>}>
              <FindingsListPage />
            </Suspense>
          }
        />
        <Route
          path="/findings/:findingId"
          element={
            <Suspense fallback={<div className="flex h-full items-center justify-center"><p className="text-sm text-muted-foreground">Loading Finding…</p></div>}>
              <FindingDetailPage />
            </Suspense>
          }
        />

        <Route path="/scans" element={<ScanManager />} />
        <Route path="/queries" element={<QueryLibrary />} />
        <Route path="/rules" element={<RulesLibrary />} />
      </Route>
    </Routes>
  );
}
