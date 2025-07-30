package llm

type Session interface {
        Chat(message string) (string, error)
}

type Server interface {
        Chat(model string) (Session, error)
}
