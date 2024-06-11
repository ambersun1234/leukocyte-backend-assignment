package service

import (
	"context"
	"errors"
	"testing"

	"leukocyte/src/logger"
	"leukocyte/src/message_queue/mocks"
	"leukocyte/src/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestProducerSuite struct {
	suite.Suite

	routingKey types.RoutingKey
	producer   *Producer
	mq         *mocks.Queue
}

func TestProducer(t *testing.T) {
	suite.Run(t, new(TestProducerSuite))
}

func (suite *TestProducerSuite) SetupTest() {
	suite.routingKey = "test"
	suite.mq = mocks.NewQueue(suite.T())
	suite.producer = NewProducer(
		context.Background(), logger.NewTestLogger(),
		suite.mq, suite.routingKey,
	)
}

func (suite *TestProducerSuite) TearDownTest() {
	suite.mq = nil
	suite.producer = nil
}

func (suite *TestProducerSuite) TestProduceMessage() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(nil)

	err := suite.producer.produce(1)
	suite.NoError(err)
}

func (suite *TestProducerSuite) TestProduceMessageFail() {
	suite.mq.On("Publish", suite.routingKey, mock.Anything).Return(errors.New("failed to publish message"))

	err := suite.producer.produce(1)
	suite.Error(err)
}
