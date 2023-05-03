package e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type E struct {
	C                    *elasticsearch.Client
	IndexName            string
	FullTextSearchColumn string
}

type I interface{}
type M map[string]I

// # Changed password for user elastic
// # PASSWORD elastic = ZvVx8w0v6CCtOFHK0477

func NewE(indexName string, searchColumn string) (E, error) {

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		// Addresses: []string{"http://127.0.0.1:9200"},
		Addresses:              []string{"https://127.0.0.1:9200"},
		Username:               "elastic",
		Password:               "ZvVx8w0v6CCtOFHK0477",
		CertificateFingerprint: "13374D6EA68778A39558E8FB53A23551B3C43B6F5A3295CAEFC9D046C5D71F51",
	})

	if err != nil {
		return E{}, err
	}

	_, err = client.Ping()
	if err != nil {
		return E{}, err
	}

	return E{
		C:                    client,
		IndexName:            indexName,
		FullTextSearchColumn: searchColumn,
	}, nil
}

func (e E) Info() (M, error) {
	res, err := e.C.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r M
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func (e E) Insert(ctx context.Context, i I) (uuid.UUID, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return uuid.UUID{}, err
	}
	id := e.GetId(i)
	req := esapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: id,
		Refresh:    "true",
		Body:       bytes.NewBuffer(data),
	}
	res, err := req.Do(ctx, e.C)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer res.Body.Close()

	return uuid.Parse(id)
}

func (e E) Search(ctx context.Context, q string) (SearchResponse, error) {
	var r SearchResponse
	var buf bytes.Buffer
	query := M{
		"query": M{
			"match": M{
				// "title": q,
				e.FullTextSearchColumn: q,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return r, err
	}

	// Perform the search request.
	res, err := e.C.Search(
		e.C.Search.WithContext(context.Background()),
		e.C.Search.WithIndex(e.IndexName),
		e.C.Search.WithBody(&buf),
		e.C.Search.WithTrackTotalHits(true),
		e.C.Search.WithPretty(),
	)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}

func (e E) Get(ctx context.Context, id string) (M, error) {
	span, s_ctx := opentracing.StartSpanFromContext(ctx, "E Get called")
	defer span.Finish()

	req := esapi.GetRequest{
		Index:      e.IndexName,
		DocumentID: id,
	}
	res, err := req.Do(s_ctx, e.C)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	var r M
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	return r, nil
}

func (e E) GetId(i I) string {
	typeOf := reflect.TypeOf(i)
	if typeOf.Kind() == reflect.Map {
		m, ok := i.(M)
		if ok {
			id, ok := m["id"]
			if ok {
				return rawId(id)
			}
		}
		m1, ok := i.(map[string]interface{})
		if ok {
			id, ok := m1["id"]
			if ok {
				return rawId(id)
			}
		}
	}
	if typeOf.Kind() == reflect.Struct {
		v := reflect.ValueOf(i)
		for i := 0; i < v.NumField(); i++ {
			field := typeOf.Field(i)
			tag, ok := field.Tag.Lookup("json")
			if ok {
				if tag == "id,omitempty" {
					idField := v.Field(i)
					idType := idField.Type()
					if idType.Kind() == reflect.String {
						return rawId(idField.String())
					}
					if idType.Kind() == reflect.Int {
						return rawId(idField.Int())
					}
					return uuid.New().String()
				}
			}
		}
	}
	return uuid.New().String()
}

func rawId(i interface{}) string {
	id, ok := i.(string)
	if ok {
		if id != "" {
			return id
		}
	}
	idInt, ok := i.(int)
	if ok {
		if idInt > 0 {
			return fmt.Sprintf("%v", idInt)
		}
	}
	return uuid.New().String()
}

type SearchResponse struct {
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
	TimedOut bool   `json:"timed_out"`
	Took     int64  `json:"took"`
}

type Hits struct {
	Hits     []Hit   `json:"hits"`
	MaxScore float64 `json:"max_score"`
	Total    Total   `json:"total"`
}

type Hit struct {
	ID     string  `json:"_id"`
	Index  string  `json:"_index"`
	Score  float64 `json:"_score"`
	Source M       `json:"_source"`
}

type Source struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Total struct {
	Relation string `json:"relation"`
	Value    int64  `json:"value"`
}

type Shards struct {
	Failed     int64 `json:"failed"`
	Skipped    int64 `json:"skipped"`
	Successful int64 `json:"successful"`
	Total      int64 `json:"total"`
}
