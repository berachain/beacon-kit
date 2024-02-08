package queue

import "errors"

var ErrAddToStoppedQueue = errors.New("Cannot add to stopped queue")
