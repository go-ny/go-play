package awsOp

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-play/common/getEnv"
	"log"
	"net/http"
	"strings"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var filepath string

func UploadImage(c *gin.Context) {
	log.Println("start to upload image...")

	sess := c.MustGet("sess").(*session.Session)

	uploader := s3manager.NewUploader(sess)
	MyBucket := getEnv.EnvWithKey("BUCKET_NAME")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error":  fmt.Sprintf("Failed to read from FormFile: %v" ,err),
		})
		return
	}

	filename := uuid.New().String() + "." +  strings.Split(header.Filename, ".")[1]

	log.Println("file name: \n", filename)

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput {
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error":    "Failed to upload file",
			"uploader": up,
		})
		fmt.Println("failed to upload file")
		fmt.Println("here's the error: ", err)
		return
	}
	//filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	filepath = "https://" + MyBucket + "." + "s3" + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})
	fmt.Println("filepath: ", filepath)
}

func ConnectAws() *session.Session {
	AccessKeyID = getEnv.EnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = getEnv.EnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = getEnv.EnvWithKey("AWS_REGION")
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func S3UploadImageAPI(c *gin.Context) (string, error) {
	log.Println("start to upload image...")

	sess := c.MustGet("sess").(*session.Session)

	uploader := s3manager.NewUploader(sess)
	MyBucket := getEnv.EnvWithKey("BUCKET_NAME")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error":  fmt.Sprintf("Failed to read from FormFile: %v" ,err),
		})
		return "nil", err
	}

	filename := uuid.New().String() + "." +  strings.Split(header.Filename, ".")[1]

	log.Println("file name: \n", filename)

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput {
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error":    "Failed to upload file",
			"uploader": up,
		})
		fmt.Println("failed to upload file")
		fmt.Println("here's the error: ", err)
		return "nil", err
	}
	//filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	filepath = "https://" + MyBucket + "." + "s3" + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})
	fmt.Println("filepath: ", filepath)

	return filepath, nil
}