# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	equipmentlending/cmd/equipment-service	0.025s
ok  	equipmentlending/internal/audit	0.002s
ok  	equipmentlending/internal/config	0.001s
--- FAIL: TestEquipmentSortKeepsAll (0.00s)
    sort_test.go:25: date sort lost records: got 2 want 3
FAIL
FAIL	equipmentlending/internal/equipment	0.001s
ok  	equipmentlending/internal/persistence	0.037s
ok  	equipmentlending/internal/reporting	0.001s
ok  	equipmentlending/internal/service	0.044s
ok  	equipmentlending/internal/transport	0.021s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/equipment-service): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/equipment-service): exit `0`
