package multi_tenant

import (
	"context"
	"testing"
)

func TestRegisterAndLookupTenant(t *testing.T) {
	m := NewManager()

	err := m.RegisterTenant(Tenant{
		ID:     "tenant_acme",
		Name:   "Acme Corp",
		Status: TenantStatusActive,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := m.GetTenant("tenant_acme")
	if !ok {
		t.Fatal("expected tenant to exist")
	}
	if got.Name != "Acme Corp" {
		t.Fatalf("expected Acme Corp, got %s", got.Name)
	}
}

func TestIsAllowed(t *testing.T) {
	m := NewManager()

	_ = m.RegisterTenant(Tenant{
		ID:     "tenant_a",
		Name:   "Tenant A",
		Status: TenantStatusActive,
	})
	_ = m.RegisterTenant(Tenant{
		ID:     "tenant_b",
		Name:   "Tenant B",
		Status: TenantStatusDisabled,
	})

	if !m.IsAllowed("tenant_a") {
		t.Fatal("expected tenant_a to be allowed")
	}
	if m.IsAllowed("tenant_b") {
		t.Fatal("expected tenant_b to be blocked")
	}
	if m.IsAllowed("missing") {
		t.Fatal("expected missing tenant to be blocked")
	}
}

func TestTenantContext(t *testing.T) {
	ctx := WithTenantID(context.Background(), "tenant_demo")
	got := TenantIDFromContext(ctx)
	if got != "tenant_demo" {
		t.Fatalf("expected tenant_demo, got %s", got)
	}
}