package stencil_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	stencil "github.com/goto/stencil/clients/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func runProtoc(
	rootDir string,
	includeImports bool,
	descSetOut string,
	filePaths []string,
) error {
	protocBinPath, err := exec.LookPath("protoc")
	if err != nil {
		return err
	}
	protocBinPath, err = filepath.EvalSymlinks(protocBinPath)
	if err != nil {
		return err
	}
	protocBinPath, err = filepath.Abs(protocBinPath)
	if err != nil {
		return err
	}
	protocIncludePath, err := filepath.Abs(filepath.Join(filepath.Dir(protocBinPath), "..", "include"))
	if err != nil {
		return err
	}
	args := []string{"-I", rootDir, "-I", protocIncludePath}
	args = append(args, fmt.Sprintf("--descriptor_set_out=%s", descSetOut))
	if includeImports {
		args = append(args, "--include_imports")
	}
	args = append(args, filePaths...)
	stderr := bytes.NewBuffer(nil)
	cmd := exec.Command(protocBinPath, args...)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s returned error: %v %v", protocBinPath, err, stderr.String())
	}
	return nil
}

func getDescriptorDataByPath(t *testing.T, includeImports bool, rootPath string) ([]byte, error) {
	root, _ := filepath.Abs(rootPath)
	fileName := filepath.Join(t.TempDir(), "file.desc")
	rootFiles, _ := filepath.Glob(filepath.Join(root, "./*.proto"))
	err := runProtoc(root, includeImports, fileName, rootFiles)
	assert.NoError(t, err)
	data, err := ioutil.ReadFile(fileName)
	return data, err
}

func getDescriptorData(t *testing.T, includeImports bool) ([]byte, error) {
	return getDescriptorDataByPath(t, includeImports, "./test_data")
}

func getUpdatedDescriptorDataAndMsgData(t *testing.T, includeImports bool) ([]byte, []byte) {
	data, err := getDescriptorDataByPath(t, includeImports, "./test_data/updated")
	assert.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer ts.Close()
	url := ts.URL
	client, err := stencil.NewClient([]string{url}, stencil.Options{})
	assert.Nil(t, err)
	assert.NotNil(t, client)
	msgDesc, err := client.GetDescriptor("test.stencil.One")
	assert.NoError(t, err)
	//construct message
	msg := dynamicpb.NewMessage(msgDesc).New()
	fieldOne := msgDesc.Fields().ByName("field_one")
	msg.Set(fieldOne, protoreflect.ValueOfInt64(200))
	fieldTwo := msgDesc.Fields().ByName("field_two")
	msg.Set(fieldTwo, protoreflect.ValueOfInt64(300))

	msgData, err := proto.Marshal(msg.Interface())
	assert.NoError(t, err)
	return data, msgData
}

func TestNewClient(t *testing.T) {
	t.Run("should return error if url is not valid", func(t *testing.T) {
		url := "h_ttp://invalidurl"
		_, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Contains(t, err.Error(), "invalid request")
	})

	t.Run("should return error if request fails", func(t *testing.T) {
		url := "ithttp://localhost"
		_, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Contains(t, err.Error(), "request failed")
	})

	t.Run("should return error if file download fails", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}))
		defer ts.Close()
		url := ts.URL
		_, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Contains(t, err.Error(), "request failed.")
	})

	t.Run("should return error if downloaded file is not valid", func(t *testing.T) {
		data := []byte("invalid")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		_, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Contains(t, err.Error(), "invalid file descriptorset file.")
	})

	t.Run("should return error if downloaded file is not fully contained file", func(t *testing.T) {
		data, err := getDescriptorData(t, false)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		_, err = stencil.NewClient([]string{url}, stencil.Options{})
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "file is not fully contained descriptor file.")
		}
	})

	t.Run("should create a client if provided descriptor file is valid", func(t *testing.T) {
		data, err := getDescriptorData(t, true)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		client, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Nil(t, err)
		assert.NotNil(t, client)
	})

	t.Run("should pass provided headers to request", func(t *testing.T) {
		data, _ := getDescriptorData(t, true)
		headers := map[string]string{
			"key": "value",
		}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, val := range headers {
				assert.Equal(t, val, r.Header.Get(key))
			}
			w.Write(data)
		}))
		client, err := stencil.NewClient([]string{ts.URL}, stencil.Options{HTTPOptions: stencil.HTTPOptions{Headers: headers}})
		assert.Nil(t, err)
		assert.NotNil(t, client)
	})

	t.Run("should refresh descriptors by specified intervals", func(t *testing.T) {
		data, _ := getDescriptorData(t, true)
		callCount := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Write(data)
		}))
		client, _ := stencil.NewClient([]string{ts.URL}, stencil.Options{AutoRefresh: true, RefreshInterval: 2 * time.Millisecond})
		// wait for interval to end
		time.Sleep(3 * time.Millisecond)
		client.GetDescriptor("test.One")
		time.Sleep(1 * time.Millisecond)
		client.Close()
		assert.Equal(t, 2, callCount)
	})
}

func TestClient(t *testing.T) {
	t.Run("GetDescriptor", func(t *testing.T) {
		data, err := getDescriptorData(t, true)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		client, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Nil(t, err)
		assert.NotNil(t, client)
		t.Run("should return notFoundErr if not found", func(t *testing.T) {
			msg, err := client.GetDescriptor("test.stencil.Two.Unknown")
			assert.Nil(t, msg)
			assert.NotNil(t, err)
			assert.Equal(t, stencil.ErrNotFound, err)
		})
		t.Run("should get nested message descriptor from fully qualified java classname", func(t *testing.T) {
			msg, err := client.GetDescriptor("test.stencil.Two.Four")
			assert.Nil(t, err)
			field := msg.Fields().ByName("recursive")
			assert.NotNil(t, field)
		})
		t.Run("should get nested message descriptor if java_option not specified", func(t *testing.T) {
			msg, err := client.GetDescriptor("test.Three")
			assert.Nil(t, err)
			field := msg.Fields().ByName("field_one")
			assert.NotNil(t, field)
		})
		t.Run("should get descriptor if package name is not defined", func(t *testing.T) {
			msg, err := client.GetDescriptor("Root")
			assert.Nil(t, err)
			field := msg.Fields().ByName("field_one")
			assert.NotNil(t, field)
		})
		t.Run("should get descriptor if proto package name is not defined but java package is defined", func(t *testing.T) {
			msg, err := client.GetDescriptor("test.stencil.Root")
			assert.Nil(t, err)
			field := msg.Fields().ByName("field_one")
			assert.NotNil(t, field)
		})
	})
	t.Run("Parse", func(t *testing.T) {
		data, err := getDescriptorData(t, true)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		client, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Nil(t, err)
		assert.NotNil(t, client)
		t.Run("should return notFoundErr if not found", func(t *testing.T) {
			msg, err := client.Parse("test.stencil.Two.Unknown", []byte(""))
			assert.Nil(t, msg)
			assert.NotNil(t, err)
			assert.Equal(t, stencil.ErrNotFound, err)
		})

		t.Run("should parse wire format data given className", func(t *testing.T) {
			msgDesc, err := client.GetDescriptor("test.stencil.One")
			assert.NoError(t, err)
			//construct message
			msg := dynamicpb.NewMessage(msgDesc).New()
			fieldOne := msgDesc.Fields().ByName("field_one")
			msg.Set(fieldOne, protoreflect.ValueOfInt64(200))

			bytesData, err := proto.Marshal(msg.Interface())

			assert.NoError(t, err)
			parsed, err := client.Parse("test.stencil.One", bytesData)
			assert.Nil(t, err)
			assert.NotNil(t, parsed)
			val := parsed.ProtoReflect().Get(fieldOne)
			assert.Equal(t, int64(200), val.Int())
			assert.Nil(t, parsed.ProtoReflect().GetUnknown())
		})

		t.Run("should parse extensions without having any unknown fields", func(t *testing.T) {
			msgDesc, err := client.GetDescriptor("test.ExtendableMessage")
			assert.NoError(t, err)
			//construct message
			msg := dynamicpb.NewMessage(msgDesc).New()
			fieldOne := msgDesc.Fields().ByName("field_extra")
			msg.Set(fieldOne, protoreflect.ValueOfInt64(200))

			extenderMsgDesc, err := client.GetDescriptor("test.Extender")
			assert.NoError(t, err)
			fieldTwoDesc := extenderMsgDesc.Extensions().ByName("field_two")
			fieldTwoType := dynamicpb.NewExtensionType(fieldTwoDesc)
			fieldTwo := fieldTwoType.TypeDescriptor()
			proto.SetExtension(msg.Interface(), fieldTwoType, "field_two_value")

			bytesData, err := proto.Marshal(msg.Interface())

			assert.NoError(t, err)
			parsed, err := client.Parse("test.ExtendableMessage", bytesData)
			assert.Nil(t, err)
			assert.NotNil(t, parsed)
			val := parsed.ProtoReflect().Get(fieldOne)
			assert.Equal(t, int64(200), val.Int())
			val2 := parsed.ProtoReflect().Get(fieldTwo)
			assert.Equal(t, "field_two_value", val2.String())
			assert.Nil(t, parsed.ProtoReflect().GetUnknown())
		})
	})

	t.Run("Serialize", func(t *testing.T) {
		desc, err := getDescriptorData(t, true)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(desc)
		}))
		defer ts.Close()
		url := ts.URL
		client, err := stencil.NewClient([]string{url}, stencil.Options{})
		assert.Nil(t, err)
		assert.NotNil(t, client)

		validData := map[string]interface{}{
			"field_one": 23,
		}

		t.Run("should return error when unable to get descriptor", func(t *testing.T) {
			result, err := client.Serialize("invalidClass", validData)
			assert.Nil(t, result)
			assert.Equal(t, stencil.ErrNotFound, err)
		})
		t.Run("should return error when unable to serialize to bytes", func(t *testing.T) {
			mapData := make(map[string]interface{})
			mapData["key1"] = "value1"

			result, err := client.Serialize("test.stencil.One", mapData)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Equal(t, stencil.ErrInvalidDescriptor, err)
		})
		t.Run("should return bytes", func(t *testing.T) {
			className := "test.stencil.One"
			bytes, err := client.Serialize(className, validData)
			assert.NoError(t, err)

			parsed, err := client.Parse(className, bytes)
			if err != nil {
				t.Fatal(err)
			}
			descriptor, err := client.GetDescriptor(className)
			if err != nil {
				t.Fatal(err)
			}
			fieldOneValue := validData["field_one"].(int)
			fieldOne := descriptor.Fields().ByName("field_one")
			val := parsed.ProtoReflect().Get(fieldOne)

			assert.Equal(t, int64(fieldOneValue), val.Int())
		})
	})
}

func TestRefreshStrategies(t *testing.T) {
	t.Run("VersionBasedRefresh", func(t *testing.T) {
		dataDownloadOneCount := 0
		dataDownloadTwoCount := 0
		versionsDownloadCount := 0
		// setup
		data, err := getDescriptorDataByPath(t, true, "./test_data")
		assert.NoError(t, err)
		versions := `{"versions": [1]}`
		mux := http.NewServeMux()
		mux.HandleFunc("/v1beta1/namespaces/test-namespace/schemas/test-schema/versions", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			rw.Write([]byte(versions))
			versionsDownloadCount++
		})
		mux.HandleFunc("/v1beta1/namespaces/test-namespace/schemas/test-schema/versions/1", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write(data)
			dataDownloadOneCount++
		})
		mux.HandleFunc("/v1beta1/namespaces/test-namespace/schemas/test-schema/versions/2", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write(data)
			dataDownloadTwoCount++
		})
		ts := httptest.NewServer(mux)

		// test
		opts := stencil.Options{AutoRefresh: true, RefreshStrategy: stencil.VersionBasedRefresh, RefreshInterval: 2 * time.Millisecond}
		client, err := stencil.NewClient([]string{fmt.Sprintf("%s/v1beta1/namespaces/test-namespace/schemas/test-schema", ts.URL)}, opts)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		// wait for refresh interval
		time.Sleep(3 * time.Millisecond)
		desc, err := client.GetDescriptor("test.stencil.One")
		assert.Nil(t, err)
		assert.NotNil(t, desc)
		time.Sleep(1 * time.Millisecond)
		assert.Equal(t, 2, versionsDownloadCount)
		assert.Equal(t, 1, dataDownloadOneCount)
		assert.Equal(t, 0, dataDownloadTwoCount)
		// simulates version update
		versions = `{"versions": [1,2]}`
		// wait for refresh interval
		time.Sleep(3 * time.Millisecond)
		desc, err = client.GetDescriptor("test.stencil.One")
		assert.Nil(t, err)
		assert.NotNil(t, desc)

		time.Sleep(1 * time.Millisecond)
		assert.Equal(t, 3, versionsDownloadCount)
		assert.Equal(t, 1, dataDownloadOneCount)
		assert.Equal(t, 1, dataDownloadTwoCount)
	})
}
