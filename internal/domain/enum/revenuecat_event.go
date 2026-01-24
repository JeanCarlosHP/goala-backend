package enum

type RevenueCatEventType string

const (
	EventInitialPurchase     RevenueCatEventType = "INITIAL_PURCHASE"
	EventRenewal             RevenueCatEventType = "RENEWAL"
	EventCancellation        RevenueCatEventType = "CANCELLATION"
	EventExpiration          RevenueCatEventType = "EXPIRATION"
	EventBillingIssue        RevenueCatEventType = "BILLING_ISSUE"
	EventProductChange       RevenueCatEventType = "PRODUCT_CHANGE"
	EventNonRenewingPurchase RevenueCatEventType = "NON_RENEWING_PURCHASE"
	EventUncancellation      RevenueCatEventType = "UNCANCELLATION"
)

func (e RevenueCatEventType) String() string {
	return string(e)
}

func (e RevenueCatEventType) IsValid() bool {
	switch e {
	case EventInitialPurchase, EventRenewal, EventCancellation,
		EventExpiration, EventBillingIssue, EventProductChange,
		EventNonRenewingPurchase, EventUncancellation:
		return true
	}
	return false
}
