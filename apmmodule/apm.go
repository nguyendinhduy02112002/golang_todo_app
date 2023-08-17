package apmmodule

import (
	"go.elastic.co/apm"
)

// InitAPM khởi tạo và cấu hình tracer của Elastic APM
func InitAPM(serviceName, serviceVersion string) *apm.Tracer {
	tracer, err := apm.NewTracer(serviceName, serviceVersion)
	if err != nil {
		panic(err)
	}
	return tracer
}
