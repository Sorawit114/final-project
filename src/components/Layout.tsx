import Link from "next/link";
import { ReactNode } from "react";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Header */}
      <header className="bg-gray-800 p-4 flex justify-between items-center shadow-md justify-end">
        <div className="flex gap-4 items-center">
          <nav className="flex gap-4">
            <Link href="/" className="hover:text-indigo-400">Home</Link>
            <Link href="/my-stock" className="hover:text-indigo-400">My Stock</Link>
            <Link href="/test-model" className="hover:text-indigo-400">Test Model</Link>
          </nav>
          {/* Login/Logout Button */}
          <Link href="/login">
            <button className="ml-4 bg-indigo-600 hover:bg-indigo-500 px-4 py-2 rounded-md shadow">
              Login
            </button>
          </Link>
        </div>
      </header>

      {/* Content */}
      <main>{children}</main>
    </div>
  );
}
