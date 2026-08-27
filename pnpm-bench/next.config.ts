import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Keep build work realistic without needing external services.
  output: "standalone",
};

export default nextConfig;
