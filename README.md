# SCAVO Exchange - Backend

## Overview

SCAVO Exchange is a hybrid exchange platform designed to combine decentralized (DEX) and centralized (CEX) models into a unified ecosystem.

This repository contains the backend system built in Go which powers:

* DEX (initial focus)
* Future hybrid trading (DEX + CEX)
* Wallet integration (self-custody first)
* Blockchain interaction (SCAVIUM network)
* Real-time communication (WebSocket)

## Current Focus

The backend is currently in DEX-first mode prioritizing:

* Non-custodial architecture
* AMM-based trading (v1)
* Direct blockchain interaction (SCAVIUM)
* Wallet-native flows

## Architecture Style

* Modular Monolith
* Event-aware backend
* REST and WebSocket APIs
* Background workers
* Blockchain-integrated services

## Tech Stack

Language: Go  
API: REST and WebSocket  
Database: PostgreSQL (baseline integrated)  
Cache: Redis (scaffolded)  
Blockchain: SCAVIUM (EVM)  
Contracts: Solidity  
Infrastructure: Docker Compose

## Core Principles

* Self-custody first
* Security by design
* Deterministic backend behavior
* Observable systems
* Phase-driven development

## Project Structure (current)

- `cmd/`
- `internal/app`
- `internal/core`
- `internal/modules`
- `migrations`
- `scripts`
- `docs`

## Roadmap

The project is structured in stages:

- Stage 0: Foundation
- Stage 1: Identity and Wallets
- Stage 2: Blockchain Integration
- Stage 3: Smart Contracts (DEX)
- Stage 4: DEX Backend Logic
- Stage 5: APIs and Realtime
- Stage 6: Hybrid Expansion
- Stage 7: Security and Operations
- Stage 8: Testing and Internal Release

See full roadmap in `docs/roadmap.md`.

## Documentation

All documentation is located in `/docs`:

* Architecture
* Flows
* Decisions
* Development
* Testing
* Roadmap
* Handoff

## Development Status

Current phase after this subphase update:

**Stage 0 - Foundation**  
**Phase 0.4 - Auth and User Stabilization**  
**Subphase 0.4.2 - Token Lifecycle and Auth Transport Hardening**

Implemented in this subphase:

- auth service boundary preserved from 0.4.1
- shared token extraction introduced for HTTP and WebSocket transports
- auth claims context formalized in core auth layer
- reusable HTTP auth middleware introduced
- `GET /auth/me` now protected through middleware
- WebSocket auth attachment aligned with the same token extraction rules
- auth transport responsibility clarified without introducing session persistence yet

## Workflow Rules

* The ZIP project is the source of truth
* Documentation must always match implementation
* Work is phase-driven
* Each step includes a commit reference

## Future Scope

SCAVO Exchange will evolve into a hybrid system including:

* DEX (non-custodial)
* CEX (custodial accounts)
* P2P trading
* Fiat ramps

## Notes

This backend is designed to be:

* Frontend-ready
* Scalable
* Secure
* Blockchain-native

## Next Step

Phase 0.4.3 - Session Evolution and Wallet Auth Preparation