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
	client, err := stencil.NewClient(url, stencil.Options{})
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
	t.Run("ParseWithRefresh", func(t *testing.T) {

		t.Run("should return notFoundErr if not found", func(t *testing.T) {
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
			msg, err := client.ParseWithRefresh("test.stencil.Two.Unknown", []byte(""))
			assert.Nil(t, msg)
			assert.NotNil(t, err)
			assert.Equal(t, stencil.ErrNotFound, err)
		})

		t.Run("should return error if refresh fails", func(t *testing.T) {
			_, msgData := getUpdatedDescriptorDataAndMsgData(t, true)
			data, err := getDescriptorData(t, true)
			assert.NoError(t, err)
			count := 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if count == 1 {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					count++
					return
				}
				count++
				w.Write(data)
			}))
			defer ts.Close()
			url := ts.URL
			client, err := stencil.NewClient(url, stencil.Options{})
			assert.Nil(t, err)
			assert.NotNil(t, client)
			_, err = client.ParseWithRefresh("test.stencil.One", msgData)
			assert.NotNil(t, err)
		})

		t.Run("should refresh proto definitions should parse without any UnknownFields", func(t *testing.T) {
			data, msgData := getUpdatedDescriptorDataAndMsgData(t, true)
			oldData, err := getDescriptorData(t, true)
			count := 0
			ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if count == 0 {
					w.Write(oldData)
					count++
					return
				}
				w.Write(data)
			}))
			defer ts2.Close()
			newClient, err := stencil.NewClient(ts2.URL, stencil.Options{})
			parsed, err := newClient.Parse("test.stencil.One", msgData)
			assert.NotNil(t, parsed.ProtoReflect().GetUnknown())
			parsed, err = newClient.ParseWithRefresh("test.stencil.One", msgData)
			assert.Nil(t, err)
			assert.Nil(t, parsed.ProtoReflect().GetUnknown())
			parsed.ProtoReflect().Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
				if field.Name() == "field_one" {
					assert.Equal(t, int64(200), value.Int())
				}
				if field.Name() == "field_two" {
					assert.Equal(t, int64(300), value.Int())
				}
				return true
			})
		})
	})

	t.Run("Serialize", func(t *testing.T) {
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

		validData := map[string]interface{}{
			"field_one": 23,
		}

		t.Run("should return error when unable to get descriptor", func(t *testing.T) {
			result, err := client.Serialize("invalidClass", validData)
			assert.Nil(t, result)
			assert.Equal(t, stencil.ErrNotFound, err)
		})
		t.Run("should return error when unable to serialize to bytes", func(t *testing.T) {
			className := "test.stencil.One"

			mapData := make(map[string]interface{})
			mapData["key1"] = "value1"

			result, err := client.Serialize(className, mapData)
			assert.Nil(t, result)
			assert.Error(t, err)
		})
		t.Run("should return bytes", func(t *testing.T) {
			className := "test.stencil.One"

			result, err := client.Serialize(className, validData)
			assert.NoError(t, err)

			expected := []byte{0x8, 0x17}
			assert.Equal(t, expected, result)
		})
	})
}
