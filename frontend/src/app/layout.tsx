"use client";

import { ConfigProvider } from "antd";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "./globals.css";

const queryClient = new QueryClient();

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ru">
      <body>
        <QueryClientProvider client={queryClient}>
          <ConfigProvider>
            {children}
          </ConfigProvider>
        </QueryClientProvider>
      </body>
    </html>
  );
}
