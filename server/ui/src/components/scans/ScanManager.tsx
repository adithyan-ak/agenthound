import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { ScanSearch, Plus, Upload } from "lucide-react";
import { fetchScans } from "@/api/scans";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ScanHistory } from "./ScanHistory";
import { NewScan } from "./NewScan";
import { ScanImport } from "./ScanImport";

export function ScanManager() {
  const [showNewScan, setShowNewScan] = useState(false);
  const [showImport, setShowImport] = useState(false);

  const { data: scans, isLoading, refetch } = useQuery({
    queryKey: ["scans"],
    queryFn: () => fetchScans(),
    staleTime: 30_000,
  });

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="flex items-center gap-2 text-lg font-semibold text-foreground">
          <ScanSearch className="h-5 w-5 text-primary" />
          Scan Manager
        </h2>
        <div className="flex gap-2">
          <Button
            onClick={() => setShowImport(true)}
            size="sm"
            variant="outline"
          >
            <Upload className="h-4 w-4 mr-1.5" />
            Import scan
          </Button>
          <Button onClick={() => setShowNewScan(true)} size="sm">
            <Plus className="h-4 w-4 mr-1.5" />
            New Scan
          </Button>
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="space-y-2 p-4">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-3/4" />
            </div>
          ) : (
            <ScanHistory scans={scans ?? []} onDeleted={() => refetch()} />
          )}
        </CardContent>
      </Card>

      <NewScan open={showNewScan} onClose={() => setShowNewScan(false)} />
      <ScanImport
        open={showImport}
        onClose={() => setShowImport(false)}
        onSuccess={() => refetch()}
      />
    </div>
  );
}
