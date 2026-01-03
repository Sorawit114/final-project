import { useState, useEffect } from "react";
import React from "react";
import Link from "next/link";
import TradingViewWidget from "../components/TradingViewWidget";
import MiniChart from "../components/MiniChart";
import { PlusCircleIcon } from "@heroicons/react/24/solid";
import Layout from "../components/Layout";
import { useToast } from "@/context/ToastContext";

const symbols = [
  "AAPL",
  "TSLA",
  "NVDA",
  "MSFT",
  "AMZN",
  "GOOGL",
  "META",
  "NFLX",
];

export default function Home() {
  const [selectedSymbol, setSelectedSymbol] = useState(symbols[0]);
  const [userId, setUserId] = useState<string | null>(null);

  useEffect(() => {
    // Check for logged in user
    const storedUserId = localStorage.getItem("userId");
    if (storedUserId) {
      setUserId(storedUserId);
    }
  }, []);

  const { toast } = useToast();

  const handleAddStock = async (symbol: string) => {
    if (!userId) {
      toast("Please login to add stocks.", "error");
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/add-stock", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userId: parseInt(userId),
          stockName: symbol,
        }),
      });

      const data = await res.json();

      if (res.ok) {
        toast(`${symbol} added to your watchlist!`, "success");
      } else {
        toast(`Failed to add ${symbol}: ${data.error}`, "error");
      }
    } catch (err) {
      toast("Error connecting to server.", "error");
    }
  };

  return (
    <Layout>
      {/* Mini Charts */}
      <div className="relative flex flex-col items-center justify-center text-center py-20 px-4 md:px-0">
        <h2 className="text-4xl md:text-5xl font-bold mb-4">
          Welcome to My Stock Dashboard
        </h2>
        <p className="text-gray-400 mb-8 max-w-xl">
          Analyze your favorite stocks, check patterns, and track your portfolio
          easily.
        </p>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6 w-full max-w-6xl">
          {symbols.map((sym) => (
            <div
              key={sym}
              onClick={() => setSelectedSymbol(sym)}
              className={`cursor-pointer relative p-4 rounded-2xl bg-white/5 border border-white/10 hover:bg-white/10 transition-all group ${selectedSymbol === sym ? "ring-2 ring-purple-500 bg-white/10" : ""
                }`}
            >
              <MiniChart symbol={sym} />
              <div className="flex justify-between items-center mt-3 px-1">
                <p className={`text-center font-bold transition-colors ${selectedSymbol === sym ? "text-purple-300" : "text-gray-300 group-hover:text-white"}`}>
                  {sym}
                </p>
                {userId && (
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleAddStock(sym);
                    }}
                    className="bg-indigo-600 hover:bg-indigo-500 text-white p-2 rounded-full shadow-lg flex items-center justify-center transition-transform hover:scale-110 hover:shadow-indigo-500/50"
                    title="Add to My Stock"
                  >
                    <PlusCircleIcon className="h-5 w-5" />
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Advanced Chart */}
      <div className="p-4 max-w-7xl mx-auto h-[600px]">
        <TradingViewWidget symbol={selectedSymbol} />
      </div>
    </Layout>
  );
}
