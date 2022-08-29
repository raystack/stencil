package schema

import "fmt"

func getNonEmpty(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}

func schemaKeyFunc(nsName, schema string, version int32) string {
	return fmt.Sprintf("%s-%s-%d", nsName, schema, version)
}

func getBytes(key interface{}) []byte {
	buf, _ := key.([]byte)
	return buf
}
