import { X } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@shared/ui/primitives/select";

interface RuleFiltersProps {
  severity: string;
  collector: string;
  onSeverityChange: (value: string) => void;
  onCollectorChange: (value: string) => void;
  onClear: () => void;
}

const triggerClass =
  "h-8 w-[140px] rounded-[3px] border-border bg-black/30 font-mono text-[11px] uppercase tracking-[0.06em] text-foreground/80";

export function RuleFilters({
  severity,
  collector,
  onSeverityChange,
  onCollectorChange,
  onClear,
}: RuleFiltersProps) {
  const hasFilters = severity !== "all" || collector !== "all";

  return (
    <div className="flex shrink-0 items-center gap-2">
      <Select value={severity} onValueChange={onSeverityChange}>
        <SelectTrigger className={triggerClass}>
          <SelectValue placeholder="Severity" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Severities</SelectItem>
          <SelectItem value="critical">Critical</SelectItem>
          <SelectItem value="high">High</SelectItem>
          <SelectItem value="medium">Medium</SelectItem>
          <SelectItem value="low">Low</SelectItem>
          <SelectItem value="info">Info</SelectItem>
        </SelectContent>
      </Select>

      <Select value={collector} onValueChange={onCollectorChange}>
        <SelectTrigger className={triggerClass}>
          <SelectValue placeholder="Collector" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Collectors</SelectItem>
          <SelectItem value="mcp">MCP</SelectItem>
          <SelectItem value="a2a">A2A</SelectItem>
          <SelectItem value="config">Config</SelectItem>
        </SelectContent>
      </Select>

      {hasFilters && (
        <button
          onClick={onClear}
          className="inline-flex h-8 items-center gap-1 rounded-[3px] px-2 font-mono text-[11px] uppercase tracking-[0.08em] text-muted-foreground transition-colors hover:bg-white/[0.05] hover:text-foreground"
        >
          <X className="h-3.5 w-3.5" />
          Clear
        </button>
      )}
    </div>
  );
}
