package internal

type Experiment struct {
	Responsive bool
	Head       bool
	Leg1       bool
	Leg2       bool
	Leg3       bool
	Leg4       bool
	Leg5       bool
	Leg6       bool
	Wing1      bool
	Wing2      bool
}

func (e Experiment) Extremities() uint8 {
	var count uint8

	if e.Leg1 {
		count += 1
	}

	if e.Leg2 {
		count += 1
	}

	if e.Leg3 {
		count += 1
	}

	if e.Leg4 {
		count += 1
	}

	if e.Leg5 {
		count += 1
	}

	if e.Leg6 {
		count += 1
	}

	if e.Wing1 {
		count += 1
	}

	if e.Wing2 {
		count += 1
	}

	return count
}
