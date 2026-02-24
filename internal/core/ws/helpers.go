package ws

import "encoding/json"

// JSON serializa un objeto a RawMessage para usar en Envelope.Data.
// (Si hay error, devuelve {} para no romper el flujo de WS).
func JSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}
