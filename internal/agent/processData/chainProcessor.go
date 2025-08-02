package processData

type ChainProcessor struct {
	Processors []DataProcessor
}

func (p *ChainProcessor) Process(data []byte, headers map[string]string) ([]byte, map[string]string, error) {
	var err error
	for _, processor := range p.Processors {
		data, headers, err = processor.Process(data, headers)
		if err != nil {
			return nil, nil, err
		}
	}
	return data, headers, nil
}
