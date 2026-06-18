import { Stack, Switcher } from "@/components/ui/layout";
import { StatCards } from "./StatCards";
import { RiskChart } from "./RiskChart";
import { AuthCoverage } from "./AuthCoverage";
import { TopRiskyEntities } from "./TopRiskyEntities";
import { CredentialExposure } from "./CredentialExposure";
import { CrossProtocol } from "./CrossProtocol";
import { TopFindings } from "./TopFindings";
import { RecentScans } from "./RecentScans";

export function Dashboard() {
  return (
    <Stack className="p-6" gap="1.5rem">
      <h2 className="text-lg font-semibold text-foreground">Dashboard</h2>
      <StatCards />
      <Switcher threshold="48rem">
        <RiskChart />
        <TopRiskyEntities />
      </Switcher>
      <Switcher threshold="48rem">
        <AuthCoverage />
        <CredentialExposure />
      </Switcher>
      <Switcher threshold="48rem">
        <CrossProtocol />
        <TopFindings />
      </Switcher>
      <RecentScans />
    </Stack>
  );
}
