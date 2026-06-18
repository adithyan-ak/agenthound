import type { APINode } from "@/api/types";
import { Grid } from "@/components/ui/layout";

const HIDDEN_KEYS = new Set([
  "objectid",
  "scan_id",
  "last_seen",
  "created_at",
  "description_hash",
  "card_hash",
]);

export function PropertiesTab({ node }: { node: APINode }) {
  const entries = Object.entries(node.properties ?? {}).filter(
    ([k, v]) => !HIDDEN_KEYS.has(k) && v !== null && v !== undefined && v !== "",
  );

  if (entries.length === 0) {
    return (
      <div className="text-sm text-muted-foreground">No properties recorded.</div>
    );
  }

  return (
    <Grid min="14rem" gap="0.75rem 2rem">
      {entries.map(([k, v]) => (
        <div key={k} className="flex flex-col gap-0.5 min-w-0">
          <div className="text-[10px] uppercase tracking-wider text-muted-foreground font-medium">
            {k.replace(/_/g, " ")}
          </div>
          <div className="truncate text-[13px] text-foreground font-mono">
            {renderValue(v)}
          </div>
        </div>
      ))}
    </Grid>
  );
}

function renderValue(v: unknown): string {
  if (v === null || v === undefined) return "—";
  if (typeof v === "boolean") return v ? "true" : "false";
  if (typeof v === "number") return String(v);
  if (typeof v === "string") return v;
  if (Array.isArray(v)) return v.map((x) => String(x)).join(", ");
  if (typeof v === "object") return JSON.stringify(v);
  return String(v);
}
