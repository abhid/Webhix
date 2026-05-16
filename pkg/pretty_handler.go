package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"gh.tarampamp.am/colors"
)

const logTimeLayout = "15:04:05.000"

var (
	logTimeStyle    = colors.FgWhite | colors.Faint
	logMessageStyle = colors.FgWhite | colors.Bold
	logSourceStyle  = colors.FgWhite | colors.Faint
	logKeyStyle     = colors.FgCyan
	logErrorStyle   = colors.FgRed | colors.Bold
)

type prettyHandler struct {
	out       io.Writer
	level     slog.Leveler
	addSource bool
	attrs     []prettyAttr
	groups    []string
	cwd       string
	mu        *sync.Mutex
}

type prettyAttr struct {
	groups []string
	attr   slog.Attr
}

func newPrettyHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {
	if out == nil {
		out = io.Discard
	}

	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	return &prettyHandler{
		out:       out,
		level:     opts.Level,
		addSource: opts.AddSource,
		cwd:       cwd,
		mu:        &sync.Mutex{},
	}
}

func (h *prettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.level != nil {
		minLevel = h.level.Level()
	}

	return level >= minLevel
}

func (h *prettyHandler) Handle(_ context.Context, record slog.Record) error {
	var buf bytes.Buffer

	if !record.Time.IsZero() {
		buf.WriteString(logTimeStyle.Wrap(record.Time.Format(logTimeLayout)))
		buf.WriteByte(' ')
	}

	buf.WriteString(formatLogLevel(record.Level))

	if record.Message != "" {
		buf.WriteByte(' ')
		buf.WriteString(logMessageStyle.Wrap(record.Message))
	}

	if h.addSource && record.PC != 0 {
		buf.WriteByte(' ')
		buf.WriteString(logSourceStyle.Wrap("source=" + h.source(record.PC)))
	}

	for _, attr := range h.attrs {
		h.appendAttr(&buf, attr.groups, attr.attr)
	}

	record.Attrs(func(attr slog.Attr) bool {
		h.appendAttr(&buf, h.groups, attr)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h.clone()
	for _, attr := range attrs {
		clone.attrs = append(clone.attrs, prettyAttr{
			groups: slicesClone(h.groups),
			attr:   attr,
		})
	}

	return clone
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	clone := h.clone()
	clone.groups = append(clone.groups, name)

	return clone
}

func (h *prettyHandler) clone() *prettyHandler {
	return &prettyHandler{
		out:       h.out,
		level:     h.level,
		addSource: h.addSource,
		attrs:     append([]prettyAttr(nil), h.attrs...),
		groups:    slicesClone(h.groups),
		cwd:       h.cwd,
		mu:        h.mu,
	}
}

func (h *prettyHandler) appendAttr(buf *bytes.Buffer, groups []string, attr slog.Attr) {
	if attr.Equal(slog.Attr{}) {
		return
	}

	attr.Value = attr.Value.Resolve()

	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		if len(groupAttrs) == 0 {
			return
		}

		if attr.Key != "" {
			groups = append(slicesClone(groups), attr.Key)
		}

		for _, groupAttr := range groupAttrs {
			h.appendAttr(buf, groups, groupAttr)
		}

		return
	}

	key := joinLogKey(groups, attr.Key)
	if key == "" {
		return
	}

	value := formatLogValue(attr.Value)

	buf.WriteByte(' ')
	buf.WriteString(logKeyStyle.Wrap(key))
	buf.WriteByte('=')

	if key == "error" || key == "err" {
		buf.WriteString(logErrorStyle.Wrap(value))
		return
	}

	buf.WriteString(value)
}

func (h *prettyHandler) source(pc uintptr) string {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()

	file := frame.File
	if h.cwd != "" {
		if rel, err := filepath.Rel(h.cwd, file); err == nil && !strings.HasPrefix(rel, "..") {
			file = rel
		}
	}

	return fmt.Sprintf("%s:%d", filepath.ToSlash(file), frame.Line)
}

func formatLogLevel(level slog.Level) string {
	label := strings.ToUpper(level.String())
	if len(label) < 5 {
		label += strings.Repeat(" ", 5-len(label))
	}

	switch {
	case level < slog.LevelInfo:
		return (colors.FgMagenta | colors.Bold).Wrap(label)
	case level < slog.LevelWarn:
		return (colors.FgGreen | colors.Bold).Wrap(label)
	case level < slog.LevelError:
		return (colors.FgYellow | colors.Bold).Wrap(label)
	default:
		return (colors.FgRed | colors.Bold).Wrap(label)
	}
}

func formatLogValue(value slog.Value) string {
	value = value.Resolve()

	switch value.Kind() {
	case slog.KindString:
		return quoteLogString(value.String())
	case slog.KindBool:
		return strconv.FormatBool(value.Bool())
	case slog.KindDuration:
		return value.Duration().String()
	case slog.KindFloat64:
		return strconv.FormatFloat(value.Float64(), 'f', -1, 64)
	case slog.KindInt64:
		return strconv.FormatInt(value.Int64(), 10)
	case slog.KindTime:
		return value.Time().Format(time.RFC3339Nano)
	case slog.KindUint64:
		return strconv.FormatUint(value.Uint64(), 10)
	case slog.KindAny:
		if err, ok := value.Any().(error); ok {
			return quoteLogString(err.Error())
		}

		return quoteLogString(fmt.Sprintf("%+v", value.Any()))
	case slog.KindGroup:
		return quoteLogString(fmt.Sprintf("%+v", value.Group()))
	case slog.KindLogValuer:
		return quoteLogString(fmt.Sprintf("%+v", value.Any()))
	default:
		return quoteLogString(value.String())
	}
}

func quoteLogString(value string) string {
	if value == "" || strings.ContainsAny(value, " \t\r\n=\"") {
		return strconv.Quote(value)
	}

	return value
}

func joinLogKey(groups []string, key string) string {
	if len(groups) == 0 {
		return key
	}

	if key == "" {
		return strings.Join(groups, ".")
	}

	parts := append(slicesClone(groups), key)

	return strings.Join(parts, ".")
}

func slicesClone[T any](items []T) []T {
	if len(items) == 0 {
		return nil
	}

	return append([]T(nil), items...)
}
