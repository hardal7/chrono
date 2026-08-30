package requestctx

type ContextKey string

const (
	UserID    ContextKey = "userID"
	SessionID ContextKey = "sessionID"
	RequestID ContextKey = "requestID"
	IP        ContextKey = "IP"
)
