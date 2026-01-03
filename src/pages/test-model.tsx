import { useState } from "react";
import Layout from "../components/Layout";

const availableSymbols = ["AAPL", "TSLA", "NVDA", "TISCO"];

export default function TestModel() {
  const [symbol, setSymbol] = useState(availableSymbols[0]);
  const [datetime, setDatetime] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const runModel = async () => {
    if (!datetime) {
      alert("กรุณาเลือกวันและเวลา (เวลาไทย)");
      return;
    }

    setLoading(true);
    setResult(null);

    try {
      const res = await fetch("http://localhost:8080/predict", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          symbol,
          datetime,
        }),
      });

      const data = await res.json();

      if (!res.ok || data.error) {
        setResult({ error: data.error || "Backend error" });
      } else {
        setResult(data);
      }
    } catch (err) {
      setResult({ error: "Cannot connect to backend" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="bg-white/10 backdrop-blur-xl border border-white/20 p-8 rounded-2xl shadow-2xl max-w-xl mx-auto space-y-6 mt-10">
        <h2 className="text-2xl font-bold text-center mb-6">Market Prediction Model</h2>

        {/* Symbol */}
        <div>
          <label className="block mb-2 text-gray-300 font-medium">Select Symbol</label>
          <div className="relative">
            <select
              className="w-full p-3 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-purple-500 appearance-none transition-colors"
              value={symbol}
              onChange={(e) => setSymbol(e.target.value)}
            >
              {availableSymbols.map((s) => (
                <option key={s} className="bg-gray-800 text-white">{s}</option>
              ))}
            </select>
            <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-4 h-4">
                <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </div>
          </div>
        </div>

        {/* Datetime */}
        <div>
          <label className="block mb-2 text-gray-300 font-medium">Select Date & Time (Thai Time)</label>
          <input
            type="datetime-local"
            className="w-full p-3 rounded-lg bg-white/5 border border-white/10 text-white focus:outline-none focus:border-purple-500 transition-colors [color-scheme:dark]"
            value={datetime}
            onChange={(e) => setDatetime(e.target.value)}
          />
          <p className="text-sm text-gray-400 mt-2 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-4 h-4">
              <path strokeLinecap="round" strokeLinejoin="round" d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z" />
            </svg>
            Analyzes 5 past 1-minute candles
          </p>
        </div>

        {/* Button */}
        <button
          onClick={runModel}
          disabled={loading}
          className="w-full bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-400 hover:to-teal-400 text-white py-3 rounded-lg font-bold shadow-lg shadow-emerald-500/20 transform hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:transform-none mt-4"
        >
          {loading ? "Analyzing Market..." : "Run Prediction Model"}
        </button>

        {/* Loading */}
        {loading && (
          <div className="text-center text-emerald-400 animate-pulse font-medium">
            Processing market data...
          </div>
        )}

        {/* Result */}
        {result && (
          <div className={`p-6 rounded-xl border ${result.error ? "bg-red-500/10 border-red-500/30" : "bg-emerald-500/10 border-emerald-500/30"} mt-6 animate-fade-in`}>
            {result.error ? (
              <div className="text-red-400 flex items-center gap-2 justify-center font-medium">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
                </svg>
                {result.error}
              </div>
            ) : (
              <div className="space-y-3">
                <div className="flex justify-between items-center border-b border-white/10 pb-2">
                  <span className="text-gray-400">Symbol</span>
                  <span className="text-xl font-bold">{result.symbol}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Pattern Detected</span>
                  <span className="font-semibold text-emerald-300">{result.pattern_name}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Confidence</span>
                  <div className="flex items-center gap-2">
                    <div className="w-24 h-2 bg-gray-700/50 rounded-full overflow-hidden">
                      <div className="h-full bg-emerald-400" style={{ width: `${result.confidence * 100}%` }}></div>
                    </div>
                    <span className="text-sm font-medium">{(result.confidence * 100).toFixed(1)}%</span>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4 mt-2 text-sm bg-white/5 p-3 rounded-lg">
                  <div>
                    <span className="text-gray-500 block text-xs uppercase tracking-wider">Thai Time</span>
                    {result.thai_time}
                  </div>
                  <div className="text-right">
                    <span className="text-gray-500 block text-xs uppercase tracking-wider">US Time</span>
                    {result.us_time}
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </Layout>
  );
}
