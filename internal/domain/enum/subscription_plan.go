package enum

type SubscriptionPlan string

const (
	PlanFree    SubscriptionPlan = "free"
	PlanMonthly SubscriptionPlan = "monthly"
	PlanYearly  SubscriptionPlan = "yearly"
)

func (p SubscriptionPlan) String() string {
	return string(p)
}

func (p SubscriptionPlan) IsValid() bool {
	switch p {
	case PlanFree, PlanMonthly, PlanYearly:
		return true
	}
	return false
}
