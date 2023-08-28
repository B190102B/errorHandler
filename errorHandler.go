package errorHandler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

var (
	logFilename string
	sentryDns   string
)

func init() {
	logFilename = os.Getenv("ERROR_LOG_FILE")
	sentryDns = os.Getenv("SENTRY_DSN")
}

func HandleError(err error, desc ...any) {
	if err != nil {
		msg := fmt.Sprintln(desc, ": ", err)

		ThrowError(msg)
	}
}

func ThrowError(msg any) {
	panic(msg)
}

func RegisterRecover() {
	// run defer on top
	if err := recover(); err != nil {
		RegisterSaveLogFile(err)
		RegisterSentry(err)
	}
}

func RegisterSaveLogFile(err any) {
	if logFilename != "" {
		log.Println("ERROR_LOG_FILE is require !!")
		return
	}

	f, _ := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("%v\n", err))

	log.Println("Success Save Log To: ", logFilename)
}

func RegisterSentry(err any) {
	if sentryDns != "" {
		log.Println("SENTRY_DSN is require !!")
		return
	}

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:       sentryDns,
		Transport: sentrySyncTransport,
	}); err != nil {
		log.Printf("%v", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureException(fmt.Errorf("%v", err))

	log.Println("Success Send Sentry")
}
