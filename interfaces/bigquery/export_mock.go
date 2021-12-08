//go:build test
// +build test

package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"reflect"
)

type mockBigqueryClientImpl struct {
	bqClient *bigquery.Client
	csItems  []ChangeStreamTableSchema
}

func (m *mockBigqueryClientImpl) putRecord(_ context.Context, _ string, _ string, csItems []ChangeStreamTableSchema) error {
	if csItems == nil {
		return fmt.Errorf("Expect csItems to not be nil.")
	}

	fmt.Printf("%v\n", csItems)
	if e, a := m.csItems, csItems; !reflect.DeepEqual(e, a) {
		return fmt.Errorf("expect %v, got %v", e, a)
	}
	return nil
}
