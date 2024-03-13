package changedetector

type ChangeRequest struct {
	NamespaceName string
	SchemaName    string
	Version       int32
	OldData       []byte
	NewData       []byte
}
