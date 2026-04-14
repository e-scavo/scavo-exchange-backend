// Package authorization defines the static authorization model introduced in
// Phase 0.10.1. It intentionally stops at role/permission primitives and their
// foundational mapping so later subphases can wire context propagation,
// policy evaluation and endpoint enforcement without changing the core model.
package authorization
