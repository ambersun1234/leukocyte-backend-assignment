package service

import (
	"encoding/json"
	"errors"
	"testing"

	"leukocyte/src/logger"
	mockQueue "leukocyte/src/message_queue/mocks"
	mockOrchestration "leukocyte/src/orchestration/mocks"
	"leukocyte/src/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestConsumerSuite struct {
	suite.Suite

	routingKey string
	consumer   *Consumer
	mq         *mockQueue.Queue
	orch       *mockOrchestration.Orchestration
}

func TestConsumer(t *testing.T) {
	suite.Run(t, new(TestConsumerSuite))
}

func (suite *TestConsumerSuite) exampleWorkerData() types.JobObject {
	return types.JobObject{
		Namespace:     "test",
		Name:          "test",
		Image:         "test",
		RestartPolicy: "test",
		Commands:      []string{"test"},
	}
}

func (suite *TestConsumerSuite) SetupTest() {
	suite.routingKey = "test"
	suite.mq = mockQueue.NewQueue(suite.T())
	suite.orch = mockOrchestration.NewOrchestration(suite.T())
	suite.consumer = NewConsumer(
		logger.NewTestLogger(), suite.orch, suite.mq, suite.routingKey,
	)
}

func (suite *TestConsumerSuite) TearDownTest() {
	suite.mq = nil
	suite.orch = nil
	suite.consumer = nil
}

func (suite *TestConsumerSuite) TestEnqueue() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(nil)

	err := suite.consumer.enqueue(suite.routingKey)
	suite.NoError(err)
}

func (suite *TestConsumerSuite) TestEnqueueFail() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(errors.New("failed to publish message"))

	err := suite.consumer.enqueue(suite.routingKey)
	suite.Error(err)
}

func (suite *TestConsumerSuite) TestWorker() {
	suite.orch.On("Schedule", mock.Anything).Return(nil)

	data := suite.exampleWorkerData()
	dataStr, err := json.Marshal(data)
	suite.NoError(err)

	err = suite.consumer.Worker(string(dataStr))
	suite.NoError(err)

	suite.mq.AssertNumberOfCalls(suite.T(), "Publish", 0)
}

func (suite *TestConsumerSuite) TestWorkerUnmarshalFail() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(nil)

	err := suite.consumer.Worker("invalid")
	suite.NoError(err)

	suite.mq.AssertNumberOfCalls(suite.T(), "Publish", 1)
}

func (suite *TestConsumerSuite) TestWorkerScheduleFail() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(nil)
	suite.orch.On("Schedule", mock.Anything).Return(errors.New("failed to schedule job"))

	data := suite.exampleWorkerData()
	dataStr, err := json.Marshal(data)
	suite.NoError(err)

	err = suite.consumer.Worker(string(dataStr))
	suite.NoError(err)

	suite.mq.AssertNumberOfCalls(suite.T(), "Publish", 1)
}

func (suite *TestConsumerSuite) TestStart() {
	suite.mq.On("Consume", suite.routingKey, mock.Anything).Return(nil)

	err := suite.consumer.Start()
	suite.NoError(err)
}

func (suite *TestConsumerSuite) TestStartFail() {
	suite.mq.On("Consume", suite.routingKey, mock.Anything).Return(errors.New("failed to consume message"))

	err := suite.consumer.Start()
	suite.Error(err)
}
