package progress

import "io"

type writeCounter struct {
	io.Writer
	task *Task
}

func (w *writeCounter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)

	w.task.Advance(int64(n))

	return n, err
}

type readCounter struct {
	io.Reader
	task *Task
}

func (w *readCounter) Read(p []byte) (n int, err error) {
	n, err = w.Reader.Read(p)

	w.task.Advance(int64(n))

	return n, err
}

// WriterCounter wraps w with an io.Writer that advances t by the bytes written
// on every call to Write.
func (t *Task) WriteCounter(w io.Writer) io.Writer {
	return &writeCounter{Writer: w, task: t}
}

// ReadCounter wraps r with an io.Reader that advances t by the bytes read on
// every call to Read.
func (t *Task) ReadCounter(r io.Reader) io.Reader {
	return &readCounter{Reader: r, task: t}
}
