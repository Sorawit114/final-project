from fastapi import FastAPI
from pydantic import BaseModel
import tensorflow as tf
import numpy as np
import yfinance as yf
import pandas as pd

# โหลดโมเดล .keras
model = tf.keras.models.load_model("model/my_stock_prediction_cnn_lstm_model.keras")

app = FastAPI(title="Stock Pattern AI")

class Query(BaseModel):
    symbol: str
    start: str  # datetime string, เช่น '2025-12-01T09:30'
    end: str

def preprocess_stock(symbol: str, start: str, end: str):
    """
    ดึงราคาหุ้นจาก Yahoo Finance
    แล้วสร้าง feature สำหรับโมเดล
    """
    df = yf.download(symbol, start=start, end=end, interval="1m")  # 1 นาที
    if df.empty:
        return None
    # ตัวอย่าง preprocess: ใช้ close price normalize
    close_prices = df["Close"].values
    features = (close_prices - np.mean(close_prices)) / np.std(close_prices)
    # reshape เป็น [1, n_features]
    return features.reshape(1, -1)

@app.post("/predict")
def predict(query: Query):
    features = preprocess_stock(query.symbol, query.start, query.end)
    if features is None:
        return {"error": "No stock data available for this range."}

    # prediction โมเดล
    prediction = model.predict(features)
    # สมมติโมเดล output เป็น probability
    has_pattern = bool(prediction[0][0] > 0.5)

    return {
        "symbol": query.symbol,
        "start": query.start,
        "end": query.end,
        "has_pattern": has_pattern
    }
