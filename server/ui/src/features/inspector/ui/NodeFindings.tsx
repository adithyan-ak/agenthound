import type { CSSProperties } from "react";
import { AlertTriangle } from "lucide-react";
import { useFindings } from "@entities/finding";
import { useGraphStore } from "../model/graph-store";
import { SEVERITY_BY_KEY } from "@shared/theme/tokens";
import { Card, CardContent } from "@shared/ui/primitives/card";
import { Badge } from "@shared/ui/primitives/badge";

interface NodeFindingsProps {
  nodeId: string;
}

const SEVERITY_BADGE_VARIANT: Record<
  string,
  "destructive" | "default" | "secondary" | "outline"
> = {
  critical: "destructive",
  high: "destructive",
  medium: "default",
  low: "secondary",
  info: "outline",
};

function severityCardInlineStyle(level: string): CSSProperties {
  const s = SEVERITY_BY_KEY[level];
  if (!s) return {};
  return { borderColor: s.border, backgroundColor: s.bg };
}

export function NodeFindings({ nodeId }: NodeFindingsProps) {
  const highlightPath = useGraphStore((s) => s.highlightPath);

  const { data: allFindings, isLoading } = useFindings();

  if (isLoading) {
    return (
      <div className="py-4 text-sm text-muted-foreground text-center">
        Loading...
      </div>
    );
  }

  const findings = (allFindings ?? []).filter(
    (f) => f.source_id === nodeId || f.target_id === nodeId,
  );

  if (findings.length === 0) {
    return (
      <div className="py-4 text-sm text-muted-foreground text-center">
        No findings for this node
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {findings.map((finding) => (
        <button
          key={finding.id}
          onClick={() => {
            const edgeKey = `${finding.source_id}->${finding.target_id}:${finding.edge_kind}`;
            highlightPath({
              nodeIds: [finding.source_id, finding.target_id],
              edgeKeys: [edgeKey],
              title: finding.title,
            });
          }}
          className="w-full text-left"
          title="Show this path on the graph"
        >
          <Card
            className="transition-colors hover:brightness-125 cursor-pointer"
            style={severityCardInlineStyle(finding.severity)}
          >
            <CardContent className="px-3 py-2">
              <div className="flex items-start gap-2">
                <AlertTriangle className="h-3.5 w-3.5 mt-0.5 flex-shrink-0" />
                <div className="min-w-0">
                  <div className="flex items-center gap-2 mb-0.5">
                    <Badge
                      variant={
                        SEVERITY_BADGE_VARIANT[finding.severity] ?? "outline"
                      }
                      className="text-[10px] uppercase"
                    >
                      {finding.severity}
                    </Badge>
                    <span className="text-xs font-medium">{finding.title}</span>
                  </div>
                  <p className="text-xs opacity-80">{finding.description}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </button>
      ))}
    </div>
  );
}
