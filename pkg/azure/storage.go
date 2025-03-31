package azure

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func GetTemplateFromBlob(blobName string) (string, error) {
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER")
	if containerName == "" {
		return "", fmt.Errorf("variável de ambiente AZURE_STORAGE_CONTAINER não definida")
	}

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	URL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("erro ao analisar a URL: %w", err)
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("erro ao criar credenciais para o Azure Blob Storage: %w", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceURL := azblob.NewServiceURL(*parsedURL, pipeline)

	containerURL := serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlobURL(blobName)

	ctx := context.Background()
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return "", fmt.Errorf("erro ao baixar o template: %w", err)
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	bodyBytes, err := ioutil.ReadAll(bodyStream)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o conteúdo do blob: %w", err)
	}

	return string(bodyBytes), nil
}
