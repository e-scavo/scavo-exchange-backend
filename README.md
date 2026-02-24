# SCAVO Exchange Backend

## Run
  SCAVO_ENV=local SCAVO_HTTP_ADDR=:8080 go run ./cmd/scavo-server

## HTTP
  GET /health
  GET /version
  GET /ws (websocket)

## WS protocol
Envelope: {"id":"...","type":"req|res|evt","action":"system.ping", "data":{...}}

Example ping:
  {"id":"1","type":"req","action":"system.ping"}

## Test WS (websocat)
  websocat ws://localhost:8080/ws
  {"id":"1","type":"req","action":"system.ping"}
