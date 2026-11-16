
import { Suspense } from "react";
import CustomerPage from "./components/CustomerPage";

export default function Home() {
  return (
    <Suspense
      fallback={
        <main style={{ padding: 24 }}>
          Memuat menu...
        </main>
      }
    >
      <CustomerPage />
    </Suspense>
  );
}
