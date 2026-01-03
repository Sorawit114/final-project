import Link from "next/link";
import { ReactNode, useEffect, useState } from "react";
import { useRouter } from "next/router";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const router = useRouter();
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("username");
    if (storedUser) {
      setUsername(storedUser);
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("userId");
    localStorage.removeItem("username");
    setUsername(null);
    // Force reload to clear all states (especially in index.tsx)
    window.location.href = "/";
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-900 via-purple-900 to-gray-900 relative overflow-hidden text-white font-sans">
      {/* Background Decor */}
      <div className="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
        <div className="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-purple-600/20 rounded-full blur-[100px]" />
        <div className="absolute bottom-[10%] right-[10%] w-[30%] h-[30%] bg-indigo-600/20 rounded-full blur-[100px]" />
      </div>

      {/* Header */}
      <header className="relative z-10 bg-white/5 backdrop-blur-md border-b border-white/10 px-6 py-4 flex justify-between items-center shadow-lg">
        <Link href="/" className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-indigo-400 hover:opacity-80 transition-opacity">
          Stock
        </Link>

        <div className="flex gap-6 items-center">
          <nav className="flex gap-6">
            <Link href="/" className="text-gray-300 hover:text-white transition-colors">Home</Link>
            <Link href="/my-stock" className="text-gray-300 hover:text-white transition-colors">My Stock</Link>
            <Link href="/test-model" className="text-gray-300 hover:text-white transition-colors">Test Model</Link>
          </nav>

          {/* Auth State */}
          {username ? (
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className="w-8 h-8 rounded-full bg-gradient-to-r from-purple-500 to-indigo-500 flex items-center justify-center text-sm font-bold shadow-md">
                  {username.charAt(0).toUpperCase()}
                </div>
                <span className="font-medium text-purple-200">{username}</span>
              </div>
              <button
                onClick={handleLogout}
                className="text-sm bg-white/10 hover:bg-white/20 text-white px-3 py-1.5 rounded-lg transition-all border border-white/10"
              >
                Logout
              </button>
            </div>
          ) : (
            <Link href="/login">
              <button className="bg-indigo-600 hover:bg-indigo-500 px-5 py-2 rounded-lg shadow-lg shadow-indigo-500/20 transition-all hover:-translate-y-0.5 font-medium">
                Login
              </button>
            </Link>
          )}
        </div>
      </header>

      {/* Content */}
      <main className="relative z-10 container mx-auto px-4 py-8">
        {children}
      </main>
    </div>
  );
}
