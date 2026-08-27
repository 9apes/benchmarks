"use client";

import { useMemo } from "react";
import { format } from "date-fns";
import _ from "lodash";

export default function Home() {
  const label = useMemo(
    () => _.startCase("cache bench pnpm install and build"),
    [],
  );

  return (
    <main style={{ padding: 24, fontFamily: "system-ui, sans-serif" }}>
      <h1>{label}</h1>
      <p>Built at {format(new Date(), "yyyy-MM-dd HH:mm:ss")}</p>
    </main>
  );
}
