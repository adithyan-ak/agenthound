import { createRoot } from "react-dom/client";
import { AppProviders } from "./providers/AppProviders";
import { AppRoutes } from "./routes";
import "@shared/styles/globals.css";

createRoot(document.getElementById("root")!).render(
  <AppProviders>
    <AppRoutes />
  </AppProviders>,
);
