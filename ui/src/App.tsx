import { useEffect } from "react";
import { Routes, Route } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { Dashboard } from "@/components/dashboard/Dashboard";
import { GraphExplorer } from "@/components/graph/GraphExplorer";
import { Pathfinder } from "@/components/pathfinder/Pathfinder";
import { ScanManager } from "@/components/scans/ScanManager";
import { QueryLibrary } from "@/components/queries/QueryLibrary";
import { LoginPage } from "@/components/auth/LoginPage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useAuthStore } from "@/store/auth";

export function App() {
  const initialize = useAuthStore((s) => s.initialize);

  useEffect(() => {
    initialize();
  }, [initialize]);

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<AppLayout />}>
          <Route path="/" element={<Dashboard />} />
          <Route path="/graph" element={<GraphExplorer />} />
          <Route path="/pathfinder" element={<Pathfinder />} />
          <Route path="/scans" element={<ScanManager />} />
          <Route path="/queries" element={<QueryLibrary />} />
        </Route>
      </Route>
    </Routes>
  );
}
