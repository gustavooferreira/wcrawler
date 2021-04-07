// Package ring provides an implementation of a ring buffer containing strings.
package ring

// Buffer represents a ring buffer.
// If the buffer is full, the last item in the buffer is replaced with the new item.
// Can only add entries to the buffer but never remove them.
type Buffer struct {
	buf          []string
	readPointer  int
	writePointer int
	// len represents the number of items in the buffer
	len int
	// size represents the total size of the buffer
	size int
}

// New returns a new RingBuffer.
func New(size int) Buffer {
	return Buffer{buf: make([]string, size), size: size}
}

// Len returns the length of the buffer (number of entries filled).
func (rb *Buffer) Len() int {
	return rb.len
}

// Add adds a new entry to the buffer.
func (rb *Buffer) Add(entry string) {
	if rb.len == rb.size {
		rb.readPointer++
		if rb.readPointer >= rb.size {
			rb.readPointer %= rb.size
		}
	} else {
		rb.len++
	}

	rb.buf[rb.writePointer] = entry
	rb.writePointer++
	if rb.writePointer >= rb.size {
		rb.writePointer %= rb.size
	}
}

// ReadAll returns a slice of strings with all the elements in the buffer
func (rb *Buffer) ReadAll() []string {
	if rb.len == 0 {
		return nil
	}

	entries := make([]string, rb.len)

	for i := 0; i < rb.len; i++ {
		index := rb.readPointer + i

		if index >= rb.size {
			index %= rb.size
		}

		entries[i] = rb.buf[index]
	}

	return entries
}
