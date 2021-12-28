package data

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
)

type AzureBlobEnv struct {
	AZUREBLOB_SECRET_ACCOUNTNAME  string
	AZUREBLOB_SECRET_MYACCOUNTKEY string
	AZUREBLOB_SECRET_MYACCOUNTURL string
	sclnt                         azblob.ServiceClient
}

type BlobSource interface {
	GetServiceClient() (azblob.ServiceClient, error)
	GetContainerClient() (azblob.ContainerClient, error)
	GetBlockBlobClient() (azblob.BlockBlobClient, error)
}

type BlobData struct {
	data string `json:"data"`
}

var blob = &AzureBlobEnv{"", "", "", azblob.ServiceClient{}}

func (b *AzureBlobEnv) GetServiceClient() (azblob.ServiceClient, error) {
	if (azblob.ServiceClient{}) == b.sclnt {
		var err error
		// load .env file
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("No Local .env file. So accessing environment variables from Kyma runtime")
		}
		b.AZUREBLOB_SECRET_ACCOUNTNAME = goDotEnvVariable("AZUREBLOB_SECRET_ACCOUNTNAME")
		b.AZUREBLOB_SECRET_MYACCOUNTKEY = goDotEnvVariable("AZUREBLOB_SECRET_MYACCOUNTKEY")
		b.AZUREBLOB_SECRET_MYACCOUNTURL = goDotEnvVariable("AZUREBLOB_SECRET_MYACCOUNTURL")
		cred, err := azblob.NewSharedKeyCredential(b.AZUREBLOB_SECRET_ACCOUNTNAME, b.AZUREBLOB_SECRET_MYACCOUNTKEY)
		if err != nil {
			fmt.Println(err)
			return b.sclnt, err
		}
		sclnt, err := azblob.NewServiceClientWithSharedKey(b.AZUREBLOB_SECRET_MYACCOUNTURL, cred, nil)
		if err != nil {
			fmt.Println(err)
			return b.sclnt, err
		}
		b.sclnt = sclnt
		return b.sclnt, nil
	} else {
		return b.sclnt, nil
	}
}

func (b *AzureBlobEnv) GetContainerClient(c string) (azblob.ContainerClient, error) {
	sclnt, err := b.GetServiceClient()
	if err != nil {
		fmt.Println(err)
		return azblob.ContainerClient{}, err
	}
	cclnt := sclnt.NewContainerClient(c)
	return cclnt, nil
}

func (b *AzureBlobEnv) GetBlockBlobClient(c string, bl string) (azblob.BlockBlobClient, error) {
	cclnt, err := b.GetContainerClient(c)
	if err != nil {
		fmt.Println(err)
		return azblob.BlockBlobClient{}, err
	}
	bclnt := cclnt.NewBlockBlobClient(bl)
	return bclnt, nil
}

func GetBlockBlobClient(c string, bl string) (azblob.BlockBlobClient, error) {
	sclnt, err := blob.GetServiceClient()
	if err != nil {
		fmt.Println(err)
		return azblob.BlockBlobClient{}, err
	}
	cclnt := sclnt.NewContainerClient(c)
	bclnt := cclnt.NewBlockBlobClient(bl)
	return bclnt, nil
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

func CreateTestBlob() {
	ctx := context.Background()
	data := "Hello Venu!!"
	bclnt, err := blob.GetBlockBlobClient("test", "HelloVenu.txt")
	if err != nil {
		fmt.Println(err)
	}
	_, err = bclnt.Upload(ctx, NopCloser(strings.NewReader(data)), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateTestCaseBlob(t RawTestCase) {
	ctx := context.Background()
	data := t.Filedata
	dataDec, _ := base64.StdEncoding.DecodeString(data)
	dataDecString := string(dataDec)
	bclnt, err := blob.GetBlockBlobClient("testcases/"+t.Name+"/"+t.Testcase, t.Filename)
	if err != nil {
		fmt.Println(err)
	}
	r, err := bclnt.Upload(ctx, NopCloser(strings.NewReader(dataDecString)), nil)
	log.Println(r.RawResponse)
	if err != nil {
		fmt.Println(err)
	}
}

func GetTestCasePayloadBlob(f string) []byte {
	ctx := context.Background()
	sarr := strings.Split(f, "/")
	bclnt, err := blob.GetBlockBlobClient("testcases/"+sarr[1]+"/"+sarr[2], sarr[3])
	if err != nil {
		fmt.Println(err)
	}
	r, err := bclnt.Download(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	// Open a buffer, reader, and then download!
	downloadedData := &bytes.Buffer{}
	reader := r.Body(azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		log.Fatalln(err)
		err = reader.Close()
	}
	var downloadedDataStr = downloadedData.String()
	b := map[string]string{"data": base64.StdEncoding.EncodeToString([]byte(downloadedDataStr))}
	bjson, err := json.Marshal(b)
	if err != nil {
		log.Fatalln(err)
	}
	return bjson
}

func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}
