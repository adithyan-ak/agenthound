import { Wrench } from "lucide-react";
import { WidgetCard } from "@shared/ui/widgets";
import type { RemediationStep } from "@entities/finding/model";
import { SIGNAL_OK } from "@shared/theme/tokens";
import { CopyableCodeBlock } from "./CopyableCodeBlock";

interface FindingRemediationProps {
  steps: RemediationStep[];
}

export function FindingRemediation({ steps }: FindingRemediationProps) {
  if (steps.length === 0) return null;

  return (
    <WidgetCard title="Remediation" icon={Wrench} accent={SIGNAL_OK}>
      <div className="space-y-3.5">
        {steps.map((step) => (
          <div key={step.step} className="flex gap-3">
            <span className="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-[3px] bg-primary/10 font-mono text-[11px] font-bold tabular-nums text-primary ring-1 ring-inset ring-primary/30">
              {String(step.step).padStart(2, "0")}
            </span>
            <div className="min-w-0 flex-1">
              <div className="text-sm font-semibold text-foreground">{step.title}</div>
              <p className="mt-0.5 text-xs leading-relaxed text-muted-foreground">
                {step.description}
              </p>
              {step.commands && step.commands.length > 0 && (
                <CopyableCodeBlock code={step.commands.join("\n")} />
              )}
            </div>
          </div>
        ))}
      </div>
    </WidgetCard>
  );
}
