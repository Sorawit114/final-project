import { useState } from "react";
import Link from "next/link";
import Layout from "../components/Layout";

const availableSymbols = ["AAPL", "TSLA", "NVDA", "TISCO"];

export default function TestModel() {
  const [selectedSymbol, setSelectedSymbol] = useState(availableSymbols[0]);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);

  const handleRunModel = async () => {
    if (!selectedSymbol || !startDate || !endDate) {
      alert("Please select symbol and date range.");
      return;
    }

    setLoading(true);
    setResult(null);

    try {
      const response = await fetch("http://localhost:8000/predict", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          symbol: selectedSymbol,
          start: startDate,
          end: endDate,
        }),
      });

      const data = await response.json();
      if (data.error) {
        setResult(`❌ ${data.error}`);
      } else {
        setResult(
          data.has_pattern
            ? `✅ Pattern detected for ${data.symbol} from ${data.start} to ${data.end}`
            : `❌ No pattern detected for ${data.symbol} from ${data.start} to ${data.end}`
        );
      }
    } catch (error) {
      setResult("❌ Error calling backend API  " + endDate);
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <h1 className="text-2xl font-bold mb-4">Test AI Model Pattern</h1>

      <div className="bg-gray-800 p-6 rounded-xl shadow-lg max-w-2xl mx-auto flex flex-col gap-4">
        {/* Select Symbol */}
        <div className="flex flex-col">
          <label className="font-semibold mb-1">Select Stock</label>
          <select
            className="p-2 rounded-md bg-gray-700 text-white"
            value={selectedSymbol}
            onChange={(e) => setSelectedSymbol(e.target.value)}
          >
            {availableSymbols.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
        </div>

        {/* Date Range */}
        <div className="flex gap-4">
          <div className="flex flex-col w-1/2">
            <label className="font-semibold">Start Date</label>
            <input
              type="datetime-local"
              className="p-2 rounded-md bg-gray-700 text-white"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
            />
          </div>
          <div className="flex flex-col w-1/2">
            <label className="font-semibold">End Date</label>
            <input
              type="datetime-local"
              className="p-2 rounded-md bg-gray-700 text-white"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
          </div>
        </div>

        {/* Run Button */}
        <button
          onClick={handleRunModel}
          className="bg-green-500 hover:bg-green-400 px-4 py-2 rounded-md font-semibold transition"
        >
          Run Model
        </button>

        {/* Spinner */}
        {loading && (
          <div className="flex justify-center mt-4">
            <div className="w-10 h-10 border-4 border-green-400 border-t-transparent rounded-full animate-spin"></div>
          </div>
        )}

        {/* Result */}
        {result && (
          <div className="bg-gray-700 p-4 rounded-md mt-4 text-green-300 font-medium whitespace-pre-wrap">
            {result}
          </div>
        )}
      </div>
    </Layout>
  );
}
