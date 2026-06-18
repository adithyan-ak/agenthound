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
      <div className="font-mono text-xs uppercase tracking-[0.1em] text-muted-foreground">
        No properties recorded.
      </div>
    );
  }

  return (
    <Grid min="14rem" gap="0.75rem 2rem">
      {entries.map(([k, v]) => (
        <div key={k} className="flex min-w-0 flex-col gap-0.5">
          <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-muted-foreground">
            {k.replace(/_/g, " ")}
          </div>
          <div className="truncate font-mono text-[13px] text-foreground/90">{renderValue(v)}</div>
        </div>
      ))}
    </Grid>
  );
}

function renderValue(v: unknown): string {
  if (v === null || v === undefined) return "\u2014";
  if (typeof v === "boolean") return v ? "true" : "false";
  if (typeof v === "number") return String(v);
  if (typeof v === "string") return v;
  if (Array.isArray(v)) return v.map((x) => String(x)).join(", ");
  if (typeof v === "object") return JSON.stringify(v);
  return String(v);
}
