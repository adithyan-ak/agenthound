import { useState } from "react";
import { Copy, Check } from "lucide-react";

interface CopyableCodeBlockProps {
  code: string;
}

export function CopyableCodeBlock({ code }: CopyableCodeBlockProps) {
  const [copied, setCopied] = useState(false);

  function handleCopy() {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="relative mt-2 rounded-[3px] border border-border bg-black/50 p-2.5">
      <button
        onClick={handleCopy}
        className="absolute right-1.5 top-1.5 rounded-[2px] p-1 text-muted-foreground transition-colors hover:bg-white/[0.06] hover:text-foreground"
        title="Copy to clipboard"
      >
        {copied ? <Check className="h-3 w-3 text-emerald-400" /> : <Copy className="h-3 w-3" />}
      </button>
      <pre className="whitespace-pre-wrap pr-6 font-mono text-[10px] leading-relaxed text-foreground/90">
        {code}
      </pre>
    </div>
  );
}
