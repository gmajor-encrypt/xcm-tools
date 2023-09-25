package tracker

import (
	"encoding/json"
	"errors"
)

type Outcome struct {
	Complete   *interface{} `json:"complete,omitempty"`
	Incomplete *struct {
		Error Enum `json:"Error"`
	} `json:"Incomplete"`
	Error *Enum `json:"error"`
}

func ParseAttempted(event Event, index int) (bool, error) {
	var outcome Outcome
	bytes, err := json.Marshal(event.Params[index].Value)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bytes, &outcome)
	if err != nil {
		return false, err
	}
	if outcome.Complete != nil {
		return true, nil
	}
	if outcome.Incomplete != nil {
		return false, errors.New(outcome.Incomplete.Error.Key())
	}
	if outcome.Error != nil {
		return false, errors.New(outcome.Error.Key())
	}
	return false, nil
}

type Enum map[string]interface{}

func (e Enum) Value() interface{} {
	for _, v := range e {
		return v
	}
	return ""
}

func (e Enum) Key() string {
	for k := range e {
		return k
	}
	return ""
}
