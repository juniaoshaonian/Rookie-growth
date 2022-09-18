package opentelmetry

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"webframe"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Builder() webframe.Middleware {
	return func(next webframe.HanleFunc) webframe.HanleFunc {
		return func(ctx *webframe.Context) {
			spanCtx, span := m.Tracer.Start(ctx.Req.Context(), "Unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			ctx.Req = ctx.Req.WithContext(spanCtx)
			next(ctx)

			defer func() {
				if ctx.MatchRouter != "" {
					// 不是 404
					span.SetName(ctx.MatchRouter)
				}
				span.SetAttributes(attribute.Int("http.status", ctx.ResponseCode))

			}()
		}
	}
}
