package sqsutil

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SendMessage to send SQS message
func (s *Service) SendMessage(queueURL string, message string) (*sqs.SendMessageOutput, error) {
	// ! double json encode
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	result, err := s.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonMsg)),
		QueueUrl:    aws.String(queueURL),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetQueueURL represents URL of the queue we want to send a message to
func (s *Service) GetQueueURL(queue string) (*sqs.GetQueueUrlOutput, error) {
	result, err := s.sqs.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMessages to get SQS message
func (s *Service) GetMessages(queueURL string) (*sqs.ReceiveMessageOutput, error) {
	msgResult, err := s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
	})
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

// DeleteMessageBatch to delete SQS message
func (s *Service) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	msgResult, err := s.sqs.DeleteMessageBatch(input)
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

// SendMessageDelay to send message to SQS with delay seconds
func (s *Service) SendMessageDelay(queueURL string, message map[string]interface{}, delaySeconds int64) (*sqs.SendMessageOutput, error) {
	// ! double json encode
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	result, err := s.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(string(jsonMsg)),
		QueueUrl:     aws.String(queueURL),
		DelaySeconds: &delaySeconds,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendMessageBatch to send batch of message to SQS
func (s *Service) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	msgResult, err := s.sqs.SendMessageBatch(input)
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}
