package multi_tenant

import (
	"fmt"
	"strings"
)

type Manager struct {
	tenants map[string]Tenant
}

func NewManager() *Manager {
	return &Manager{
		tenants: map[string]Tenant{},
	}
}

func (m *Manager) RegisterTenant(t Tenant) error {
	t.ID = strings.TrimSpace(t.ID)
	t.Name = strings.TrimSpace(t.Name)

	if t.ID == "" {
		return fmt.Errorf("tenant id is required")
	}
	if t.Name == "" {
		return fmt.Errorf("tenant name is required")
	}
	if t.Status == "" {
		t.Status = TenantStatusActive
	}

	m.tenants[t.ID] = t
	return nil
}

func (m *Manager) GetTenant(id string) (Tenant, bool) {
	t, ok := m.tenants[strings.TrimSpace(id)]
	return t, ok
}

func (m *Manager) IsAllowed(id string) bool {
	t, ok := m.GetTenant(id)
	if !ok {
		return false
	}
	return t.Status == TenantStatusActive
}