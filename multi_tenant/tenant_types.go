package multi_tenant

type Tenant struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Status   string                 `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

const (
	TenantStatusActive   = "active"
	TenantStatusDisabled = "disabled"
)