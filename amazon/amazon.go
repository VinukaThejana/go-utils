// Package amazon is used to create a
// AWS session and initiate diffrent AWS
// services on the AWS session
package amazon

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VinukaThejana/go-utils/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// AWS is a struct that contains the methods that
// needs to be implemented on AWS
type AWS struct {
	Session *session.Session
}

// InitAWS is a funtion to intialize the AWS session
func InitAWS(AwsS3Region, AwsAccessKeyID, AwsSecretAccsessKey string) (*session.Session, errors.Status) {
	config := aws.Config{
		Region: aws.String(AwsS3Region),
		Credentials: credentials.NewStaticCredentials(
			AwsAccessKeyID,
			AwsSecretAccsessKey,
			"",
		),
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		log.Println("Failed to intialize the connection with Amazon")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	return sess, errors.Okay
}

// Upload is a method on AWS to upload the given image/video
// to the amazon s3 storage bucket
func (a AWS) Upload(AwsS3StorageBucketName, AwsS3Region, fileName string) (*string, errors.Status) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Failed to open the file")
		log.Println(err)
		return nil, errors.InternalServerError
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Failed to get file stat")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	fileSize := fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Println("Failed to generate UUID")
		log.Println(err)
		return nil, errors.InternalServerError
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
		log.Println("Failed to upload !")
		log.Println(err)
		return nil, errors.InternalServerError
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AwsS3StorageBucketName, AwsS3Region, uid.String())
	return &url, errors.Okay
}
