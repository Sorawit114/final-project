import { useEffect, useState } from "react";
import Link from "next/link";
import Layout from "../components/Layout";

interface Stock {
  symbol: string;
  quantity: number;
  avgPrice: number;
}

interface StockWithPrice extends Stock {
  currentPrice: number;
}

export default function MyStock() {
  const [stocks, setStocks] = useState<StockWithPrice[]>([]);
  const [loading, setLoading] = useState(true);
  const [userId, setUserId] = useState<string | null>(null);

  useEffect(() => {
    // Check for user
    const storedUserId = localStorage.getItem("userId");
    if (storedUserId) {
      setUserId(storedUserId);
      fetchStocks(storedUserId);
    } else {
      setLoading(false);
    }
  }, []);

  const fetchStocks = async (uid: string) => {
    try {
      const res = await fetch(`http://localhost:8080/user-stocks?userId=${uid}`);
      const data = await res.json();

      if (res.ok && data.stocks) {
        // Map backend data to frontend model
        // Backend only gives name, so we mock price/qty for now or default to 0
        const mappedStocks: StockWithPrice[] = data.stocks.map((s: any) => ({
          symbol: s.stockShortName,
          quantity: 100, // Default mock
          avgPrice: 150.00, // Default mock
          currentPrice: 155.00 + (Math.random() * 10 - 5), // Mock fluctuation
        }));
        setStocks(mappedStocks);
      }
    } catch (err) {
      console.error("Failed to fetch stocks", err);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdatePrice = (symbol: string) => {
    // Mock update price
    setStocks((prev) =>
      prev.map((s) =>
        s.symbol === symbol
          ? { ...s, currentPrice: s.currentPrice * (1 + (Math.random() * 0.05 - 0.025)) }
          : s
      )
    );
  };

  return (
    <Layout>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-white bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-indigo-400">
          My Stock Portfolio
        </h1>
      </div>

      {!userId ? (
        <div className="text-center py-20 bg-white/5 rounded-xl border border-white/10 backdrop-blur-md">
          <h2 className="text-xl text-gray-300 mb-4">Please login to view your portfolio</h2>
          <Link href="/login">
            <button className="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-2 rounded-full font-bold transition-all shadow-lg hover:-translate-y-1">
              Login Now
            </button>
          </Link>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl bg-white/5 backdrop-blur-md border border-white/10 shadow-xl">
          {loading ? (
            <div className="p-10 text-center text-gray-400">Loading your entries...</div>
          ) : stocks.length === 0 ? (
            <div className="p-10 text-center text-gray-400">
              <p className="mb-4">You haven't added any stocks yet.</p>
              <Link href="/">
                <button className="text-indigo-400 hover:text-indigo-300 font-semibold hover:underline">
                  Browse Stocks to Add
                </button>
              </Link>
            </div>
          ) : (
            <table className="w-full table-auto text-left border-collapse">
              <thead>
                <tr className="bg-white/5 border-b border-white/10">
                  <th className="p-4 font-semibold text-gray-200">Symbol</th>
                  <th className="p-4 font-semibold text-gray-200">Quantity</th>
                  <th className="p-4 font-semibold text-gray-200">Avg Price</th>
                  <th className="p-4 font-semibold text-gray-200">Current Price</th>
                  <th className="p-4 font-semibold text-gray-200">P&L (%)</th>
                  <th className="p-4 font-semibold text-gray-200">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {stocks.map((s, index) => {
                  const pl = ((s.currentPrice - s.avgPrice) / s.avgPrice) * 100;
                  return (
                    <tr key={`${s.symbol}-${index}`} className="hover:bg-white/5 transition-colors">
                      <td className="p-4 font-bold text-white tracking-wide">{s.symbol}</td>
                      <td className="p-4 text-gray-300">{s.quantity}</td>
                      <td className="p-4 text-gray-300">${s.avgPrice.toFixed(2)}</td>
                      <td className="p-4 font-medium text-white">${s.currentPrice.toFixed(2)}</td>
                      <td className={`p-4 font-bold ${pl >= 0 ? "text-emerald-400" : "text-rose-400"}`}>
                        {pl >= 0 ? "+" : ""}
                        {pl.toFixed(2)}%
                      </td>
                      <td className="p-4">
                        <button
                          onClick={() => handleUpdatePrice(s.symbol)}
                          className="bg-purple-600/20 hover:bg-purple-600/40 text-purple-300 hover:text-white border border-purple-500/30 px-3 py-1 rounded-lg text-sm transition-all"
                        >
                          Refresh
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>
      )}
    </Layout>
  );
}
