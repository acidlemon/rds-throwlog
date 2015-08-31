package main

import (
	"io"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rds"
)

// implementation of DownloadCompleteDBLogFile

const opDownloadCompleteDBLogFile = "DownloadCompleteDBLogFile"

func rdsDownloadCompleteDBLogFile(c *rds.RDS, input *DownloadCompleteDBLogFileInput) (*DownloadCompleteDBLogFileOutput, error) {
	req, out := rdsDownloadCompleteDBLogFileRequest(c, input)
	err := req.Send()
	return out, err
}

func rdsDownloadCompleteDBLogFileRequest(c *rds.RDS, input *DownloadCompleteDBLogFileInput) (req *request.Request, output *DownloadCompleteDBLogFileOutput) {
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
