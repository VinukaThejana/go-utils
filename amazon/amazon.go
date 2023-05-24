// Package amazon is used to create a
// AWS session and initiate diffrent AWS
// services on the AWS session
package amazon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

var log logger.Logger

// AWS is a struct that contains the methods that
// needs to be implemented on AWS
type AWS struct {
	Session *session.Session
}

// Init is a funtion to intialize the AWS session
func (a *AWS) Init(AwsAccessKeyID, AwsSecretAccsessKey, AwsRegion string) error {
	config := aws.Config{
		Region: aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			AwsAccessKeyID,
			AwsSecretAccsessKey,
			"",
		),
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return err
	}

	a.Session = sess
	return nil
}

// Upload is a method on AWS to upload the given image/video
// to the amazon s3 storage bucket
func (a AWS) Upload(AwsS3StorageBucketName, AwsS3Region, fileName string) (*string, *string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	fileSize := fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}

	_, err = s3.New(a.Session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AwsS3StorageBucketName),
		Key:                  aws.String(uid.String()),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(int64(fileSize)),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return nil, nil, err
	}

	uidString := uid.String()
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AwsS3StorageBucketName, AwsS3Region, uid.String())
	return &url, &uidString, nil
}

// SendMessage is a method on AWS to send messages to SQS queue
func (a AWS) SendMessage(sqsURL string, payload interface{}) {
	// Marshal the payload of the data sent to SQS
	jsonData, err := json.Marshal(payload)
	if err != nil {
		errMsg := "Failed to marshal the data"
		log.Error(err, &errMsg)
		return
	}

	sqsClient := sqs.New(a.Session)
	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &sqsURL,
		MessageBody: aws.String(string(jsonData)),
	})
}
