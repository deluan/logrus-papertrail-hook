package logrus_papertrail

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	format = "Jan 2 15:04:05"
)

const (
	ConnTCP = "tcp"
	ConnUDP = "udp"
)

// private interface to make tests with TCP conn
type conn_int interface {
	Write([]byte) (int, error)
}

// PapertrailHook to send logs to a logging service compatible with the Papertrail API.
type Hook struct {
	// Connection Details
	Host     string
	Port     int
	ConnType string

	// App Details
	Appname  string
	Hostname string

	levels []logrus.Level

	conn conn_int
}

// NewPapertrailHook creates a UDP hook to be added to an instance of logger.
func NewPapertrailHook(hook *Hook) (*Hook, error) {
	if hook.ConnType == ConnTCP {
		return newPapertrailTCPHook(hook)
	}
	var err error
	hook.ConnType = ConnUDP
	hook.conn, err = net.Dial(hook.ConnType, fmt.Sprintf("%s:%d", hook.Host, hook.Port))
	return hook, err
}

// NewPapertrailTCPHook creates a TCP/TLS hook to be added to an instance of logger.
// Deprecated. Use NewPapertrailHook with hook.ConnType = ConnTCP
func NewPapertrailTCPHook(hook *Hook) (*Hook, error) {
	return newPapertrailTCPHook(hook)
}

func newPapertrailTCPHook(hook *Hook) (*Hook, error) {
	var err error
	hook.ConnType = ConnTCP
	hook.conn, err = tls.Dial(hook.ConnType, fmt.Sprintf("%s:%d", hook.Host, hook.Port), nil)
	return hook, err
}

// Fire is called when a log event is fired.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	date := time.Now().Format(format)
	msg, _ := entry.String()
	payload := fmt.Sprintf("<22> %s %s %s: %s", date, hook.Hostname, hook.Appname, msg)

	bytesWritten, err := hook.conn.Write([]byte(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to send log line to Papertrail via %s. Wrote %d bytes before error: %v", hook.ConnType, bytesWritten, err)
		return err
	}

	return nil
}

// SetLevels specify nessesary levels for this hook.
func (hook *Hook) SetLevels(lvs []logrus.Level) {
	hook.levels = lvs
}

// Levels returns the available logging levels.
func (hook *Hook) Levels() []logrus.Level {

	if hook.levels == nil {
		return []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		}
	}

	return hook.levels
}
