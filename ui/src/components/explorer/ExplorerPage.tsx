import { ReactFlowProvider } from "@xyflow/react";
import { ExplorerCanvas } from "./ExplorerCanvas";
import { LensBar } from "./LensBar";
import { InfoCard } from "./InfoCard";
import { Legend } from "./Legend";
import { StatusStrip } from "./StatusStrip";

export function ExplorerPage() {
  return (
    <div className="relative h-full w-full overflow-hidden bg-[#050B18]">
      <ReactFlowProvider>
        <ExplorerCanvas />
        <LensBar />
        <InfoCard />
        <Legend />
        <StatusStrip />
      </ReactFlowProvider>
    </div>
  );
}
