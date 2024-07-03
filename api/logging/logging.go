package api

import(
	"github.com/go-kit/log"

)

type loggingMiddleware struct {
	logger log.Logger
	next   HelloWorldService
}
