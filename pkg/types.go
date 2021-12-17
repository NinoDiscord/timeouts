package pkg

type OperationType int

const (
	Ready OperationType = iota
	Apply
	Request
	RequestAll
	RequestAllBack
)

type Message struct {
	OP   OperationType `json:"op"`
	Data interface{}   `json:"d"`
}

type Timeout struct {
	Type        string `json:"type"`
	GuildId     string `json:"guild_id"`
	UserId      string `json:"user_id"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresAt   int64  `json:"expires_at"`
	ModeratorId string `json:"moderator_id"`
	Reason      string `json:"reason,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
