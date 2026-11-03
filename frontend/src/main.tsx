import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App";

import "./styles/variables.css";
import "./styles/reset.css";
import "./styles/global.css";
import "./styles/layout.css";
import "./styles/components.css";
import "./styles/dashboard.css";
import "./styles/auth.css";
import "./styles/responsive.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);