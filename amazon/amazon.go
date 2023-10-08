// Package amazon is used to initialize the AWS session and provide functionality
// to various AWS services
package amazon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

// AWS struct contains the AWS session
type AWS struct {
	Session *session.Session
}

// Init is a function that is used to initialize the session connection to AWS
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

// Upload provides the ability to upload an Image or a video to the configured
// amazon s3 storage bucket
func (a AWS) Upload(AwsS3StorageBucketName, AwsS3Region, fileName string) (url *string, uid *string, err error) {
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

	itemUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}

	_, err = s3.New(a.Session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AwsS3StorageBucketName),
		Key:                  aws.String(itemUUID.String()),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(int64(fileSize)),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return nil, nil, err
	}

	uidString := itemUUID.String()
	*url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AwsS3StorageBucketName, AwsS3Region, itemUUID.String())
	return url, &uidString, nil
}

// SendMessage provides the ability to send messaged to a given SQS queue
func (a AWS) SendMessage(sqsURL string, payload interface{}) (err error) {
	// Marshal the payload of the data sent to SQS
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal the data")
	}

	sqsClient := sqs.New(a.Session)
	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &sqsURL,
		MessageBody: aws.String(string(jsonData)),
	})

	return err
}
