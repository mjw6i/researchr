package internal

import (
	"errors"
)

type SuccessStore struct{}

func (store SuccessStore) getResult() (Result, error) {
	result := Result{
		RemainedResponsivePercent:         "57.03",
		RemainedResponsiveHeadlessPercent: "1.63",
		AverageExtremitiesRemoved:         "3.72",
		RemainedResponsiveMissingPercent: [9]string{
			"100",
			"95.88",
			"87.80",
			"75.61",
			"65.91",
			"34.82",
			"24.27",
			"11.56",
			"0.03",
		},
	}

	return result, nil
}

func (store SuccessStore) storeExperiment(e Experiment) error {
	return nil
}

type FailureStore struct{}

func (store FailureStore) getResult() (Result, error) {
	return Result{}, errors.New("Store error")
}

func (store FailureStore) storeExperiment(e Experiment) error {
	return errors.New("Store error")
}
