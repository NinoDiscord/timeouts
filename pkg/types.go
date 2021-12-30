// Copyright (c) 2021 Nino
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pkg

type OperationType int

const (
	Ready OperationType = iota
	Apply
	Request
	RequestAll
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
