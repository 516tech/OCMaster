# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

超频大师 (OCMaster) — a cross-platform hardware info collection tool + merchant-side overclocking suggestion platform. Two modules:

- **C端 (Client)**: Single-file portable desktop app (Electron) that scans hardware and optionally uploads data with a share code. Supports Windows / Linux / macOS.
- **S端 (Merchant Web)**: Admin platform where merchants look up hardware by share code, reference a hardware database, fill in overclocking suggestions, and export a PDF report.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| C端 Client | Electron 28+ / Vue 3 / C++ Native Addon (node-addon-api) |
| Backend API | Golang / chi / GORM / zerolog |
| Database | MySQL 8.0 |
| S端 Frontend | Vue 3 + Element Plus + TypeScript |
| DI | Google Wire (compile-time) |
| Deployment | Docker Compose (backend + MySQL + Nginx) |

## Architecture

```
C端 (Electron) ──HTTPS──→ Nginx ──→ Backend (Golang/chi) ──→ MySQL
                              ↑
S端 (Vue 3 SPA) ──HTTPS─────┘
```

## Backend: DDD Layered Architecture

```
s-end/backend/
├── domain/           # Entities, value objects, repository interfaces (no deps)
├── application/      # Use cases, service orchestration
├── infrastructure/   # GORM repos, JWT, PDF generation
├── interfaces/http/  # chi handlers, middleware, DTOs
├── pkg/config/       # viper config with defaults (no config file needed)
└── cmd/server/       # main.go + wire.go (DI entry point)
```

Wire injection chain: `Config → DB → Repos → Services → Handlers → Router → Server`

## Key Design Decisions

- **C端 single-file**: Electron packaged with electron-builder → .exe/.dmg/.AppImage.
- **C端 hardware detection**: C++ native addon per platform (WMI on Win, sysfs on Linux, IOKit on macOS), fallback to manual input for PSU/cooler.
- **S端 config**: all defaults in code, viper loads config file/env vars optionally. App starts without any config file.
- **Testing**: all tests use in-code defaults (no config file dependency). Integration tests use testcontainers-go for real MySQL.
- **Share codes**: 6-digit numeric codes, valid 7 days, user-revocable.
- **Data lifecycle**: hardware uploads deleted after 7 days, merchant query history deleted after 30 days, merchant accounts permanent (cancellable).
- **HTTPS everywhere**, no IP/geolocation collection, opt-in upload only.

## Hardware Info Collected (C端)

CPU (model, cores/threads, base freq, BIOS version), Motherboard (brand, model, chipset, BIOS version), RAM (total, per-stick capacity, count, frequency, timings, die type via SPD, channel count), GPU (model, VRAM), PSU (rated wattage via SMBus or manual input), Cooler (air/AIO, manual if undetectable).

## Merchant Platform (S端)

- Phone + password + bcrypt auth, admin-reviewed registration, JWT sessions.
- Share code lookup → view hardware → reference DB shows common OC ranges for that CPU/RAM die/cooler tier → fill template → generate PDF.
- Customizable risk-warning templates per merchant.

## API Routes

```
POST   /api/v1/hardware/upload
GET    /api/v1/hardware/:code
DELETE /api/v1/hardware/:code
POST   /api/v1/merchant/register
POST   /api/v1/merchant/login
GET    /api/v1/merchant/profile
PUT    /api/v1/merchant/profile
PUT    /api/v1/merchant/password
POST   /api/v1/suggestions
GET    /api/v1/suggestions/:id
GET    /api/v1/suggestions/history
GET    /api/v1/suggestions/:id/pdf
GET    /api/v1/reference/cpu
GET    /api/v1/reference/ram
GET    /api/v1/reference/cooler
```
