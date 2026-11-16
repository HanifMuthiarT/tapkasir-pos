import type { Metadata } from "next";
import type { ReactNode } from "react";
import { Nunito } from "next/font/google";
import "./globals.css";

const bodyFont = Nunito({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-body",
});

export const metadata: Metadata = {
  title: "TapKasir Pos",
  description: "Pilih menu favorit dan pesan langsung dari meja melalui QR Code.",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="id">
      <body className={bodyFont.variable}>
        {children}
      </body>
    </html>
  );
}
