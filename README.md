<p align="center">
  <img src="logo.png" width="250" alt="BlueIce Logo">
</p>

# 🧊 BlueIce Server
**A high performance Minecraft Server implementation written in Go.**

This project is work in progress.

The server doesn't have any real features yet, but we just finished the login and configuration state so the player
can join the server. There is no online mode verification yet and nothing to do on the server since the player
just spawns in an empty world. But we're working on it.

## 🚀 The Vision
BlueIce isn't just another server. We are building:
* **Cloud-Native Architecture:** Designed for horizontal scaling.
* **WASM Plugin System:** High-performance, language-agnostic plugins.
* **Go Core:** Leveraging Go's concurrency for ultra-low latency.

## 🛠 Current Status
- [x] Handshake & Login State
- [x] NBT-Registry Writer (Custom Implementation)
- [x] Configuration State
- [x] **Joinable Void** (The player can officially connect!)
- [ ] Chunk Rendering
- [ ] World Saving and Loading

## 🛠 Getting Started
Since this is a Go project, you can try it out by cloning the repo and running:

```bash
go run cmd/blueice/main.go
```

## Credits
Nano Banana 2 generated the server logo. It features the Go-Gopher design created by Renee French
under the CC-BY 4.0 license.
