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
  const [newStock, setNewStock] = useState<Stock>({
    symbol: "",
    quantity: 0,
    avgPrice: 0,
  });

  // โหลดข้อมูลจาก localStorage
  useEffect(() => {
    const stored = localStorage.getItem("myStocks");
    if (stored) {
      const parsed: StockWithPrice[] = JSON.parse(stored);
      setStocks(parsed);
    }
  }, []);

  // บันทึกไป localStorage
  useEffect(() => {
    localStorage.setItem("myStocks", JSON.stringify(stocks));
  }, [stocks]);

  // mock function ดึงราคาปัจจุบัน (แทน TradingView API)
  const fetchCurrentPrice = (symbol: string) => {
    // สำหรับ mock ใช้ random ราคา
    const stock = stocks.find((s) => s.symbol === symbol);
    return stock ? stock.avgPrice * (0.9 + Math.random() * 0.2) : 0;
  };

  const handleAddStock = () => {
    if (!newStock.symbol) return;
    const stockWithPrice: StockWithPrice = {
      ...newStock,
      currentPrice: fetchCurrentPrice(newStock.symbol),
    };
    setStocks([...stocks, stockWithPrice]);
    setNewStock({ symbol: "", quantity: 0, avgPrice: 0 });
  };

  const handleUpdatePrice = (symbol: string) => {
    setStocks((prev) =>
      prev.map((s) =>
        s.symbol === symbol
          ? { ...s, currentPrice: fetchCurrentPrice(symbol) }
          : s
      )
    );
  };

  return (
    <Layout>

      <h1 className="text-2xl font-bold mb-4">My Stock</h1>

      {/* ตารางหุ้น */}
      <table className="w-full table-auto text-left border-collapse">
        <thead>
          <tr>
            <th className="border-b border-gray-700 p-2">Symbol</th>
            <th className="border-b border-gray-700 p-2">Quantity</th>
            <th className="border-b border-gray-700 p-2">Avg Price</th>
            <th className="border-b border-gray-700 p-2">Current Price</th>
            <th className="border-b border-gray-700 p-2">P&L (%)</th>
            <th className="border-b border-gray-700 p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((s) => {
            const pl = ((s.currentPrice - s.avgPrice) / s.avgPrice) * 100;
            return (
              <tr key={s.symbol}>
                <td className="border-b border-gray-700 p-2">{s.symbol}</td>
                <td className="border-b border-gray-700 p-2">{s.quantity}</td>
                <td className="border-b border-gray-700 p-2">
                  {s.avgPrice.toFixed(2)}
                </td>
                <td className="border-b border-gray-700 p-2">
                  {s.currentPrice.toFixed(2)}
                </td>
                <td
                  className={`border-b border-gray-700 p-2 ${
                    pl >= 0 ? "text-green-500" : "text-red-500"
                  }`}
                >
                  {pl >= 0 ? "+" : ""}
                  {pl.toFixed(2)}%
                </td>
                <td className="border-b border-gray-700 p-2">
                  <button
                    onClick={() => handleUpdatePrice(s.symbol)}
                    className="bg-gray-700 px-2 py-1 rounded hover:bg-gray-600"
                  >
                    Update
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </Layout>
  );
}
