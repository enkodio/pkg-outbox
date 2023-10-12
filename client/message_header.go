package client

type MessageHeader struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

func (m MessageHeader) GetKey() string {
	return m.Key
}

func (m MessageHeader) GetValue() []byte {
	return m.Value
}

type MessageHeaders []MessageHeader

func NewMessageHeaders(headers Headers) MessageHeaders {
	messageHeaders := make(MessageHeaders, 0, len(headers))
	for _, h := range headers {
		if h == nil {
			continue
		}
		messageHeaders = append(messageHeaders, MessageHeader{
			Key:   h.GetKey(),
			Value: h.GetValue(),
		})
	}
	return messageHeaders
}

func (m MessageHeaders) ToHeaders() Headers {
	headers := make(Headers, len(m))
	for i, h := range m {
		headers[i] = h
	}
	return headers
}

func (m MessageHeaders) GetValueByKey(key string) []byte {
	for _, header := range m {
		if header.GetKey() == key {
			return header.GetValue()
		}
	}
	return nil
}

func (m *MessageHeaders) SetHeader(key string, value []byte) {
	*m = append(*m, MessageHeader{
		Key:   key,
		Value: value,
	})
}
