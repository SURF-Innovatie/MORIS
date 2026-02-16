import React from "react";
import ReactDOM from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { RouterProvider } from "react-router-dom";

import "./index.css";
import { createAppRouter } from "./router";
import { AuthProvider } from "./contexts/AuthContext";
import { BackendStatusProvider } from "./contexts/BackendStatusContext";
import { Toaster } from "@/components/ui/sonner";
import { BackendOfflineAlert } from "@/components/status/BackendOfflineAlert";
import { SessionExpiredAlert } from "@/components/status/SessionExpiredAlert";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 1000 * 30,
    },
  },
});

const router = createAppRouter();

const rootElement = document.getElementById("root");
if (!rootElement) {
  throw new Error("Root container missing in index.html");
}

ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <BackendStatusProvider>
        <AuthProvider>
          <RouterProvider router={router} />
          <Toaster />
          <BackendOfflineAlert />
          <SessionExpiredAlert />
          <ReactQueryDevtools
            initialIsOpen={false}
            buttonPosition="bottom-right"
          />
        </AuthProvider>
      </BackendStatusProvider>
    </QueryClientProvider>
  </React.StrictMode>,
);
