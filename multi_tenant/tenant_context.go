package multi_tenant

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const TenantContextKey contextKey = "tenant_id"

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantContextKey, strings.TrimSpace(tenantID))
}

func TenantIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(TenantContextKey)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func TenantIDFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
}