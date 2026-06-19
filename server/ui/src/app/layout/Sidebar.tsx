import { X } from "lucide-react";
import { useUIStore } from "@shared/model/ui-store";
import { EntityInspector } from "@features/inspector";
import { ScrollArea } from "@shared/ui/primitives/scroll-area";
import { Button } from "@shared/ui/primitives/button";

export function Sidebar() {
  const closeSidebar = useUIStore((s) => s.closeSidebar);

  return (
    <aside className="flex w-[380px] flex-shrink-0 flex-col border-l border-border bg-carbon-900">
      <div className="flex items-center justify-between border-b border-border px-4 py-2">
        <span className="font-mono text-[11px] font-semibold uppercase tracking-[0.16em] text-muted-foreground">
          <span className="text-primary/70">[</span> Inspector <span className="text-primary/70">]</span>
        </span>
        <Button onClick={closeSidebar} variant="ghost" size="icon" className="h-7 w-7 rounded-[3px]">
          <X className="h-4 w-4" />
        </Button>
      </div>
      <ScrollArea className="h-full">
        <EntityInspector />
      </ScrollArea>
    </aside>
  );
}
