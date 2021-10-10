package main

type DataStore interface {
	getResult() (Result, error)
}

type Result struct {
	remainedResponsivePercent float64
}

type MockDataStore struct{}

func (store MockDataStore) getResult() (Result, error) {
	result := Result{
		remainedResponsivePercent: 13.5,
	}

	return result, nil
}
