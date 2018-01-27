package util

import (
	"context"
	"fmt"

	"github.com/op/go-logging"
)

type Logger string
const loggerID = "logger_id"
var log = logging.MustGetLogger("vischain")

func (l Logger) Debug(s string) {
	log.Debugf("[node=%s] %s", l, s)
}
func (l Logger) Debugf(s string, args ...interface{}) {
	log.Debugf("[node=%s] %s", l, fmt.Sprintf(s, args))
}

func (l Logger) Info(s string) {
	log.Infof("[node=%s] %s", l, s)
}
func (l Logger) Infof(s string, args ...interface{}) {
	log.Infof("[node=%s] %s", l, fmt.Sprintf(s, args))
}

func (l Logger) Warning(s string) {
	log.Warningf("[node=%s] %s", l, s)
}
func (l Logger) Warningf(s string, args ...interface{}) {
	log.Warningf("[node=%s] %s", l, fmt.Sprintf(s, args))
}

func (l Logger) Error(s string) {
    log.Errorf("[node=%s] %s", l, s)
}
func (l Logger) Errorf(s string, args ...interface{}) {
	log.Errorf("[node=%s] %s", l, fmt.Sprintf(s, args))
}

func (l Logger) Fatal(s string) {
    log.Fatalf("[node=%s] %s", l, s)
}
func (l Logger) Fatalf(s string, args ...interface{}) {
	log.Fatalf("[node=%s] %s", l, fmt.Sprintf(s, args))
}

func CreateLogger(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, "vischain_node_logger", id)
}

func GetLogger(ctx context.Context) Logger {
    return Logger(ctx.Value("vischain_node_logger").(string))
}