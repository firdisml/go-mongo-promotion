package configs

import (
	"bytes"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Storage *session.Session = ConnectStorage()

func ConnectStorage() *session.Session {

	credential := credentials.NewStaticCredentials(
		Env("AWS_S3_ACCESS_ID"),
		Env("AWS_S3_SECRET_KEY"),
		"",
	)

	aws_config := aws.Config{
		Region:      aws.String(Env("AWS_S3_REGION")),
		Credentials: credential,
	}

	session, err := session.NewSession(&aws_config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Storage")

	return session

}

func UploadFile(session *session.Session, file_size int64, file_buffer []byte, promotion_id primitive.ObjectID) *s3.PutObjectOutput {
	result, upload_error := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String("promotion-app-storage"),
		Key:           aws.String(promotion_id.Hex()),
		Body:          bytes.NewReader(file_buffer),
		ContentLength: aws.Int64(int64(file_size)),
		ContentType:   aws.String(http.DetectContentType(file_buffer)),
		StorageClass:  aws.String("INTELLIGENT_TIERING"),
	})
	if upload_error != nil {
		log.Fatal(upload_error)
	}

	return result

}
