import { ReactFlowProvider } from "@xyflow/react";
import { ExplorerCanvas } from "./ExplorerCanvas";

export function ExplorerPage() {
  return (
    <div className="relative h-full w-full overflow-hidden bg-[#050B18]">
      <ReactFlowProvider>
        <ExplorerCanvas />
      </ReactFlowProvider>
    </div>
  );
}
