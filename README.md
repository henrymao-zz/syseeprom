# platform

Go multi-vendor EEPROM hardware integration using consumer-defined interfaces and the registry/factory pattern.

## Structure

```
cmd/syseeprom/main.go           # Entrypoint — auto-detects vendor from /etc/machine.conf
eeprom/
  eeprom.go                    # EEPROM interface
  registry.go                  # Register() / GetDriver() factory
sysconf/
  sysconf.go                   # Parses /etc/machine.conf, maps onie_machine → driver
  sysconf_test.go              # Unit tests
vendors/
  broadcom/
    broadcom.go                # Broadcom ONIE EEPROM driver
  dell/s5232f/
    s5232f.go                  # Dell S5232f ONIE TLV EEPROM driver
    s5232f_test.go             # Unit tests
  mellanox/
    mellanox.go                # Mellanox I2C EEPROM driver
```

## Interface

```go
type EEPROM interface {
    GetBaseMAC() (string, error)
    GetSerialNumber() (string, error)
    GetModel() (string, error)
}
```

## How it works

1. Vendor packages implement the `eeprom.EEPROM` interface implicitly.
2. Each vendor calls `eeprom.Register()` in its `init()` function at startup.
3. `main()` reads `onie_machine` from `/etc/machine.conf` and resolves the driver.
4. The application interacts purely with the interface — no vendor-specific code.

### Driver resolution

| Priority | Source | Example |
|---|---|---|
| 1 | CLI argument | `./syseeprom dell-s5232f` |
| 2 | `/etc/machine.conf` → `onie_machine` | `onie_machine=dellemc_s5232f_c3538` |
| 3 | Fallback | `mellanox` |

### machine.conf → driver mapping

| onie_machine | Driver |
|---|---|
| `dellemc_s5232f_c3538` | `dell-s5232f` |
| `dellemc_s5200_c3538` | `dell-s5232f` |

## Supported platforms

| Vendor | Platform | Driver | EEPROM Format |
|---|---|---|---|
| Dell | S5232f-ON | `dell-s5232f` | ONIE TLV (TlvInfo) |
| Mellanox | SN2700 | `mellanox` | I2C EEPROM |
| Broadcom | Generic | `broadcom` | ONIE EEPROM |

### ONIE TLV Format (Dell S5232f)

The Dell S5232f driver parses the ONIE TlvInfo EEPROM format:

```
Offset  Size  Field
0       8     ID string ("TlvInfo\x00")
8       1     Version (0x01)
9       2     Total TLV length (big-endian uint16)
11      N     TLV entries
```

Each TLV entry:
```
Offset  Size  Field
0       1     Type code
1       1     Length (N)
2       N     Value
```

TLV type codes used:

| Code | Field |
|---|---|
| `0x21` | Product Name (model) |
| `0x22` | Part Number |
| `0x23` | Serial Number |
| `0x24` | Base MAC Address (6 bytes) |
| `0x26` | Device Version |
| `0x2F` | Service Tag |
| `0xFE` | CRC-32 |

## Build

```bash
go build -o syseeprom ./cmd/syseeprom/
```

## Test

```bash
go test ./...
```

## Run

```bash
./syseeprom                  # auto-detect from /etc/machine.conf
./syseeprom dell-s5232f      # explicit driver
./syseeprom broadcom
./syseeprom mellanox
```

## Adding a vendor

1. Create `vendors/<vendor>/<platform>/<platform>.go`
2. Implement `eeprom.EEPROM` (`GetBaseMAC`, `GetSerialNumber`, `GetModel`)
3. Register in `init()` with `eeprom.Register("<name>", factory)`
4. Add a blank import in `cmd/syseeprom/main.go`:
   ```go
   _ "platform/vendors/<vendor>/<platform>"
   ```
 5. Map `onie_machine` → driver name in `sysconf/sysconf.go`

## Debian Package

```bash
# Copy packaging files to project root and build
cp -r package/debian .
dpkg-buildpackage -us -uc -b

# Source-only build for PPA submission
dpkg-buildpackage -S -sa
dput ppa:henrymao/ubuntu-nos ../syseeprom_*_source.changes
```

Package files are maintained in `package/debian/` to keep them separate
from the Go source tree. When building, copy them to the project root
where `dpkg-buildpackage` expects them.
