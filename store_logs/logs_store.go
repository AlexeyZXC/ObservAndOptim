package store_logs

import (
	"context"
	//"elastic/e"
	"elastic/m"
)

type LogsStore struct {
	E E
}

func NewLogsStore() (LogsStore, error) {
	e, err := NewE("logs")
	if err != nil {
		return LogsStore{}, err
	}
	return LogsStore{E: e}, nil
}

func (s LogsStore) Add(ctx context.Context, logs m.Logs) error {
	return s.E.Insert(ctx, logs)
}

// func (s LogsStore) Search(ctx context.Context, query string) ([]m.Logs, error) {
// 	result, err := s.E.Search(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hits := result.Hits.Hits
// 	logs := []m.Logs{}
// 	for _, hit := range hits {
// 		var log m.Logs
// 		//map[string]interface{} -> struct
// 		err = mapstructure.Decode(hit.Source, &log)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Id = hit.ID
// 		logs = append(logs, log)
// 	}
// 	return logs, nil
// }

// func (s LogsStore) Get(ctx context.Context, id string) (m.Logs, error) {
// 	result, err := s.E.Get(ctx, id)
// 	if err != nil {
// 		return m.Logs{}, err
// 	}
// 	l.L(result)
// 	var log m.Logs

// 	return log, nil
// }
