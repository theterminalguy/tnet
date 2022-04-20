package billing

type PaymentPlan string

const Free PaymentPlan = "free"

// TODO: add more payment plans
// const Premium PaymentPlan = "premium"
// const Enterprise PaymentPlan = "enterprise"

func (PaymentPlan) Values() (kinds []string) {
	for _, k := range []PaymentPlan{Free} {
		kinds = append(kinds, string(k))
	}
	return
}
