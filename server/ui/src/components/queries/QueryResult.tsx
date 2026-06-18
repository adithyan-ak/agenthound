import type { PreBuiltQuery } from "@/api/types";

interface QueryResultProps {
  rows: Record<string, unknown>[];
  query: PreBuiltQuery;
}

function formatCell(value: unknown): string {
  if (value == null) return "\u2014";
  if (typeof value === "boolean") return value ? "Yes" : "No";
  if (Array.isArray(value)) return value.join(", ") || "\u2014";
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

export function QueryResult({ rows, query }: QueryResultProps) {
  if (rows.length === 0) {
    return (
      <div className="py-4 text-center font-mono text-xs uppercase tracking-[0.1em] text-muted-foreground">
        No results for "{query.name}"
      </div>
    );
  }

  const columns = Object.keys(rows[0]!);

  return (
    <div>
      <div className="mb-2 font-mono text-[10px] uppercase tracking-[0.12em] text-muted-foreground">
        {rows.length} row{rows.length !== 1 ? "s" : ""}
      </div>
      <div className="overflow-x-auto rounded-[3px] border border-border/70">
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-border bg-black/30">
              {columns.map((col) => (
                <th
                  key={col}
                  className="px-3 py-1.5 font-mono text-[10px] font-semibold uppercase tracking-[0.1em] text-muted-foreground"
                >
                  {col}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr
                key={i}
                className="border-b border-border/50 transition-colors last:border-0 hover:bg-white/[0.03]"
              >
                {columns.map((col) => (
                  <td
                    key={col}
                    className="max-w-[300px] truncate px-3 py-1.5 font-mono text-[11px] text-foreground/90"
                  >
                    {formatCell(row[col])}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
