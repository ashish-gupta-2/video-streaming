import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Video Streaming App',
  description: 'Stream videos using HLS',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}

