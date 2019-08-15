package logs

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// NewInstanceMsg sets the message to indicate the start of the log
	NewInstanceMsg = "START"
	// EndInstanceMsg sets the message to indicate the end of the log
	EndInstanceMsg = "END"
	// LogLevelDebug defines a normal debug log
	LogLevelDebug = "DEBUG"
	// LogLevelPanic defines a panic log
	LogLevelPanic = "PANIC"
	// LogLevelFatal defines a fatal log
	LogLevelFatal = "FATAL"
	// DateFormat defines the log date format
	DateFormat = time.RFC3339
	// RuneAlpha enumerates Alpha case sensitive runes
	RuneAlpha = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
)

var (
	osExit = os.Exit
)

// Log represents information about a rest server log.
type Log struct {
	entries    []Entry
	folder     string
	identifier string
	showCount  bool
}

// Entry represents information about a rest server log entry.
type Entry struct {
	Level   string
	Message string
	Time    time.Time
}

type JsonEntry struct {
	Timestamp      string `json:"timestamp"`
	Level          string `json:"level"`
	Identifier     string `json:"identifier"`
	SequenceNumber string `json:"sequence_number"`
	Message        string `json:"message"`
}

type JsonEntries []JsonEntry

func (l Log) getDate(t time.Time) string {
	return t.Format(DateFormat)
}

// New creates new instance of Log
func New(folder string) *Log {
	var log Log

	log.folder = folder

	if exist, _ := Exist(log.folder); !exist {
		Mkdir(log.folder)
	}

	log.entries = make([]Entry, 1)
	log.entries[0] = Entry{
		Message: NewInstanceMsg,
		Time:    time.Now(),
	}

	log.Identify()

	return &log
}

// Identify log
func (l *Log) Identify(tags ...string) {
	l.identifier = MD5(fmt.Sprint(time.Now().UnixNano()))
}

// ShowCount log
func (l *Log) ShowCount(show bool) {
	l.showCount = show
}

// GetIdentify log
func (l *Log) GetIdentify() string {
	return l.identifier
}

// GetCount log
func (l *Log) GetCount() bool {
	return l.showCount
}

// SetIdentify log
func (l *Log) SetIdentify(tag string) {
	l.identifier = tag
}

func (l *Log) addEntry(level string, v ...interface{}) {
	l.entries = append(
		l.entries,
		Entry{
			Level:   level,
			Message: fmt.Sprint(v...),
			Time:    time.Now(),
		},
	)
}

// Entries returns all the entries
func (l *Log) Entries() []Entry {
	return l.entries
}

// Print a regular log
func (l *Log) Print(v ...interface{}) {
	l.addEntry(LogLevelDebug, v...)
}

// Panic then throws a panic with the same message afterwards
func (l *Log) Panic(v ...interface{}) {
	l.addEntry(LogLevelPanic, v...)
	panic(fmt.Sprint(v...))
}

// ThrowFatalTest allows Fatal to be testable
func (l *Log) ThrowFatalTest(msg string) {
	defer func() { osExit = os.Exit }()
	osExit = func(int) {}
	l.Fatal(msg)
}

// Fatal is equivalent to Print() and followed by a call to os.Exit(1)
func (l *Log) Fatal(v ...interface{}) {
	l.addEntry(LogLevelFatal, v...)
	l.Dump()
	osExit(1)
}

// LastEntry returns the last inserted log
func (l *Log) LastEntry() Entry {
	return l.entries[len(l.entries)-1]
}

// Count returns number of inserted logs
func (l *Log) Count() int {
	return len(l.entries)
}

// Dump will print all the messages to the io.
func (l *Log) Dump() {
	var (
		filename, line, lines string
	)

	l.addEntry("", EndInstanceMsg)
	len := len(l.entries)

	for i := 0; i < len; i++ {
		if l.identifier != "" {
			if l.showCount {
				line = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", l.getDate(l.entries[i].Time), l.entries[i].Level, l.identifier, strconv.Itoa(i), l.entries[i].Message)
			} else {
				line = fmt.Sprintf("%s\t%s\t%s\t%s\n", l.getDate(l.entries[i].Time), l.entries[i].Level, l.identifier, l.entries[i].Message)
			}
		} else {
			if l.showCount {
				line = fmt.Sprintf("%s\t%s\t%s\t%s\n", l.getDate(l.entries[i].Time), l.entries[i].Level, strconv.Itoa(i), l.entries[i].Message)
			} else {
				line = fmt.Sprintf("%s\t%s\t%s\n", l.getDate(l.entries[i].Time), l.entries[i].Level, l.entries[i].Message)
			}
		}
		lines = lines + line

		fmt.Print(line)
	}

	go func() {

		if l.identifier != "" {
			filename = fmt.Sprintf("%s/%d.%s.log", l.folder, time.Now().UnixNano(), l.identifier)
		} else {
			filename = fmt.Sprintf("%s/%d.%s.log", l.folder, time.Now().UnixNano(), strings.ToLower(RandomString(20, []rune(RuneAlpha))))
		}

		file, _ := os.Create(filename)
		defer file.Close()
		l.entries = make([]Entry, 0)
		file.WriteString(lines)
	}()
}

// Exist checks if folder or file exist
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, err
}

// Mkdir make directory if directory does not exist
func Mkdir(path string) error {
	if exist, _ := Exist(path); !exist {
		return os.Mkdir(path, 0777)
	}
	return nil
}

// Delete delete file
func Delete(filename string) error {
	return os.Remove(filename)
}

func RandomString(max int, runes []rune) string {
	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]rune, max)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}

// MD5 calculates the MD5 hash of a string
func MD5(raw string) string {
	m := md5.New()
	io.WriteString(m, raw)
	return fmt.Sprintf("%x", m.Sum(nil))
}
