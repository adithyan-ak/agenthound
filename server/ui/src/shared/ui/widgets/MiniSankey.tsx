import { NODE_KIND_COLORS, ACCENT, INSTRUMENT } from "@shared/theme/tokens";

/** One agent -> shared host -> resource cross-protocol pivot path. */
export interface PivotPath {
  agent: string;
  host: string;
  resource: string;
}

function truncate(s: string, n = 15): string {
  return s.length > n ? s.slice(0, n - 1) + "\u2026" : s;
}

/**
 * Three-lane Sankey-style diagram (A2A agent -> shared host -> MCP resource)
 * drawn in pure SVG. Promoted from the dashboard CrossProtocol widget so it
 * can be reused across surfaces.
 */
export function MiniSankey({ pivots }: { pivots: PivotPath[] }) {
  const agents = [...new Set(pivots.map((p) => p.agent))];
  const hosts = [...new Set(pivots.map((p) => p.host))];
  const resources = [...new Set(pivots.map((p) => p.resource))];

  const colW = 108;
  const svgW = 380;
  const gap = (svgW - colW * 3) / 2;
  const nodeH = 24;
  const nodeGap = 6;

  const yOf = (items: string[]) => items.map((_, i) => i * (nodeH + nodeGap) + 8);
  const agentYs = yOf(agents);
  const hostYs = yOf(hosts);
  const resourceYs = yOf(resources);

  const totalH =
    Math.max(
      agents.length * (nodeH + nodeGap),
      hosts.length * (nodeH + nodeGap),
      resources.length * (nodeH + nodeGap),
    ) + 16;

  const col1X = 0;
  const col2X = colW + gap;
  const col3X = colW * 2 + gap * 2;

  return (
    <svg width="100%" viewBox={`0 0 ${svgW} ${totalH}`} className="overflow-visible">
      <defs>
        <linearGradient id="xproto-link" x1="0" y1="0" x2="1" y2="0">
          <stop offset="0%" stopColor={ACCENT} stopOpacity={0.5} />
          <stop offset="100%" stopColor={NODE_KIND_COLORS.MCPResource} stopOpacity={0.6} />
        </linearGradient>
      </defs>

      {pivots.map((p, i) => {
        const ay = (agentYs[agents.indexOf(p.agent)] ?? 0) + nodeH / 2;
        const hy = (hostYs[hosts.indexOf(p.host)] ?? 0) + nodeH / 2;
        const ry = (resourceYs[resources.indexOf(p.resource)] ?? 0) + nodeH / 2;
        const x1 = col1X + colW;
        const x2 = col2X;
        const x3 = col2X + colW;
        const x4 = col3X;
        return (
          <g key={i}>
            <path
              d={`M${x1},${ay} C${x1 + gap / 2},${ay} ${x2 - gap / 2},${hy} ${x2},${hy}`}
              fill="none"
              stroke="url(#xproto-link)"
              strokeWidth={1.5}
            />
            <path
              d={`M${x3},${hy} C${x3 + gap / 2},${hy} ${x4 - gap / 2},${ry} ${x4},${ry}`}
              fill="none"
              stroke="url(#xproto-link)"
              strokeWidth={1.5}
            />
          </g>
        );
      })}

      {agents.map((name, i) => (
        <g key={`a-${name}`}>
          <rect x={col1X} y={agentYs[i]} width={colW} height={nodeH} rx={2} fill={INSTRUMENT.panel} stroke={NODE_KIND_COLORS.A2AAgent} strokeWidth={1} />
          <rect x={col1X} y={agentYs[i]} width={3} height={nodeH} fill={NODE_KIND_COLORS.A2AAgent} />
          <text x={col1X + 10} y={(agentYs[i] ?? 0) + 16} fill={INSTRUMENT.sankeyAgent} fontSize={9} className="select-none font-mono">
            {truncate(name)}
          </text>
        </g>
      ))}
      {hosts.map((name, i) => (
        <g key={`h-${name}`}>
          <rect x={col2X} y={hostYs[i]} width={colW} height={nodeH} rx={2} fill={INSTRUMENT.panel} stroke={INSTRUMENT.grayDim} strokeWidth={1} />
          <rect x={col2X} y={hostYs[i]} width={3} height={nodeH} fill={INSTRUMENT.grayDim} />
          <text x={col2X + 10} y={(hostYs[i] ?? 0) + 16} fill={INSTRUMENT.sankeyHost} fontSize={9} className="select-none font-mono">
            {truncate(name)}
          </text>
        </g>
      ))}
      {resources.map((name, i) => (
        <g key={`r-${name}`}>
          <rect x={col3X} y={resourceYs[i]} width={colW} height={nodeH} rx={2} fill={INSTRUMENT.panel} stroke={NODE_KIND_COLORS.MCPResource} strokeWidth={1} />
          <rect x={col3X} y={resourceYs[i]} width={3} height={nodeH} fill={NODE_KIND_COLORS.MCPResource} />
          <text x={col3X + 10} y={(resourceYs[i] ?? 0) + 16} fill={INSTRUMENT.sankeyResource} fontSize={9} className="select-none font-mono">
            {truncate(name)}
          </text>
        </g>
      ))}
    </svg>
  );
}
