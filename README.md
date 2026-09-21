# Gatewatch

EVE Online system intelligence service written in Go.

Discord:
!dscan <system>

Example:
!dscan Tama

Returns:

- system / region / constellation
- security status
- planets / stargates / stations
- ship jumps in the last hour
- ship kills in the last hour
- pod kills in the last hour
- NPC kills in the last hour

HTTP:
GET /api/dscan?system=Tama

Data source:
EVE Online ESI public API
