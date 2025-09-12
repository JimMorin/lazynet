# LazyNet - A Modern Terminal-Based Network Monitor

**LazyNet** is an interactive, terminal-based utility for monitoring network connections and processes.  
It provides a real-time, filterable, and sortable view of your system's network activity, with integrated actions to manage processes and diagnose connections.

---

## Table of Contents

- [LazyNet - A Modern Terminal-Based Network Monitor](#lazynet---a-modern-terminal-based-network-monitor)
  - [Table of Contents](#table-of-contents)
  - [Main Interface](#main-interface)
  - [Features](#features)
    - [Process Management](#process-management)
    - [Network Diagnostics](#network-diagnostics)
    - [Export to CSV](#export-to-csv)
    - [In-App Help Panel](#in-app-help-panel)
  - [Help Panel \& Ping Utility](#help-panel--ping-utility)
  - [Installation](#installation)
- [Clone the repository](#clone-the-repository)
- [Tidy dependencies](#tidy-dependencies)
- [Build the binary (using Makefile)](#build-the-binary-using-makefile)

---

## Main Interface

The main view provides a **two-pane layout**:

- **Left Pane:** Live-updating list of connections  
- **Right Pane:** Detailed information for the selected connection  

This allows you to monitor and analyze network activity efficiently.

---

## Features

- **Live Monitoring:** Automatically refreshes the connection list every few seconds.  
- **Interactive TUI:** Fully keyboard-driven interface with optional mouse support.  
- **Two-Pane Layout:** View all connections and details simultaneously.  
- **Real-time Filtering:** Press `/` to instantly filter connections by process name, PID, status, or address.  
- **Column Sorting:**  
  - Press `s` to cycle through sortable columns  
  - Press `S` to toggle ascending/descending order  
- **Color-Coded Status:** Quickly identify connection states (`ESTABLISHED`, `LISTEN`, `CLOSE_WAIT`) with colors.  

### Process Management

- `k` – Gracefully kill the process associated with the selected connection  
- `K` – Forcefully kill the process (SIGKILL)  

### Network Diagnostics

- `p` – Ping the remote address of a connection and view live output in a modal  

### Export to CSV

- `e` – Export the currently visible connections to a `lazynet_export.csv` file for analysis  

### In-App Help Panel

- `h` – Toggle a detailed, colorful panel explaining all keybindings  

---

## Help Panel & Ping Utility

The in-app help (`h`) provides a clear overview of all commands.  
Actions like ping (`p`) run in a convenient modal overlay, making diagnostics easier.

---

## Installation

### Option 1: Build from Source

Ensure you have **Go (version 1.21 or newer)** installed on your system.

```bash
# Clone the repository
git clone https://github.com/sureshkrishnan-v/lazynet.git
cd lazynet

# Tidy dependencies
go mod tidy

# Build the binary (using Makefile)
make build
```

### Option 2: Using Nix

If you have [Nix](https://nixos.org/download.html) installed, you can run LazyNet directly without cloning:

```bash
# Run directly from GitHub
nix run 'git+https://github.com/sureshkrishnan-v/lazynet'

# Or with SSH
nix run 'git+ssh://git@github.com/sureshkrishnan-v/lazynet'
```
