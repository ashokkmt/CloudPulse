import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'CloudPulse',
  description: 'Minimal CloudPulse demo app for the DigitalOcean bootcamp roadmap.',
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
