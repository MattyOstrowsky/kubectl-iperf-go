1. `go install`
2. `export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))`
3. `kubectl-iperf --help`
4. `kubectl plugin list`
5. `kubectl iperf --help`