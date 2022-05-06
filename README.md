# Devcon
A simple package to ssh into devices.

## TODO
- [x] Pull the client bits out of RunCommand and pull that logic into the Factory function
- [ ] Test against a test server and not a lab SRX

## Profile
- go test -v -run Run -cpuprofile cpu.prof -memprofile mem.prof -bench .

## Testing
### Unit tests
#### All tests
go test -v
#### A specific test
go test -v -run RunCommand
#### Specific tests matching a pattern
go test -v -run Run
### Benchmarks
go test -run RunCommand -bench=.