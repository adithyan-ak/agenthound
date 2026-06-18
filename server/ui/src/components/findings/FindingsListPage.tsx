import { useMemo, useState, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { ArrowRight, Search, ShieldAlert } from "lucide-react";
import { fetchFindings } from "@/api/analysis";
import type { Finding } from "@/api/types";
import { Skeleton } from "@/components/ui/skeleton";
import { MiniHexIcon } from "@/components/findings/MiniHexIcon";
import { MeterBar, WidgetCard } from "@/components/dashboard/kit";
import { cn } from "@/lib/utils";
import { ACCENT, SEVERITY, SEVERITY_BY_KEY, severityColor } from "@/theme/tokens";

const SEVERITY_RANK: Record<string, number> = { critical: 0, high: 1, medium: 2, low: 3 };
const SEVERITY_LEVELS = ["critical", "high", "medium", "low"] as const;

function pad2(n: number): string {
  return String(n).padStart(2, "0");
}

export function FindingsListPage() {
  const navigate = useNavigate();
  const [activeSeverities, setActiveSeverities] = useState<Set<string>>(new Set());
  const [search, setSearch] = useState("");

  const { data: findings, isLoading } = useQuery({
    queryKey: ["findings"],
    queryFn: () => fetchFindings(),
    staleTime: 30_000,
  });

  function toggleSeverity(level: string) {
    setActiveSeverities((prev) => {
      const next = new Set(prev);
      if (next.has(level)) next.delete(level);
      else next.add(level);
      return next;
    });
  }

  const counts = useMemo(() => {
    const c: Record<string, number> = { critical: 0, high: 0, medium: 0, low: 0 };
    for (const f of findings ?? []) c[f.severity] = (c[f.severity] ?? 0) + 1;
    return c;
  }, [findings]);

  const filtered = useMemo(() => {
    if (!findings) return [];
    const q = search.toLowerCase();
    return findings
      .filter((f) => {
        if (activeSeverities.size > 0 && !activeSeverities.has(f.severity)) return false;
        if (q) {
          return (
            f.title.toLowerCase().includes(q) ||
            f.description.toLowerCase().includes(q) ||
            f.source_name.toLowerCase().includes(q) ||
            f.target_name.toLowerCase().includes(q)
          );
        }
        return true;
      })
      .sort(
        (a, b) =>
          (SEVERITY_RANK[a.severity] ?? 4) - (SEVERITY_RANK[b.severity] ?? 4) ||
          b.confidence - a.confidence,
      );
  }, [findings, activeSeverities, search]);

  const total = findings?.length ?? 0;

  return (
    <div className="dashboard-bg min-h-full p-3 sm:p-4 lg:p-5">
      <div className="mx-auto max-w-[1600px] space-y-3">
        {/* ---------- Command header ---------- */}
        <header className="space-y-3">
          <div className="min-w-0">
            <p className="font-mono text-[10px] font-semibold uppercase tracking-[0.18em] text-muted-foreground">
              Security Findings <span className="text-primary/60">//</span> Threat Register
            </p>
            <h1 className="mt-1.5 flex items-center gap-2.5 font-mono text-2xl font-bold uppercase tracking-[0.04em] text-foreground sm:text-[26px]">
              <span className="flex h-7 w-7 items-center justify-center rounded-[3px] bg-primary/10 ring-1 ring-inset ring-primary/30">
                <ShieldAlert className="h-4 w-4 text-primary" />
              </span>
              <span className="text-primary">&#9656;</span>
              Findings
              <span className="font-mono text-base font-semibold tabular-nums text-muted-foreground">
                {pad2(total)}
              </span>
              <span className="blink-caret text-primary" aria-hidden>
                _
              </span>
            </h1>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Detected attack paths and exposures across your agent, MCP, and A2A infrastructure.
            </p>
          </div>

          {/* SOC register summary strip */}
          <div className="card-elevated relative flex flex-wrap items-center overflow-hidden rounded-md">
            <span aria-hidden className="absolute left-0 top-0 h-px w-14 bg-primary/80" />
            <div className="flex flex-wrap items-center divide-x divide-border/70">
              <RegisterSeg label="Total" value={pad2(total)} color="#E9ECF0" />
              {SEVERITY_LEVELS.map((lvl) => {
                const n = counts[lvl] ?? 0;
                return (
                  <RegisterSeg
                    key={lvl}
                    label={SEVERITY[lvl].label}
                    value={pad2(n)}
                    color={severityColor(lvl)}
                    dim={n === 0}
                  />
                );
              })}
            </div>
            <div className="ml-auto flex items-center gap-2 self-stretch border-l border-border/70 px-3.5 py-2">
              <span className="font-mono text-[10px] uppercase tracking-[0.16em] text-muted-foreground">
                Showing
              </span>
              <span className="font-mono text-[11px] font-semibold tabular-nums text-primary">
                {pad2(filtered.length)}
              </span>
              <span className="font-mono text-[10px] tabular-nums text-muted-foreground">
                / {pad2(total)}
              </span>
            </div>
          </div>
        </header>

        {/* ---------- Toolbar: filters + search ---------- */}
        <div className="flex flex-wrap items-center gap-2">
          <FilterChip
            active={activeSeverities.size === 0}
            onClick={() => setActiveSeverities(new Set())}
            color={ACCENT}
            label="All"
          />
          {SEVERITY_LEVELS.map((lvl) => (
            <FilterChip
              key={lvl}
              active={activeSeverities.has(lvl)}
              onClick={() => toggleSeverity(lvl)}
              color={severityColor(lvl)}
              label={SEVERITY[lvl].label}
            />
          ))}

          <div className="relative ml-auto">
            <Search className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
            <input
              type="text"
              placeholder="search findings…"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="h-8 w-72 rounded-[3px] border border-border bg-black/30 pl-8 pr-3 font-mono text-xs text-foreground placeholder:text-muted-foreground/70 focus-visible:border-primary/50 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-primary/40"
            />
          </div>
        </div>

        {/* ---------- Register table ---------- */}
        {isLoading ? (
          <WidgetCard title="Threat Register" icon={ShieldAlert} flush>
            <div className="space-y-px p-2">
              {Array.from({ length: 8 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full rounded-[2px]" />
              ))}
            </div>
          </WidgetCard>
        ) : filtered.length === 0 ? (
          <WidgetCard title="Threat Register" icon={ShieldAlert}>
            <div className="flex flex-col items-center justify-center gap-3 py-16 text-center">
              <span className="flex h-12 w-12 items-center justify-center rounded-[4px] bg-primary/10 ring-1 ring-inset ring-primary/20">
                <ShieldAlert className="h-6 w-6 text-primary" />
              </span>
              <p className="font-mono text-xs uppercase tracking-[0.12em] text-muted-foreground">
                {total > 0 ? "No findings match the current filters" : "No findings detected"}
              </p>
              {total === 0 && (
                <button
                  onClick={() => navigate("/scans")}
                  className="font-mono text-xs uppercase tracking-[0.08em] text-primary transition-colors hover:text-primary/80"
                >
                  &#9656; Go to Scans
                </button>
              )}
            </div>
          </WidgetCard>
        ) : (
          <WidgetCard
            title="Threat Register"
            icon={ShieldAlert}
            flush
            action={
              <span className="font-mono text-[10px] uppercase tracking-[0.12em] text-muted-foreground">
                {pad2(filtered.length)} <span className="text-muted-foreground/50">/</span> {pad2(total)}
              </span>
            }
          >
            <div className="overflow-x-auto">
              <table className="w-full border-collapse text-left">
                <thead>
                  <tr className="border-b border-border bg-black/20">
                    <Th className="w-10 pr-2 text-right">#</Th>
                    <Th>Severity</Th>
                    <Th>Finding</Th>
                    <Th>Flow</Th>
                    <Th>Category</Th>
                    <Th>OWASP</Th>
                    <Th className="text-right">Confidence</Th>
                  </tr>
                </thead>
                <tbody>
                  {filtered.map((f, i) => (
                    <FindingRow
                      key={f.id}
                      index={i}
                      finding={f}
                      onClick={() => navigate(`/findings/${f.id}`)}
                    />
                  ))}
                </tbody>
              </table>
            </div>
          </WidgetCard>
        )}
      </div>
    </div>
  );
}

function RegisterSeg({
  label,
  value,
  color,
  dim,
}: {
  label: string;
  value: string;
  color: string;
  dim?: boolean;
}) {
  return (
    <div className="flex items-center gap-2 px-3 py-2">
      <span
        className="h-2 w-2 shrink-0 rounded-[1px]"
        style={{
          backgroundColor: dim ? "rgb(var(--mauve-8-raw))" : color,
          boxShadow: dim ? undefined : `0 0 6px -1px ${color}`,
        }}
      />
      <span className="font-mono text-[10px] uppercase tracking-[0.14em] text-muted-foreground">
        {label}
      </span>
      <span
        className="font-mono text-sm font-bold tabular-nums"
        style={{ color: dim ? "rgb(var(--mauve-9-raw))" : color }}
      >
        {value}
      </span>
    </div>
  );
}

function FilterChip({
  active,
  onClick,
  color,
  label,
}: {
  active: boolean;
  onClick: () => void;
  color: string;
  label: string;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "inline-flex items-center gap-1.5 rounded-[3px] border px-2.5 py-1 font-mono text-[11px] font-medium uppercase tracking-[0.08em] transition-colors",
        active
          ? "border-transparent text-foreground"
          : "border-border bg-black/30 text-muted-foreground hover:border-mauve-7 hover:text-foreground",
      )}
      style={active ? { backgroundColor: `${color}1A`, boxShadow: `inset 0 0 0 1px ${color}66` } : undefined}
    >
      <span
        className="h-2 w-2 rounded-[1px]"
        style={{ backgroundColor: active ? color : "rgb(var(--mauve-8-raw))" }}
      />
      {label}
    </button>
  );
}

function Th({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <th
      className={cn(
        "px-3 py-2 font-mono text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground",
        className,
      )}
    >
      {children}
    </th>
  );
}

function FindingRow({
  index,
  finding: f,
  onClick,
}: {
  index: number;
  finding: Finding;
  onClick: () => void;
}) {
  const color = (SEVERITY_BY_KEY[f.severity] ?? SEVERITY.low).solid;
  const conf = Math.round(f.confidence * 100);

  return (
    <tr
      onClick={onClick}
      className="cursor-pointer border-b border-border/60 transition-colors last:border-0 hover:bg-white/[0.03]"
    >
      <td
        className="px-3 py-3 text-right align-middle font-mono text-[10px] tabular-nums text-muted-foreground/60"
        style={{ boxShadow: `inset 2px 0 0 0 ${color}` }}
      >
        {pad2(index + 1)}
      </td>
      <td className="px-3 py-3 align-middle">
        <span className="inline-flex items-center gap-1.5">
          <span
            className="h-2.5 w-2.5 rounded-[1px]"
            style={{ backgroundColor: color, boxShadow: `0 0 6px -1px ${color}` }}
          />
          <span
            className="font-mono text-[10px] font-bold uppercase tracking-[0.1em]"
            style={{ color }}
          >
            {f.severity}
          </span>
        </span>
      </td>
      <td className="max-w-sm px-3 py-3 align-middle">
        <p className="truncate text-[13px] font-medium text-foreground">{f.title}</p>
        <p className="truncate text-xs text-muted-foreground">{f.description}</p>
      </td>
      <td className="px-3 py-3 align-middle">
        <div className="flex items-center gap-1.5 font-mono text-[11px]">
          <MiniHexIcon kind={f.source_kind} />
          <span className="max-w-[110px] truncate text-foreground/80">{f.source_name}</span>
          <ArrowRight className="h-3 w-3 shrink-0 text-primary/50" />
          <MiniHexIcon kind={f.target_kind} />
          <span className="max-w-[110px] truncate text-foreground/80">{f.target_name}</span>
        </div>
      </td>
      <td className="px-3 py-3 align-middle">
        <span className="font-mono text-[11px] text-muted-foreground">{f.category}</span>
      </td>
      <td className="px-3 py-3 align-middle">
        <div className="flex flex-wrap gap-1">
          {(f.owasp_map ?? []).map((tag) => (
            <span
              key={tag}
              className="rounded-[2px] border border-border bg-black/40 px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-[0.06em] text-muted-foreground"
            >
              {tag}
            </span>
          ))}
        </div>
      </td>
      <td className="px-3 py-3 align-middle">
        <div className="flex items-center justify-end gap-2">
          <MeterBar value={conf} max={100} color={color} height={4} className="w-14" />
          <span className="w-9 text-right font-mono text-[11px] font-semibold tabular-nums text-foreground">
            {conf}%
          </span>
        </div>
      </td>
    </tr>
  );
}
