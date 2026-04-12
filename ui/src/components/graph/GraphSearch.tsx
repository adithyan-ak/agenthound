import { useEffect, useMemo, useRef, useState } from "react";
import { Search, X } from "lucide-react";
import type { Node } from "@xyflow/react";
import { cn } from "@/lib/utils";
import { NODE_COLORS } from "@/lib/node-styles";

interface GraphSearchProps {
  nodes: Node[];
  onSelect: (nodeId: string) => void;
}

interface Match {
  id: string;
  label: string;
  kind: string;
  score: number;
}

function scoreMatch(label: string, query: string): number {
  const l = label.toLowerCase();
  const q = query.toLowerCase();
  if (l === q) return 1000;
  if (l.startsWith(q)) return 500 - l.length;
  const idx = l.indexOf(q);
  if (idx >= 0) return 200 - idx - l.length * 0.1;
  return -1;
}

export function GraphSearch({ nodes, onSelect }: GraphSearchProps) {
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (e.key === "/" && document.activeElement?.tagName !== "INPUT") {
        e.preventDefault();
        inputRef.current?.focus();
        setOpen(true);
      }
      if (e.key === "Escape") {
        setOpen(false);
        inputRef.current?.blur();
      }
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, []);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as globalThis.Node)
      ) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const matches = useMemo<Match[]>(() => {
    if (query.length < 1) return [];
    const results: Match[] = [];
    for (const n of nodes) {
      const d = n.data as Record<string, unknown>;
      const label = String(d?.label ?? n.id);
      const kind = String(d?.kind ?? "Unknown");
      const score = scoreMatch(label, query);
      if (score >= 0) {
        results.push({ id: n.id, label, kind, score });
      }
    }
    results.sort((a, b) => b.score - a.score);
    return results.slice(0, 15);
  }, [nodes, query]);

  useEffect(() => {
    setActiveIndex(0);
  }, [matches]);

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => Math.min(i + 1, matches.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter" && matches[activeIndex]) {
      e.preventDefault();
      onSelect(matches[activeIndex].id);
      setQuery("");
      setOpen(false);
      inputRef.current?.blur();
    }
  }

  return (
    <div
      ref={containerRef}
      className="absolute top-4 left-4 z-20 w-[340px]"
    >
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
        <input
          ref={inputRef}
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={handleKeyDown}
          placeholder="Search nodes...    /"
          className="w-full h-9 pl-9 pr-8 rounded-md bg-card/95 border border-border text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40 backdrop-blur shadow-md"
        />
        {query && (
          <button
            onClick={() => {
              setQuery("");
              inputRef.current?.focus();
            }}
            className="absolute right-2 top-1/2 -translate-y-1/2 h-5 w-5 flex items-center justify-center rounded text-muted-foreground hover:text-foreground"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
      {open && matches.length > 0 && (
        <div className="mt-1 rounded-md bg-card/95 border border-border shadow-lg backdrop-blur overflow-hidden">
          <ul className="max-h-[340px] overflow-y-auto">
            {matches.map((m, i) => (
              <li key={m.id}>
                <button
                  onMouseEnter={() => setActiveIndex(i)}
                  onClick={() => {
                    onSelect(m.id);
                    setQuery("");
                    setOpen(false);
                    inputRef.current?.blur();
                  }}
                  className={cn(
                    "w-full text-left flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent",
                    i === activeIndex && "bg-accent",
                  )}
                >
                  <span
                    className="h-2 w-2 rounded-full flex-shrink-0"
                    style={{
                      backgroundColor: NODE_COLORS[m.kind] ?? "#666",
                    }}
                  />
                  <span className="text-foreground truncate flex-1">
                    {m.label}
                  </span>
                  <span className="text-[10px] text-muted-foreground uppercase tracking-wide">
                    {m.kind}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
      {open && query.length >= 1 && matches.length === 0 && (
        <div className="mt-1 rounded-md bg-card/95 border border-border shadow-lg backdrop-blur px-3 py-2 text-xs text-muted-foreground">
          No nodes match "{query}"
        </div>
      )}
    </div>
  );
}
