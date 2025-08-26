````markdown
# go-portscanner 🔎

A fast concurrent TCP port scanner written in [Go](https://go.dev/).  
Built step-by-step while learning Go (goroutines, channels, worker pools, context, benchmarking).

---

## Features
- Scan IP addresses **or hostnames** (DNS resolution supported).
- Customizable **port range** (`-start`, `-end`).
- Adjustable **concurrency** with `-workers`.
- Configurable **timeout** per connection (`-timeout 500ms`, `-timeout 2s`).
- Colored output (thanks to [fatih/color](https://github.com/fatih/color)).

---

## Installation
Clone the repo:
```bash
git clone https://github.com/Borhanxj/go-portscanner.git
cd go-portscanner
````

Build:

```bash
go build ./cmd/portscan
```

Or run directly:

```bash
go run ./cmd/portscan -target 127.0.0.1 -start 20 -end 100
```

---

## Usage

```bash
./portscan -target <host> -start <port> -end <port> -workers <n> -timeout <duration>
```

### Examples

Scan localhost for SSH:

```bash
./portscan -target 127.0.0.1 -start 50 -end 200
```

Scan hostname with 1000 workers:

```bash
./portscan -target scanme.nmap.org -start 1 -end 1024 -workers 1000 -timeout 1s
```

---

## Example Output

```
Scanning target: scanme.nmap.org
Port range: 1-1024
Workers: 1000
Timeout: 1s
Open ports (2):
scanme.nmap.org:22 open
scanme.nmap.org:80 open
```

---

## Disclaimer ⚠️

This tool is for **educational and authorized security testing only**.
Do not scan hosts you do not own or have permission to test.

---

## License

MIT License © 2025 \ Borhan Javadian