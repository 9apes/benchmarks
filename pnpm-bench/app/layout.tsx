import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "cache-bench pnpm",
  description: "Synthetic Next.js app for CI cache benchmarking",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
