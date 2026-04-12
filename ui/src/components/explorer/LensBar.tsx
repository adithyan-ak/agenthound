import { useState } from "react";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { ChevronDown, Check } from "lucide-react";
import { LENS_LIST, type LensDefinition } from "@/lib/explorer/lens-config";
import { useExplorerStore, type LensId } from "@/store/explorer";
import { cn } from "@/lib/utils";

interface LensPillProps {
  lens: LensDefinition;
  active: boolean;
  onClick: () => void;
}

function LensPill({ lens, active, onClick }: LensPillProps) {
  const Icon = lens.icon;
  const subPresets = useExplorerStore((s) => s.subPresets[lens.id] ?? []);
  const toggleSubPreset = useExplorerStore((s) => s.toggleSubPreset);
  const [open, setOpen] = useState(false);

  const hasSubPresets = lens.subPresets.length > 0;

  const pillClasses = cn(
    "group relative flex items-center gap-1.5 rounded-full px-3 py-1.5 text-xs font-medium",
    "transition-all duration-150 ease-out border",
    "whitespace-nowrap select-none",
    active
      ? cn(lens.accentClass, "border-transparent shadow-lg")
      : "bg-slate-900/60 text-slate-300 border-slate-700/60 hover:bg-slate-800 hover:text-white hover:border-slate-600",
  );

  const pillStyle = active
    ? { boxShadow: `0 0 20px -4px ${lens.activeTint}70` }
    : undefined;

  if (!hasSubPresets) {
    return (
      <button className={pillClasses} style={pillStyle} onClick={onClick}>
        <Icon className="h-3.5 w-3.5" strokeWidth={2.25} />
        <span>{lens.label}</span>
      </button>
    );
  }

  return (
    <DropdownMenu.Root open={open} onOpenChange={setOpen}>
      <div className={pillClasses} style={pillStyle}>
        <button onClick={onClick} className="flex items-center gap-1.5">
          <Icon className="h-3.5 w-3.5" strokeWidth={2.25} />
          <span>{lens.label}</span>
        </button>
        <DropdownMenu.Trigger asChild>
          <button
            className={cn(
              "flex items-center justify-center rounded-full -mr-1 ml-0.5 h-4 w-4",
              "hover:bg-black/20 transition-colors",
            )}
            aria-label={`${lens.label} sub-presets`}
          >
            <ChevronDown className="h-3 w-3" strokeWidth={2.5} />
          </button>
        </DropdownMenu.Trigger>
      </div>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          align="start"
          sideOffset={8}
          className={cn(
            "min-w-[240px] rounded-lg border border-slate-700/80 bg-slate-900/98 p-1.5 shadow-2xl",
            "backdrop-blur",
            "data-[state=open]:animate-in data-[state=closed]:animate-out",
            "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
            "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
            "data-[side=bottom]:slide-in-from-top-2 z-50",
          )}
        >
          <div className="px-2 py-1.5 text-[10px] uppercase tracking-widest text-slate-500">
            {lens.label} · sub-presets
          </div>
          {lens.subPresets.map((sp) => {
            const enabled = subPresets.includes(sp.id);
            return (
              <DropdownMenu.Item
                key={sp.id}
                onSelect={(e) => {
                  e.preventDefault();
                  toggleSubPreset(lens.id, sp.id);
                }}
                className={cn(
                  "group/item flex cursor-pointer items-start gap-2 rounded-md px-2 py-2 text-xs",
                  "outline-none focus:bg-slate-800/80 data-[highlighted]:bg-slate-800/80",
                  "transition-colors",
                )}
              >
                <div
                  className={cn(
                    "mt-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-sm border",
                    enabled
                      ? "bg-blue-500 border-blue-400"
                      : "border-slate-600 bg-transparent",
                  )}
                >
                  {enabled && (
                    <Check className="h-2.5 w-2.5 text-white" strokeWidth={3} />
                  )}
                </div>
                <div className="flex flex-col">
                  <span className="font-medium text-white">{sp.label}</span>
                  <span className="text-[10px] text-slate-400 leading-tight">
                    {sp.description}
                  </span>
                </div>
              </DropdownMenu.Item>
            );
          })}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}

export function LensBar() {
  const activeLens = useExplorerStore((s) => s.activeLens);
  const setActiveLens = useExplorerStore((s) => s.setActiveLens);

  return (
    <div className="pointer-events-auto absolute left-1/2 top-4 z-30 -translate-x-1/2">
      <div
        className={cn(
          "flex items-center gap-1.5 rounded-full border border-slate-800/80 bg-slate-950/90 px-2 py-1.5",
          "shadow-2xl backdrop-blur-md",
        )}
      >
        <div className="px-2 text-[10px] font-semibold uppercase tracking-widest text-slate-500">
          Lens
        </div>
        {LENS_LIST.map((lens) => (
          <LensPill
            key={lens.id}
            lens={lens}
            active={activeLens === lens.id}
            onClick={() => setActiveLens(lens.id as LensId)}
          />
        ))}
      </div>
    </div>
  );
}
