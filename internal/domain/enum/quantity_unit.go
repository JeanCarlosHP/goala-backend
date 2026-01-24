package enum

type QuantityUnit string

const (
	QuantityUnitGram    QuantityUnit = "g"
	QuantityUnitML      QuantityUnit = "ml"
	QuantityUnitServing QuantityUnit = "serving"
)

func (q QuantityUnit) IsValid() bool {
	switch q {
	case QuantityUnitGram, QuantityUnitML, QuantityUnitServing:
		return true
	}
	return false
}

func (q QuantityUnit) String() string {
	return string(q)
}
