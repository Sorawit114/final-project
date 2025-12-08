import { useState } from "react";
import React from "react";
import Link from "next/link";
import TradingViewWidget from "../components/TradingViewWidget";
import MiniChart from "../components/MiniChart";
import { PlusCircleIcon } from "@heroicons/react/24/solid";
import Layout from "../components/Layout";

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
  const [showLogin, setShowLogin] = useState(false);

  const handleAddStock = (symbol: string) => {
    const stored = localStorage.getItem("myStocks");
    const stocks = stored ? JSON.parse(stored) : [];
    if (!stocks.find((s: any) => s.symbol === symbol)) {
      const newStock = { symbol, quantity: 0, avgPrice: 0, currentPrice: 0 };
      localStorage.setItem("myStocks", JSON.stringify([...stocks, newStock]));
      alert(`${symbol} added to My Stock!`);
    } else {
      alert(`${symbol} is already in My Stock.`);
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
        <div className="p-4 overflow-x-auto flex gap-4 justify-between">
          {symbols.map((sym) => (
            <div
              key={sym}
              onClick={() => setSelectedSymbol(sym)}
              className="cursor-pointer relative"
            >
              <MiniChart symbol={sym} />
              <div className="flex justify-between items-center mt-1">
                <p className="text-center w-full">{sym}</p>
                <button
                  onClick={() => handleAddStock(sym)}
                  className="bg-indigo-500 hover:bg-indigo-400 text-white p-1 rounded-full shadow-lg flex items-center justify-center"
                  title="Add to My Stock"
                >
                  <PlusCircleIcon className="h-5 w-5" />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Advanced Chart */}
      <div className="p-4">
        <TradingViewWidget symbol={selectedSymbol} />
      </div>

      {showLogin && (
        <>
          {/* Backdrop blur */}
          <div className="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm z-40"></div>

          {/* Modal */}
          <div className="fixed inset-0 flex items-center justify-center z-50">
            <div className="bg-gray-800 p-8 rounded-xl shadow-lg w-80">
              <h3 className="text-xl font-bold mb-4 text-center">Login</h3>
              <input
                type="text"
                placeholder="Username"
                className="w-full p-2 mb-3 rounded-md bg-gray-700 text-white"
              />
              <input
                type="password"
                placeholder="Password"
                className="w-full p-2 mb-4 rounded-md bg-gray-700 text-white"
              />
              <button className="w-full bg-indigo-600 hover:bg-indigo-500 py-2 rounded-md mb-2">
                Login
              </button>
              <p className="text-sm text-gray-400 text-center">
                Don't have an account?{" "}
                <a href="/register" className="text-indigo-400 hover:underline">
                  Register
                </a>
              </p>
              <button
                onClick={() => setShowLogin(false)}
                className="absolute top-2 right-2 text-gray-400 hover:text-white font-bold text-lg"
              >
                &times;
              </button>
            </div>
          </div>
        </>
      )}
    </Layout>
  );
}
