import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";

const poppins = Poppins({
  variable: "--font-poppins",
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "GNS - (Go, NextJS, ShadCn)",
  description: "Base project build with Go, NextJS, ShadCn",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en"
      className="dark"
      style={{ colorScheme: "dark" }}
    >
      <body
        className={`${poppins.variable} antialiased`}
      >
        <Providers>
          {children}
        </Providers>
      </body>
    </html >
  );
}
