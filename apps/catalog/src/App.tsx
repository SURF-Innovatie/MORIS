import { QueryClient } from "@tanstack/react-query";
import { PersistQueryClientProvider } from "@tanstack/react-query-persist-client";
import { createSyncStoragePersister } from "@tanstack/query-sync-storage-persister";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { CatalogProvider } from "./components/CatalogProvider";
import { Layout } from "./components/Layout";
import { Home } from "./pages/Home";
import { Outputs } from "./pages/Outputs";
import { ProjectDetail } from "./pages/ProjectDetail";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24, // 24 hours
    },
  },
});

const persister = createSyncStoragePersister({
  storage: window.localStorage,
});

function App() {
  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister }}
    >
      <CatalogProvider>
        <BrowserRouter>
          <Layout>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/outputs" element={<Outputs />} />
              <Route path="/projects/:id" element={<ProjectDetail />} />
            </Routes>
          </Layout>
        </BrowserRouter>
      </CatalogProvider>
    </PersistQueryClientProvider>
  );
}

export default App;
