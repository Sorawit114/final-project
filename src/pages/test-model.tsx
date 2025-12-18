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
      const res = await fetch("http://localhost:8000/predict", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ symbol, datetime }),
      });

      const data = await res.json();
      setResult(data);
    } catch (err) {
      setResult({ error: "Cannot connect to backend" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <h1 className="text-2xl font-bold mb-4">Test AI Pattern Model</h1>

      <div className="bg-gray-800 p-6 rounded-xl max-w-xl mx-auto space-y-4">
        {/* Symbol */}
        <select
          className="w-full p-2 rounded bg-gray-700"
          value={symbol}
          onChange={(e) => setSymbol(e.target.value)}
        >
          {availableSymbols.map((s) => (
            <option key={s}>{s}</option>
          ))}
        </select>

        {/* Datetime */}
        <div>
          <label className="block mb-1">
            เลือกวันและเวลา (เวลาไทย)
          </label>
          <input
            type="datetime-local"
            className="w-full p-2 rounded bg-gray-700"
            value={datetime}
            onChange={(e) => setDatetime(e.target.value)}
          />
          <p className="text-sm text-gray-400 mt-1">
            ระบบจะวิเคราะห์แท่ง 1 นาที ย้อนหลัง 5 แท่ง<br />
            เวลาที่แนะนำ: ประมาณ 20:30 – 03:00 (ตลาด US เปิด)
          </p>
        </div>

        {/* Button */}
        <button
          onClick={runModel}
          disabled={loading}
          className="w-full bg-green-500 hover:bg-green-400 p-2 rounded font-semibold disabled:opacity-50"
        >
          Run Model
        </button>

        {/* Loading */}
        {loading && (
          <div className="text-center text-green-400">
            Running model...
          </div>
        )}

        {/* Result */}
        {result && (
          <div className="bg-gray-700 p-4 rounded whitespace-pre-wrap">
            {result.error ? (
              <span className="text-red-400">
                ❌ {result.error}
              </span>
            ) : (
              <>
                <div>📈 Symbol: {result.symbol}</div>
                <div>🧠 Pattern: {result.pattern_name}</div>
                <div>
                  🎯 Confidence: {(result.confidence * 100).toFixed(2)}%
                </div>
                <div>🇹🇭 Thai Time: {result.thai_time}</div>
                <div>🇺🇸 US Time: {result.us_time}</div>
              </>
            )}
          </div>
        )}
      </div>
    </Layout>
  );
}
