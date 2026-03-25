import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { store } from "./store/store";
import "./theme/mac-os-1.css";
import "./index.css";
import App from "./App.tsx";

async function boot() {
  const useMsw = import.meta.env.DEV && import.meta.env.VITE_USE_MSW === "1";

  if (useMsw) {
    const { worker } = await import("./mocks/browser");
    await worker.start({ onUnhandledRequest: "bypass" });
  }

  createRoot(document.getElementById("root")!).render(
    <StrictMode>
      <Provider store={store}>
        <App />
      </Provider>
    </StrictMode>
  );
}

boot();
