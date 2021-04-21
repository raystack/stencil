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

	stencil "github.com/odpf/stencil/clients/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func runProtoc(
	rootDir string,
	includeImports bool,
	descSetOut string,
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
	protoFiles, _ := filepath.Glob(filepath.Join(rootDir, "./**/*.proto"))
	rootFiles, _ := filepath.Glob(filepath.Join(rootDir, "./*.proto"))
	args = append(args, rootFiles...)
	args = append(args, protoFiles...)
	stderr := bytes.NewBuffer(nil)
	cmd := exec.Command(protocBinPath, args...)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s returned error: %v %v", protocBinPath, err, stderr.String())
	}
	return nil
}

func getDescriptorData(t *testing.T, includeImports bool) ([]byte, error) {
	root, _ := filepath.Abs("./test_data")
	fileName := filepath.Join(t.TempDir(), "file.desc")
	err := runProtoc(root, includeImports, fileName)
	assert.NoError(t, err)
	data, err := ioutil.ReadFile(fileName)
	return data, err
}

func TestNewClient(t *testing.T) {
	t.Run("should return error if url is not valid", func(t *testing.T) {
		url := "h_ttp://invalidurl"
		_, err := stencil.NewClient(url, stencil.Options{})
		assert.Contains(t, err.Error(), "invalid request")
	})

	t.Run("should return error if request fails", func(t *testing.T) {
		url := "ithttp://localhost"
		_, err := stencil.NewClient(url, stencil.Options{})
		assert.Contains(t, err.Error(), "request failed")
	})

	t.Run("should return error if file download fails", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}))
		defer ts.Close()
		url := ts.URL
		_, err := stencil.NewClient(url, stencil.Options{})
		assert.Contains(t, err.Error(), "request failed.")
	})

	t.Run("should return error if downloaded file is not valid", func(t *testing.T) {
		data := []byte("invalid")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		_, err := stencil.NewClient(url, stencil.Options{})
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
		_, err = stencil.NewClient(url, stencil.Options{})
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
		client, err := stencil.NewClient(url, stencil.Options{})
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
		client, err := stencil.NewClient(ts.URL, stencil.Options{HTTPOptions: stencil.HTTPOptions{Headers: headers}})
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
		client, _ := stencil.NewClient(ts.URL, stencil.Options{AutoRefresh: true, RefreshInterval: 5 * time.Millisecond})
		time.Sleep(6 * time.Millisecond)
		client.Close()
		assert.Equal(t, 2, callCount)
	})
}

func TestNewMultiURLClient(t *testing.T) {
	data, err := getDescriptorData(t, true)
	assert.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer ts.Close()
	url := ts.URL
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	defer ts2.Close()
	url2 := ts2.URL
	_, err = stencil.NewMultiURLClient([]string{url, url2}, stencil.Options{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request failed.")
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
		client, err := stencil.NewClient(url, stencil.Options{})
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
	})
	t.Run("Parse", func(t *testing.T) {
		data, err := getDescriptorData(t, true)
		assert.NoError(t, err)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(data)
		}))
		defer ts.Close()
		url := ts.URL
		client, err := stencil.NewClient(url, stencil.Options{})
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
}
