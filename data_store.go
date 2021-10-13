package main

type DataStore interface {
	getResult() (Result, error)
	storeExperiment(Experiment) error
}

type Result struct {
	RemainedResponsivePercent         string
	RemainedResponsiveHeadlessPercent string
	AverageExtremitiesRemoved         string
	RemainedResponsive1MissingPercent string
	RemainedResponsive2MissingPercent string
	RemainedResponsive3MissingPercent string
	RemainedResponsive4MissingPercent string
	RemainedResponsive5MissingPercent string
	RemainedResponsive6MissingPercent string
	RemainedResponsive7MissingPercent string
	RemainedResponsive8MissingPercent string
}

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

type ExperimentForm struct {
	Responsive string
	Head       string
	Leg1       string
	Leg2       string
	Leg3       string
	Leg4       string
	Leg5       string
	Leg6       string
	Wing1      string
	Wing2      string
}

type SuccessStore struct{}

func (store SuccessStore) getResult() (Result, error) {
	result := Result{
		RemainedResponsivePercent:         "13.59",
		RemainedResponsiveHeadlessPercent: "6.31",
		AverageExtremitiesRemoved:         "3.72",
		RemainedResponsive1MissingPercent: "95.88",
		RemainedResponsive2MissingPercent: "85.72",
		RemainedResponsive3MissingPercent: "75.61",
		RemainedResponsive4MissingPercent: "65.91",
		RemainedResponsive5MissingPercent: "35.82",
		RemainedResponsive6MissingPercent: "24.27",
		RemainedResponsive7MissingPercent: "11.56",
		RemainedResponsive8MissingPercent: "0.03",
	}

	return result, nil
}

func (store SuccessStore) storeExperiment(e Experiment) error {
	return nil
}
