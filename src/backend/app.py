from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import tensorflow as tf
import numpy as np
import yfinance as yf
from datetime import datetime
import pytz

app = FastAPI()

# ---------- CORS ----------
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ---------- TIMEZONE ----------
TH_TZ = pytz.timezone("Asia/Bangkok")
US_TZ = pytz.timezone("US/Eastern")

# ---------- LOAD MODEL ----------
model = tf.keras.models.load_model(
    "model/my_stock_prediction_cnn_lstm_model.keras"
)

# ---------- REQUEST ----------
class Query(BaseModel):
    symbol: str
    datetime: str  # เวลาไทย เช่น "2024-01-05T21:30"

# ---------- HELPERS ----------
def prepare_input(df):
    df = df[["Open", "High", "Low", "Close"]]

    # normalize (ต้องเหมือนตอน train)
    df = (df - df.mean()) / df.std()

    if len(df) < 5:
        return None

    window = df.values[-5:]  # 5 แท่งล่าสุด
    return window.reshape(1, 5, 4)

# ---------- API ----------
@app.post("/predict")
def predict(query: Query):
    # 1️⃣ แปลงเวลาไทย → US
    try:
        th_time = TH_TZ.localize(datetime.fromisoformat(query.datetime))
        us_time = th_time.astimezone(US_TZ).replace(tzinfo=None)
    except:
        return {"error": "Invalid datetime format"}

    # 2️⃣ กันเลือกเวลาอนาคต (US)
    now_us = datetime.now(US_TZ).replace(tzinfo=None)
    if us_time > now_us:
        return {"error": "Selected time is in the future (US market time)"}

    # 3️⃣ เช็กเวลาตลาด
    market_open = us_time.replace(hour=9, minute=30)
    market_close = us_time.replace(hour=16, minute=0)

    if us_time < market_open or us_time > market_close:
        return {
            "error": "Selected time is outside US market hours (09:30–16:00 ET)"
        }

    # 4️⃣ โหลดข้อมูลหุ้นทั้งวัน (US)
    start_day = us_time.replace(hour=0, minute=0, second=0)
    end_day = us_time.replace(hour=23, minute=59, second=59)

    df = yf.download(
        query.symbol,
        start=start_day,
        end=end_day,
        interval="1m",
        progress=False
    )

    if df.empty:
        return {"error": "No data available for this day"}

    df.index = df.index.tz_localize(None)

    # 5️⃣ เอาแท่งก่อนเวลาที่เลือก
    df = df[df.index <= us_time]

    if len(df) < 5:
        return {"error": "Not enough candles before selected time"}

    features = prepare_input(df)
    if features is None:
        return {"error": "Invalid input window"}

    # 6️⃣ Predict
    prediction = model.predict(features, verbose=0)
    confidence = float(prediction[0][0])

    return {
        "symbol": query.symbol,
        "pattern_name": "Custom Pattern",
        "confidence": round(confidence, 4),
        "thai_time": th_time.strftime("%Y-%m-%d %H:%M"),
        "us_time": us_time.strftime("%Y-%m-%d %H:%M")
    }
