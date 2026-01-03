# Final Project

This project consists of a **Next.js Frontend**, a **Go Backend (Authentication & User Data)**, and a **Python Service (Deep Learning Model for Stock Prediction)**.

## Prerequisites

Ensure you have the following installed:
- **Node.js** (v18 or later)
- **Go** (v1.20 or later)
- **Python** (v3.10 or later)
- **PostgreSQL** (running on port 5432)

---

## 1. Database Setup (PostgreSQL)

Ensure PostgreSQL is running and creates a database named `Project`.

The Go Backend expects the following connection string by default:
`postgres://postgres:1234@localhost:5432/Project`

If your credentials differ, update the `.env` file or `Config/config.go` in the backend.

---

## 2. Backend Setup (Go)

This service handles User Authentication (Login/Register) and Stock Portfolio management.

### **Installation & Run**

```bash
cd src/backend
# Install dependencies
go mod tidy

# Run the server
go run main.go
# OR if compiled:
# ./server.exe
```

> **Port:** `8080`

---

## 3. Model Service Setup (Python)

This service runs the CNN-LSTM model for stock pattern prediction.

### **Installation & Run**

```bash
cd src/backend
# Create virtual environment (optional but recommended)
python -m venv venv
# Activate venv
# Windows:
.\venv\Scripts\Activate.ps1
# Mac/Linux:
# source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Run the FastAPI server
uvicorn app:app --reload --host 0.0.0.0 --port 8000
```

> **Port:** `8000`

---

## 4. Frontend Setup (Next.js)

This is the user interface.

### **Installation & Run**

```bash
# Root directory
npm install

# Run development server
npm run dev
```

> **Port:** `3000`
> Open [http://localhost:3000](http://localhost:3000)

---

## Usage Guide

1.  **Start Database:** Ensure Postgres is running.
2.  **Start Go Backend:** Terminal 1 -> `go run main.go`
3.  **Start Python Model:** Terminal 2 -> `uvicorn app:app --reload`
4.  **Start Frontend:** Terminal 3 -> `npm run dev`
5.  **Access:** Go to `localhost:3000`. Login/Register to add stocks to your portfolio.
