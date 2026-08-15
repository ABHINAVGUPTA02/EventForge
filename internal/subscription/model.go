package subscription

Type Subscription struct {
    ID string
    TenantID string
    EventTypes []string
    Endpoint string
    MaxConcurrency int
    RateLimit int
}