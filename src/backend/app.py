from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import tensorflow as tf
import yfinance as yf
from datetime import datetime
from dateutil import parser
import pytz
import numpy as np
import pandas as pd

# ---------- APP ----------
app = FastAPI()

# ---------- CORS ----------
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://127.0.0.1:3000",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ---------- TIMEZONE ----------
TH_TZ = pytz.timezone("Asia/Bangkok")
US_TZ = pytz.timezone("US/Eastern")

# ---------- MODEL ----------
model = tf.keras.models.load_model(
    "model/my_stock_prediction_cnn_lstm_model.keras"
)

ALLOWED_SYMBOLS = [
  "AAPL",
  "TSLA",
  "NVDA",
  "MSFT",
  "AMZN",
  "GOOGL",
  "META",
  "NFLX",
]

# 🔑 IMPORTANT: ต้องตรงกับตอน train 100%
PATTERN_CLASSES = [
    'Doji', 
    'Spinning Top', 
    'Harami', 
    'Engulfing',
    'Hammer', 
    'Three Outside Up/Down', 
    'Three Inside Up/Down'
]

# optional: threshold ต่อ pattern
PATTERN_THRESHOLDS = {
    "Doji": 0.60,
    "Spinning Top": 0.60,
    "Harami": 0.65,
    "Engulfing": 0.70,
    "Hammer": 0.65,
    "Three Outside Up/Down": 0.75,
    "Three Inside Up/Down": 0.75
}

# ---------- REQUEST ----------
class Query(BaseModel):
    symbol: str
    datetime: str  # "2025-01-05T21:30"

# ---------- HELPERS ----------
def prepare_input(df):
    df = df[["Open", "High", "Low", "Close"]]

    if len(df) < 5:
        return None

    window = df.values[-5:]

    # normalize per-window
    mean = window.mean(axis=0)
    std = window.std(axis=0) + 1e-8
    window = (window - mean) / std

    return window.reshape(1, 5, 4)


def is_market_open(us_time: datetime):
    market_open = us_time.replace(hour=9, minute=30, second=0)
    market_close = us_time.replace(hour=16, minute=0, second=0)
    return market_open <= us_time <= market_close


# ---------- API ----------
@app.post("/predict")
def predict(query: Query):

    # 1️⃣ validate symbol
    if query.symbol not in ALLOWED_SYMBOLS:
        return {"error": "Invalid symbol"}

    # 2️⃣ parse Thai datetime → US
    try:
        naive_time = parser.parse(query.datetime)
        th_time = TH_TZ.localize(naive_time)
        us_time = th_time.astimezone(US_TZ)
    except Exception:
        return {"error": "Invalid datetime format"}

    # 3️⃣ future check
    if us_time > datetime.now(US_TZ):
        return {"error": "Selected time is in the future (US time)"}

    # 4️⃣ market hours check
    if not is_market_open(us_time):
        return {"error": "Outside US market hours (09:30 PM – 04:00 AM)"}

    # 5️⃣ download intraday data
    start = us_time.replace(hour=0, minute=0, second=0)
    end = us_time.replace(hour=23, minute=59, second=59)

    df = yf.download(
        query.symbol,
        start=start,
        end=end,
        interval="1m",
        progress=False,
    )

    if df.empty:
        return {"error": "No data for selected day"}

    # 🔧 FIX: Handle MultiIndex or Duplicate Columns from yfinance
    if isinstance(df.columns, pd.MultiIndex):
        df.columns = df.columns.get_level_values(0)
    
    # Remove duplicate columns (keep first)
    df = df.loc[:, ~df.columns.duplicated()]

    # normalize timezone
    if df.index.tz is not None:
        df.index = df.index.tz_convert(US_TZ).tz_localize(None)
    else:
        df.index = df.index.tz_localize(None)

    # 6️⃣ candles before selected time
    us_time_naive = us_time.replace(tzinfo=None)
    df = df[df.index <= us_time_naive]

    if len(df) < 5:
        return {"error": "Not enough candles before selected time"}

    features = prepare_input(df)
    if features is None:
        return {"error": "Invalid input window"}

    # 7️⃣ predict (MULTI-CLASS)
    probs = model.predict(features, verbose=0)[0]  # shape (9,)

    pattern_index = int(np.argmax(probs))
    pattern_name = PATTERN_CLASSES[pattern_index]
    confidence = float(probs[pattern_index])

    threshold = PATTERN_THRESHOLDS.get(pattern_name, 0.6)
    is_pattern = pattern_name != "No Pattern" and confidence >= threshold

    # 8️⃣ response
    return {
        "symbol": query.symbol,
        "pattern_name": pattern_name,
        "pattern_index": pattern_index,
        "is_pattern": is_pattern,
        "confidence": round(confidence, 4),
        "thai_time": th_time.strftime("%Y-%m-%d %H:%M"),
        "us_time": us_time.strftime("%Y-%m-%d %H:%M"),
        "all_probabilities": {
            PATTERN_CLASSES[i]: round(float(probs[i]), 4)
            for i in range(len(PATTERN_CLASSES))
        },
    }
