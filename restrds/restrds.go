package restrds

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/internal/protocol/rest"
	"github.com/aws/aws-sdk-go/service/rds"
)

func Fetch(dbid, path string) (io.ReadCloser, error) {
	svc := rds.New(&aws.Config{
		Region: aws.String("ap-northeast-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	// RDSのHandlerはQuery APIになっているのでをREST APIに変更
	svc.Handlers.Build.Clear()
	svc.Handlers.Build.PushBack(rest.Build)
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.Unmarshal.PushBack(rest.Unmarshal)

	out, err := DownloadCompleteDBLogFile(svc, &DownloadCompleteDBLogFileInput{
		DBInstanceIdentifier: dbid,
		LogFileName:          path,
	})
	if err != nil {
		return nil, err
	}
	return out.Body.Close(), nil
}

// implementation of DownloadCompleteDBLogFile

const opDownloadCompleteDBLogFile = "DownloadCompleteDBLogFile"

func DownloadCompleteDBLogFile(c *rds.RDS, input *DownloadCompleteDBLogFileInput) (*DownloadCompleteDBLogFileOutput, error) {
	req, out := DownloadCompleteDBLogFileRequest(c, input)
	err := req.Send()
	return out, err
}

func DownloadCompleteDBLogFileRequest(c *rds.RDS, input *DownloadCompleteDBLogFileInput) (req *request.Request, output *DownloadCompleteDBLogFileOutput) {
	if input == nil {
		input = &DownloadCompleteDBLogFileInput{}
	}

	op := &request.Operation{
		Name:       opDownloadCompleteDBLogFile,
		HTTPMethod: "GET",
		HTTPPath:   "/v13/downloadCompleteLogFile/{DBInstanceIdentifier}/{LogFileName+}",
	}

	req = c.NewRequest(op, input, output)
	output = &DownloadCompleteDBLogFileOutput{}
	req.Data = output
	return
}

type DownloadCompleteDBLogFileInput struct {
	DBInstanceIdentifier *string `location:"uri" locationName:"DBInstanceIdentifier" type:"string" required:"true"`

	LogFileName *string `location:"uri" locationName:"LogFileName" type:"string" required:"true"`
}

type DownloadCompleteDBLogFileOutput struct {
	Body io.ReadCloser `type:"blob"`

	metadataDownloadCompleteDBLogFileOutput `json:"-" xml:"-"`
}

type metadataDownloadCompleteDBLogFileOutput struct {
	SDKShapeTraits bool `type:"structure" payload:"Body"`
}
