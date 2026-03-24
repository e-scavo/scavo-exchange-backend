# 🏗️ Architecture Overview

## 🎯 Goal

Build a backend capable of supporting:

- DEX (initial)
- Hybrid exchange (future)

---

## 🧱 Architecture Style

Modular Monolith

---

## 🔌 Components

- API Layer (REST + WS)
- Domain Modules
- Blockchain Services
- Workers
- Persistence Layer

---

## 🔁 Communication

- Internal: function calls
- External: HTTP + WS
- Async: workers + queues (future)

---

## 🔗 Blockchain

- Direct RPC to SCAVIUM
- Optional dedicated RPC

---

## 🔐 Security

- JWT-based auth
- Wallet signatures
- Rate limiting (future)

---

## 📦 Deployment

- Single binary
- Dockerized