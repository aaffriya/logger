package file

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
)

type fileHandler struct {
	file  *os.File
	writer io.Writer
	mu     *sync.Mutex
}

type FileHandler interface {
	Handle(r slog.Record) error
}

func NewFileHandler(w io.Writer, file *os.File) FileHandler {
	return &fileHandler{
		file:   file,
		writer: w,
		mu:     &sync.Mutex{},
	}
}

func (h *fileHandler) Handle(r slog.Record) error {
	
	logData := map[string]any{
		"timestamp": r.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		"level":     r.Level.String(),
		"message":   r.Message,
	}
	
	r.Attrs(func(a slog.Attr) bool {
		logData[a.Key] = a.Value.Any()
		return true
	})
	
	jsonBytes, err := json.Marshal(logData)
	if err != nil {
		return err
	}
	
	jsonBytes = append(jsonBytes, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.writer.Write(jsonBytes)
	
	return err
}
