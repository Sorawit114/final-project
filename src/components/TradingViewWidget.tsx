import { useEffect } from "react";

export default function TradingViewChart({ symbol }: { symbol: string }) {
  useEffect(() => {
    const script = document.createElement("script");
    script.src = "https://s3.tradingview.com/tv.js";
    script.async = true;
    script.onload = () => {
      // @ts-ignore
      new window.TradingView.widget({
        container_id: `tradingview_${symbol}`,
        width: "100%",
        height: 600,
        symbol: symbol,
        interval: "1", // 1 minute
        timezone: "Asia/Bangkok",
        theme: "dark",
        style: "1",
        locale: "en",
        toolbar_bg: "#f1f3f6",
        enable_publishing: false,
        hide_top_toolbar: false,
        withdateranges: true,
        allow_symbol_change: true,
        hide_side_toolbar: false,
        details: true,
      });
    };
    document.getElementById(`tradingview_${symbol}`)?.appendChild(script);
    return () => {
      const container = document.getElementById(`tradingview_${symbol}`);
      if (container) container.innerHTML = "";
    };
  }, [symbol]);

  return <div id={`tradingview_${symbol}`} className="w-full" />;
}
