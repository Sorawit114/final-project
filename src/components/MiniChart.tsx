import React, { useEffect, useRef, memo } from "react";

interface TradingViewWidgetProps {
  symbol: string;
}

function TradingViewWidget({ symbol }: TradingViewWidgetProps) {
  const container = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!container.current) return;

    // เคลียร์ container ก่อน
    container.current.innerHTML = "";

    const script = document.createElement("script");
    script.src =
      "https://s3.tradingview.com/external-embedding/embed-widget-mini-symbol-overview.js";
    script.type = "text/javascript";
    script.async = true;
    script.innerHTML = JSON.stringify({
      symbol: symbol,
      chartOnly: false,
      dateRange: "12M",
      noTimeScale: false,
      colorTheme: "dark",
      isTransparent: false,
      locale: "en",
      width: "100%",
      height: 180, // เพิ่มความสูงให้เห็นกราฟชัดขึ้น
      autosize: false, // ปิด autosize เพื่อให้ width/height ใช้งานจริง
      showVolume: false,
    });
    container.current.appendChild(script);
  }, [symbol]);

  return (
    <div className="tradingview-widget-container w-[200px] flex-shrink-0">
      <div className="tradingview-widget-container__widget" ref={container}></div>
      <div className="tradingview-widget-copyright text-xs text-gray-400">
      </div>
    </div>
  );
}

export default memo(TradingViewWidget);
