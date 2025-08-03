package processData

type DataProcessor interface {
	Process(data []byte, headers map[string]string) ([]byte, map[string]string, error)
}
