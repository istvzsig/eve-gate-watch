# Gatewatch

EVE Online system intelligence service written in Go.

## Discord

```text
!dscan <system>
```

Example:

```text
!dscan Tama
```

Returns:

- system, region, and constellation
- security status
- planets, stargates, and stations
- ship jumps in the last hour
- ship kills in the last hour
- pod kills in the last hour
- NPC kills in the last hour

## HTTP API

```text
GET /api/dscan?system=Tama
GET /api/pvp?system=Tama
```

## Data Source

EVE Online ESI public API.

## Start HTTP Server and Discord Bot

```bash
chmod +x ./start.sh
./start.sh
```

## Test HTTP Server

```bash
chmod +x ./test_requests.sh
./test_requests.sh jita
```

CLI:

```bash
curl -sS -X GET "http://localhost:8080/api/dscan?system=<system>"
curl -sS -X GET "http://localhost:8080/api/pvp?system=<system>"
```

The test script queries the Gatewatch HTTP API and returns the system intelligence data for the requested system.
