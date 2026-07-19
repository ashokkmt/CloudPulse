import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "CloudPulse Dashboard",
  description: "Real-time analytics and task management",
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
