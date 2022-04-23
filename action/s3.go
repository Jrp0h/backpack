package action

import (
	"fmt"
	"os"
	"path"

	"github.com/Jrp0h/backpack/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type s3Action struct {
	bucket       string
	region       string
	clientID     string
	clientSecret string
	token        string
}

func (action *s3Action) createConnection() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      &action.region,
		Credentials: credentials.NewStaticCredentials(action.clientID, action.clientSecret, action.token),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (action *s3Action) CanValidateConnection() bool {
	return false // TODO: Fix this
}

func (action *s3Action) TestConnection() error {
	_, err := action.createConnection()
	return err
}

func (action *s3Action) Upload(fileData *utils.FileData) error {
	session, err := action.createConnection()

	if err != nil {
		return err
	}

	f, err := os.Open(fileData.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if file with the same name already exists
	_, err = s3.New(session).HeadObject(&s3.HeadObjectInput{
		Bucket: &action.bucket,
		Key:    &fileData.Name,
	})
	if err == nil {
		return fmt.Errorf("file %s already exists", fileData.Path)
	}

	uploader := s3manager.NewUploader(session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &action.bucket,
		Key:    &fileData.Name,
		Body:   f,
	})

	if err != nil {
		return err
	}

	return nil
}

func (action *s3Action) ListFiles() ([]string, error) {
	session, err := action.createConnection()
	if err != nil {
		return nil, err
	}

	result, err := s3.New(session).ListObjects(&s3.ListObjectsInput{
		Bucket: &action.bucket,
	})
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, obj := range result.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

func (action *s3Action) Fetch(file string) (string, error) {
	session, err := action.createConnection()

	if err != nil {
		return "", err
	}

	outPath := path.Join(os.TempDir(), uuid.NewString()+".zip")
	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	uploader := s3manager.NewDownloader(session)
	_, err = uploader.Download(outFile, &s3.GetObjectInput{
		Bucket: &action.bucket,
		Key:    &file,
	})

	if err != nil {
		return "", err
	}

	return outPath, nil
}

func loadS3Action(data *map[string]string) (Action, error) {
	// Required
	bucket, err := utils.ValueOrErrorString(data, "bucket", "action/s3")
	if err != nil {
		return nil, err
	}

	region, err := utils.ValueOrErrorString(data, "region", "action/s3")
	if err != nil {
		return nil, err
	}

	id, err := utils.ValueOrErrorString(data, "client_id", "action/s3")
	if err != nil {
		return nil, err
	}

	secret, err := utils.ValueOrErrorString(data, "client_secret", "action/s3")
	if err != nil {
		return nil, err
	}

	token := utils.ValueOrDefaultString(data, "token", "")

	return &s3Action{
		bucket,
		region,
		id,
		secret,
		token,
	}, nil
}
