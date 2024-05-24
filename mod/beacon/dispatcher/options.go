package dispatcher

type DispatcherOption func(*Dispatcher) error

func WithHandler(eventType EventType, handler Handler) DispatcherOption {
	return func(d *Dispatcher) error {
		d.mu.Lock()
		defer d.mu.Unlock()

		// TODO: Do we want multiple handlers for the same event type?
		if _, ok := d.registry[eventType]; ok {
			return ErrHandlerAlreadyExists
		}

		d.registry[eventType] = handler
		return nil
	}
}
