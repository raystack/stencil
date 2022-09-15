package cmd

var dict = map[string]string{
	"COMPATIBILITY_BACKWARD":    "backward",
	"COMPATIBILITY_FORWARD":     "forward",
	"COMPATIBILITY_FULL":        "full",
	"COMPATIBILITY_UNSPECIFIED": "-",
	"FORMAT_PROTOBUF":           "protobuf",
	"FORMAT_JSON":               "json",
	"FORMAT_AVRO":               "avro",
}

var (
	formats = []string{
		"FORMAT_JSON",
		"FORMAT_PROTOBUF",
		"FORMAT_AVRO",
	}

	comps = []string{
		"COMPATIBILITY_BACKWARD",
		"COMPATIBILITY_FORWARD",
		"COMPATIBILITY_FULL",
	}
)
