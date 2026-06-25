package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// isTerminal reports whether w is a character device (a TTY). Non-file
// writers (a bytes.Buffer in tests, a pipe, a regular file) return false so
// the rewriting progress line never pollutes piped or redirected output with
// carriage returns.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// progressReporter renders a single rewriting status line to a TTY. When the
// destination is not a terminal (or --quiet is set) every method is a no-op,
// so CI logs and piped output stay clean.
//
// It is NOT safe for concurrent use. Callers drive it from a single
// goroutine — the scanners invoke their Progress callback from one dedicated
// reporter goroutine, and the fingerprint loop is sequential.
type progressReporter struct {
	w       io.Writer
	label   string
	enabled bool
	lastLen int
}

func newProgressReporter(w io.Writer, label string, quiet bool) *progressReporter {
	return &progressReporter{
		w:       w,
		label:   label,
		enabled: !quiet && isTerminal(w),
	}
}

// update redraws the progress line with the current done/total counts. The
// line is padded to overwrite any longer previous render before the carriage
// return, so a shrinking count never leaves stale characters behind.
func (p *progressReporter) update(done, total int) {
	if !p.enabled || total <= 0 {
		return
	}
	pct := done * 100 / total
	line := fmt.Sprintf("%s %d/%d (%d%%)", p.label, done, total, pct)
	pad := ""
	if p.lastLen > len(line) {
		pad = strings.Repeat(" ", p.lastLen-len(line))
	}
	_, _ = fmt.Fprintf(p.w, "\r%s%s", line, pad)
	p.lastLen = len(line)
}

// clear erases the current progress line so a following summary prints on a
// clean line. No-op when disabled or when nothing has been drawn.
func (p *progressReporter) clear() {
	if !p.enabled || p.lastLen == 0 {
		return
	}
	_, _ = fmt.Fprintf(p.w, "\r%s\r", strings.Repeat(" ", p.lastLen))
	p.lastLen = 0
}
