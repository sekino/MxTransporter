//go:build test
// +build test

package application

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	interfaceForBigquery "mxtransporter/interfaces/bigquery"
	interfaceForKinesisStream "mxtransporter/interfaces/kinesis-stream"
	interfaceForPubsub "mxtransporter/interfaces/pubsub"
	interfaceForResumeToken "mxtransporter/usecases/resume-token"
)

type mockChangeStremsWatcherClientImpl struct {
	mongoClient            *mongo.Client
	csExporter             ChangeStreamsExporterImpl
	resumeToken            string
	resumeAfterExistence   bool
	bqPassCheck            string
	pubsubPassCheck        string
	kinesisStreamPassCheck string
}

func (m *mockChangeStremsWatcherClientImpl) newBigqueryClient(_ context.Context, _ string) (*bigquery.Client, error) {
	m.bqPassCheck = "OK"
	return nil, nil
}

func (m *mockChangeStremsWatcherClientImpl) newPubsubClient(_ context.Context, _ string) (*pubsub.Client, error) {
	m.pubsubPassCheck = "OK"
	return nil, nil
}

func (m *mockChangeStremsWatcherClientImpl) newKinesisClient(_ context.Context) (*kinesis.Client, error) {
	m.kinesisStreamPassCheck = "OK"
	return nil, nil
}

func (m *mockChangeStremsWatcherClientImpl) watch(_ context.Context, ops *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if ops.ResumeAfter != nil {
		if ops.ResumeAfter.(map[string]string)["_data"] == m.resumeToken {
			m.resumeAfterExistence = true
		}
	} else {
		m.resumeAfterExistence = false
	}

	return nil, nil
}

func (c *mockChangeStremsWatcherClientImpl) setCsExporter(_ ChangeStreamsExporterImpl) {
}

func (c *mockChangeStremsWatcherClientImpl) exportChangeStreams(_ context.Context) error {
	return nil
}

type mockChangeStreamsExporterClientImpl struct {
	cs                     primitive.M
	bq                     interfaceForBigquery.BigqueryImpl
	pubsub                 interfaceForPubsub.PubsubImpl
	kinesisStream          interfaceForKinesisStream.KinesisStreamImpl
	resumeToken            interfaceForResumeToken.ResumeTokenImpl
	bqPassCheck            string
	pubsubPassCheck        string
	kinesisStreamPassCheck string
	csCursorFlag           bool
}

func (m *mockChangeStreamsExporterClientImpl) next(_ context.Context) bool {
	return m.csCursorFlag
}

func (m *mockChangeStreamsExporterClientImpl) decode() (primitive.M, error) {
	return m.cs, nil
}

func (_ *mockChangeStreamsExporterClientImpl) close(_ context.Context) error {
	return nil
}

func (m *mockChangeStreamsExporterClientImpl) exportToBigquery(_ context.Context, _ primitive.M) error {
	m.bqPassCheck = "OK"
	return nil
}

func (m *mockChangeStreamsExporterClientImpl) exportToPubsub(_ context.Context, _ primitive.M) error {
	m.pubsubPassCheck = "OK"
	return nil
}

func (m *mockChangeStreamsExporterClientImpl) exportToKinesisStream(_ context.Context, _ primitive.M) error {
	m.kinesisStreamPassCheck = "OK"
	return nil
}

func (m *mockChangeStreamsExporterClientImpl) saveResumeToken(_ string) error {
	m.csCursorFlag = false
	return nil
}
