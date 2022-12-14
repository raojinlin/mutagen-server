package prompting

import (
	"log"
)

// NoopPrompter empty prompter
type NoopPrompter struct {
}

func (n *NoopPrompter) Message(message string) error {
	log.Println("prompter Message", message)
	return nil
}

func (n *NoopPrompter) Prompt(message string) (string, error) {
	log.Println("prompter Prompt", message)
	return "", nil
}
