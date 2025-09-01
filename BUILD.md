# Commands

## Build with versions

go build -ldflags "-X main.Version=0.5.3 -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o crush .
