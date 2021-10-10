package main

type DataStore interface {
	getResult() (Result, error)
}

type Result struct {
	RemainedResponsivePercent string
}

type MockDataStore struct{}

func (store MockDataStore) getResult() (Result, error) {
	result := Result{
		RemainedResponsivePercent: "13.59",
	}

	return result, nil
}
