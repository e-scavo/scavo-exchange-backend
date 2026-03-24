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
Database: PostgreSQL (planned)
Cache: Redis (planned)
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

cmd/
internal/
pkg/

The structure will evolve during Phase 0.2

## Roadmap

The project is structured in stages:

Stage 0: Foundation
Stage 1: Identity and Wallets
Stage 2: Blockchain Integration
Stage 3: Smart Contracts (DEX)
Stage 4: DEX Backend Logic
Stage 5: APIs and Realtime
Stage 6: Hybrid Expansion
Stage 7: Security and Operations
Stage 8: Testing and Internal Release

See full roadmap in docs/roadmap.md

## Documentation

All documentation is located in /docs:

* Architecture
* Flows
* Decisions
* Development
* Roadmap

## Development Status

Current Phase:

Phase 0.1.1 - Baseline Audit and Documentation Foundation

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

Phase 0.1.2 - Architecture Definition
