package flowcus

const (
	FORMAT  string  = "2006-01-2 15:04:05 (MST)"
	VERSION float64 = 0.5
)

func NewEploratoryBox() Exploratory {
	return newExploratory()
}

func NewLinearBox() Linear {
	return newLinear()
}

func NewNonLinearBox() NonLinear {
	return newNonLinear()
}
