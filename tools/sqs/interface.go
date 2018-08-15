package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-msgqueue/msgqueue"
	"github.com/go-msgqueue/msgqueue/azsqs"
	"github.com/tokenme/tmm/common"
)

type MsgType = uint

const (
	TransactionMsg MsgType = 0
	GasMsg         MsgType = 1
)

func NewManager(config common.SQSConfig) msgqueue.Manager {
	return azsqs.NewManager(awsSQS(config), config.AccountId)
}

func NewQueue(m msgqueue.Manager, opt *msgqueue.Options) msgqueue.Queue {
	return m.NewQueue(opt)
}

func awsSQS(config common.SQSConfig) *sqs.SQS {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AK, config.Secret, config.Token),
	}))
	return sqs.New(sess)
}

type Message interface {
	String() string
}

type Queue interface {
	Start()
	Stop()
}
