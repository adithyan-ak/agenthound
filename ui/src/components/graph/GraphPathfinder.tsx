import { useEffect, useMemo, useRef, useState } from "react";
import {
  Route,
  ChevronUp,
  ChevronDown,
  Search,
  ArrowUpDown,
  X,
  Loader2,
  ArrowRight,
} from "lucide-react";
import { useReactFlow, type Node } from "@xyflow/react";
import {
  useShortestPath,
  useAllPaths,
  useWeightedPath,
} from "@/hooks/usePathfinding";
import { useGraphStore } from "@/store/graph";
import { NODE_COLORS } from "@/lib/node-styles";
import type { NodeKind, Path, PathRequest } from "@/api/types";
import { cn } from "@/lib/utils";

const NODE_KINDS: NodeKind[] = [
  "AgentInstance",
  "MCPServer",
  "MCPTool",
  "MCPResource",
  "MCPPrompt",
  "A2AAgent",
  "A2ASkill",
  "Identity",
  "Credential",
  "Host",
  "ConfigFile",
  "InstructionFile",
];

type Algorithm = "shortest" | "all" | "weighted";

interface NameSuggestion {
  id: string;
  name: string;
  kind: string;
}

interface NameAutocompleteProps {
  value: string;
  onChange: (v: string) => void;
  onPick: (name: string) => void;
  placeholder: string;
  filterKind?: string;
}

function NameAutocomplete({
  value,
  onChange,
  onPick,
  placeholder,
  filterKind,
}: NameAutocompleteProps) {
  const reactFlow = useReactFlow();
  const [open, setOpen] = useState(false);
  const [activeIdx, setActiveIdx] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);

  const suggestions = useMemo<NameSuggestion[]>(() => {
    if (value.length < 1) return [];
    const query = value.toLowerCase();
    const all = reactFlow.getNodes() as Node[];
    const out: NameSuggestion[] = [];
    for (const n of all) {
      const d = n.data as Record<string, unknown>;
      const label = String(d?.label ?? "");
      const kind = String(d?.kind ?? "");
      if (filterKind && kind !== filterKind) continue;
      if (!label.toLowerCase().includes(query)) continue;
      out.push({ id: n.id, name: label, kind });
      if (out.length >= 8) break;
    }
    return out;
  }, [value, filterKind, reactFlow]);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as globalThis.Node)
      ) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  function handleKey(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIdx((i) => Math.min(i + 1, suggestions.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIdx((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter" && suggestions[activeIdx]) {
      e.preventDefault();
      onPick(suggestions[activeIdx].name);
      setOpen(false);
    }
  }

  return (
    <div className="relative" ref={containerRef}>
      <input
        type="text"
        value={value}
        onChange={(e) => {
          onChange(e.target.value);
          setOpen(true);
          setActiveIdx(0);
        }}
        onFocus={() => setOpen(true)}
        onKeyDown={handleKey}
        placeholder={placeholder}
        className="w-full h-8 px-2 rounded bg-background border border-border text-xs text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/60"
      />
      {open && suggestions.length > 0 && (
        <div className="absolute z-30 mt-1 w-full rounded border border-border bg-card shadow-lg max-h-[200px] overflow-y-auto">
          {suggestions.map((s, i) => (
            <button
              key={s.id}
              type="button"
              onMouseEnter={() => setActiveIdx(i)}
              onClick={() => {
                onPick(s.name);
                setOpen(false);
              }}
              className={cn(
                "w-full text-left flex items-center gap-2 px-2 py-1.5 text-xs hover:bg-accent",
                i === activeIdx && "bg-accent",
              )}
            >
              <span
                className="h-1.5 w-1.5 rounded-full flex-shrink-0"
                style={{ backgroundColor: NODE_COLORS[s.kind] ?? "#666" }}
              />
              <span className="text-foreground truncate flex-1">{s.name}</span>
              <span className="text-[9px] text-muted-foreground uppercase">
                {s.kind}
              </span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

export function GraphPathfinder() {
  const [expanded, setExpanded] = useState(true);
  const [sourceKind, setSourceKind] = useState<NodeKind>("AgentInstance");
  const [sourceName, setSourceName] = useState("");
  const [targetKind, setTargetKind] = useState<NodeKind | "">("");
  const [targetName, setTargetName] = useState("");
  const [algorithm, setAlgorithm] = useState<Algorithm>("shortest");
  const [maxHops, setMaxHops] = useState(6);
  const [results, setResults] = useState<Path[] | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const highlightPath = useGraphStore((s) => s.highlightPath);
  const clearHighlight = useGraphStore((s) => s.clearHighlight);

  const shortest = useShortestPath();
  const all = useAllPaths();
  const weighted = useWeightedPath();

  const active =
    algorithm === "shortest"
      ? shortest
      : algorithm === "all"
        ? all
        : weighted;

  const isPending = shortest.isPending || all.isPending || weighted.isPending;

  function swap() {
    const sk = sourceKind;
    const sn = sourceName;
    setSourceKind((targetKind || "AgentInstance") as NodeKind);
    setSourceName(targetName);
    setTargetKind(sk);
    setTargetName(sn);
  }

  function applyPathHighlight(path: Path) {
    const nodeIds = path.nodes.map((n) => n.id);
    const edgeKeys = path.edges.map(
      (e) => `${e.source}->${e.target}:${e.kind}`,
    );
    highlightPath({ nodeIds, edgeKeys, title: "Pathfinder result" });
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErrorMsg(null);
    if (!sourceName.trim()) {
      setErrorMsg("Source name is required");
      return;
    }

    const req: PathRequest = {
      source: sourceName.trim(),
      target: targetName.trim(),
      source_kind: sourceKind,
      ...(targetKind && { target_kind: targetKind }),
      max_hops: maxHops,
      limit: 20,
    };

    active.mutate(req, {
      onSuccess: (response) => {
        setResults(response.paths);
        if (response.paths.length > 0) {
          applyPathHighlight(response.paths[0]!);
        } else {
          setErrorMsg("No paths found");
          clearHighlight();
        }
      },
      onError: (err) => {
        setErrorMsg(err instanceof Error ? err.message : "Request failed");
      },
    });
  }

  function handleReset() {
    setSourceName("");
    setTargetName("");
    setTargetKind("");
    setResults(null);
    setErrorMsg(null);
    clearHighlight();
  }

  return (
    <div className="absolute bottom-4 left-4 z-10 w-[340px]">
      <div className="rounded-lg border border-border bg-card/95 shadow-xl backdrop-blur overflow-hidden">
        {/* Header */}
        <button
          type="button"
          onClick={() => setExpanded(!expanded)}
          className="w-full flex items-center justify-between px-3 py-2 border-b border-border hover:bg-accent/50 transition-colors"
        >
          <div className="flex items-center gap-2">
            <Route className="h-4 w-4 text-primary" />
            <span className="text-sm font-semibold text-foreground">
              Pathfinder
            </span>
            {results && results.length > 0 && (
              <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-primary/20 text-primary">
                {results.length}
              </span>
            )}
          </div>
          {expanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          ) : (
            <ChevronUp className="h-4 w-4 text-muted-foreground" />
          )}
        </button>

        {expanded && (
          <div className="max-h-[65vh] overflow-y-auto">
            <form onSubmit={handleSubmit} className="p-3 space-y-2.5">
              {/* Source */}
              <div className="space-y-1">
                <label className="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                  Source
                </label>
                <div className="flex gap-1.5">
                  <select
                    value={sourceKind}
                    onChange={(e) =>
                      setSourceKind(e.target.value as NodeKind)
                    }
                    className="h-8 w-[130px] flex-shrink-0 px-1.5 rounded bg-background border border-border text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/60"
                  >
                    {NODE_KINDS.map((k) => (
                      <option key={k} value={k}>
                        {k}
                      </option>
                    ))}
                  </select>
                  <div className="flex-1">
                    <NameAutocomplete
                      value={sourceName}
                      onChange={setSourceName}
                      onPick={setSourceName}
                      placeholder="Source name..."
                      filterKind={sourceKind}
                    />
                  </div>
                </div>
              </div>

              {/* Swap */}
              <div className="flex justify-center">
                <button
                  type="button"
                  onClick={swap}
                  title="Swap source and target"
                  className="h-6 w-6 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
                >
                  <ArrowUpDown className="h-3.5 w-3.5" />
                </button>
              </div>

              {/* Target */}
              <div className="space-y-1">
                <label className="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                  Target <span className="normal-case opacity-70">(optional)</span>
                </label>
                <div className="flex gap-1.5">
                  <select
                    value={targetKind || "__any__"}
                    onChange={(e) =>
                      setTargetKind(
                        e.target.value === "__any__"
                          ? ""
                          : (e.target.value as NodeKind),
                      )
                    }
                    className="h-8 w-[130px] flex-shrink-0 px-1.5 rounded bg-background border border-border text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/60"
                  >
                    <option value="__any__">Any kind</option>
                    {NODE_KINDS.map((k) => (
                      <option key={k} value={k}>
                        {k}
                      </option>
                    ))}
                  </select>
                  <div className="flex-1">
                    <NameAutocomplete
                      value={targetName}
                      onChange={setTargetName}
                      onPick={setTargetName}
                      placeholder="Target name (opt.)..."
                      filterKind={targetKind || undefined}
                    />
                  </div>
                </div>
              </div>

              {/* Algorithm */}
              <div className="space-y-1">
                <label className="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                  Algorithm
                </label>
                <div className="grid grid-cols-3 gap-1 p-0.5 rounded border border-border bg-background">
                  {(["shortest", "all", "weighted"] as const).map((alg) => (
                    <button
                      key={alg}
                      type="button"
                      onClick={() => setAlgorithm(alg)}
                      className={cn(
                        "text-[11px] py-1 rounded capitalize transition-colors",
                        algorithm === alg
                          ? "bg-primary text-primary-foreground font-medium"
                          : "text-muted-foreground hover:text-foreground",
                      )}
                    >
                      {alg}
                    </button>
                  ))}
                </div>
              </div>

              {/* Max Hops */}
              <div className="space-y-1">
                <div className="flex items-center justify-between">
                  <label className="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                    Max Hops
                  </label>
                  <span className="text-[11px] font-mono text-foreground">
                    {maxHops}
                  </span>
                </div>
                <input
                  type="range"
                  min={1}
                  max={20}
                  value={maxHops}
                  onChange={(e) => setMaxHops(Number(e.target.value))}
                  className="w-full accent-primary"
                />
              </div>

              {/* Error */}
              {errorMsg && (
                <div className="rounded bg-red-950/50 border border-red-900 px-2 py-1.5 text-[11px] text-red-300">
                  {errorMsg}
                </div>
              )}

              {/* Buttons */}
              <div className="flex gap-1.5 pt-1">
                <button
                  type="submit"
                  disabled={isPending || !sourceName.trim()}
                  className="flex-1 h-8 rounded bg-primary text-primary-foreground text-xs font-medium flex items-center justify-center gap-1.5 hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {isPending ? (
                    <Loader2 className="h-3.5 w-3.5 animate-spin" />
                  ) : (
                    <Search className="h-3.5 w-3.5" />
                  )}
                  Find Paths
                </button>
                {(results || sourceName) && (
                  <button
                    type="button"
                    onClick={handleReset}
                    title="Reset"
                    className="h-8 w-8 rounded border border-border text-muted-foreground hover:text-foreground hover:bg-accent flex items-center justify-center"
                  >
                    <X className="h-3.5 w-3.5" />
                  </button>
                )}
              </div>
            </form>

            {/* Results */}
            {results && results.length > 0 && (
              <div className="border-t border-border p-3 space-y-1.5">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-[10px] font-semibold uppercase tracking-wide text-muted-foreground">
                    Results ({results.length})
                  </span>
                  <span className="text-[9px] text-muted-foreground">
                    click to highlight
                  </span>
                </div>
                {results.map((path, i) => (
                  <button
                    key={i}
                    type="button"
                    onClick={() => applyPathHighlight(path)}
                    className="w-full rounded border border-border bg-background/40 p-2 text-left hover:bg-accent hover:border-primary/40 transition-colors"
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-[10px] font-mono text-foreground">
                        #{i + 1}
                      </span>
                      <span className="text-[10px] text-muted-foreground">
                        {path.hops} hop{path.hops !== 1 ? "s" : ""}
                      </span>
                      {path.weight != null && (
                        <span className="text-[10px] text-muted-foreground">
                          w={path.weight.toFixed(2)}
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-1 text-[10px] text-foreground">
                      <span
                        className="h-1.5 w-1.5 rounded-full flex-shrink-0"
                        style={{
                          backgroundColor:
                            NODE_COLORS[path.nodes[0]?.kinds[0] ?? ""] ??
                            "#666",
                        }}
                      />
                      <span className="truncate max-w-[90px]">
                        {path.nodes[0]?.name ?? "?"}
                      </span>
                      <ArrowRight className="h-2.5 w-2.5 text-muted-foreground flex-shrink-0" />
                      <span className="text-muted-foreground">
                        {path.nodes.length - 2 > 0
                          ? `+${path.nodes.length - 2}`
                          : ""}
                      </span>
                      {path.nodes.length - 2 > 0 && (
                        <ArrowRight className="h-2.5 w-2.5 text-muted-foreground flex-shrink-0" />
                      )}
                      <span
                        className="h-1.5 w-1.5 rounded-full flex-shrink-0"
                        style={{
                          backgroundColor:
                            NODE_COLORS[
                              path.nodes[path.nodes.length - 1]?.kinds[0] ?? ""
                            ] ?? "#666",
                        }}
                      />
                      <span className="truncate max-w-[90px]">
                        {path.nodes[path.nodes.length - 1]?.name ?? "?"}
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            )}

            {results && results.length === 0 && !errorMsg && (
              <div className="border-t border-border p-3 text-center text-[11px] text-muted-foreground">
                No paths found
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
