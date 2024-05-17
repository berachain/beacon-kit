package events

// EventType is the type that defines the type of event.
type EventType int

// Data is the event that is sent with operation feed updates.
type Data[T any] struct {
	// Type is the type of event.
	Type EventType
	// Data is event-specific data.
	Data T
}
