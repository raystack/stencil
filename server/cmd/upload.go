package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/odpf/stencil/server/config"
	"github.com/spf13/cobra"
)

var (
	filePath  string
	uploadCmd = &cobra.Command{
		Use:   "upload [namespace] [name] [version]",
		Args:  cobra.MinimumNArgs(3),
		Short: "upload file descriptor to stencil server",
		RunE:  upload,
	}
)

func init() {
	uploadCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to descriptor file")
	rootCmd.AddCommand(uploadCmd)
}

func upload(cmd *cobra.Command, args []string) error {
	config := config.LoadConfig()
	namespace := args[0]
	name := args[1]
	version := args[2]
	return uploadDescriptor(config, namespace, name, version, filePath)
}

func uploadDescriptor(config *config.Config, namespace string, name string, version string, descriptorFile string) error {
	url := url.URL{
		Scheme: config.Scheme,
		Host:   config.Host + ":" + config.Port,
		Path:   path.Join("v1", "namespaces", namespace, "descriptors/"),
	}
	req, err := newfileUploadRequest(
		url.String(),
		map[string]string{
			"name":    name,
			"version": version,
			"latest":  "true",
		},
		"file",
		descriptorFile,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+config.AuthBearerToken)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("Error: respose code %d from stencil server: %s", res.StatusCode, string(body))
	}
	fmt.Println(string(body))
	return nil
}

// Creates a new file upload http request with extra params
func newfileUploadRequest(uri string, params map[string]string, fileParamName string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParamName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}
