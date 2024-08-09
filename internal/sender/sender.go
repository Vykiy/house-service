package sender

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Sender struct{}

func New() *Sender {
	return &Sender{}
}

func (s *Sender) SendEmail(ctx context.Context, recipient string, message string) {
	// Имитация отправки сообщения
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	// Имитация неуспешной отправки сообщения
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		fmt.Printf("failed to send message '%s' to '%s'\n", message, recipient)
		return
	}

	fmt.Printf("send message '%s' to '%s'\n", message, recipient)
}
