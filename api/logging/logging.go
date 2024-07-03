package api

import(
	"github.com/go-kit/log"
    miniotvt "github.com/chuminh2001100/goPractice/minio"
)

type loggingMiddleware struct {
	logger log.Logger
	next   miniotvt.Service
}
