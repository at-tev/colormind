package errors

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

// internalErr is an error message for output to users
const InternalErr = "Internal Error"

func HandleError(msg string, err error) {
	sentry.CaptureException(err)
	fmt.Printf("%s: %s", msg, err.Error())
}
