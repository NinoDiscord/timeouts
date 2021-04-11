package pkg

type OPCode int
type PunishmentType string

const (
	Ready OPCode = iota
	Apply
	Request
	Acknowledged
)

type WSMessage struct {
	Op OPCode      `json:"op"`
	D  interface{} `json:"d"`
}

type Timeout struct {
	Type    PunishmentType `json:"type"`
	Guild   string         `json:"guild"`
	User    string         `json:"user"`
	Issued  int64          `json:"issued"`
	Expired int64          `json:"expired"`
}
