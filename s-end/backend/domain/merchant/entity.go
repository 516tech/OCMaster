package merchant

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
)

type Merchant struct {
	ID           uint64
	Phone        string
	PasswordHash string
	Status       Status
	RiskTemplate string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *Merchant) IsActive() bool {
	return m.Status == StatusActive
}
