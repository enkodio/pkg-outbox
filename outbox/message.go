package outbox

type Message struct {
	Headers MessageHeaders `json:"headers"`
	Body    interface{}    `json:"body"`
	Topic   string         `json:"topic"`
}

type MessageHeaders []MessageHeader

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

func (m *MessageHeaders) ToMap() map[string][]byte {
	headers := make(map[string][]byte, len(*m))
	for _, h := range *m {
		headers[h.Key] = h.Value
	}
	return headers
}

func (m *MessageHeaders) GetValueByKey(key string) []byte {
	for _, header := range *m {
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

func NewMessage(topic string, body interface{}, headers ...map[string][]byte) Message {
	return Message{
		Topic:   topic,
		Body:    body,
		Headers: NewMessageHeaders(headers...),
	}
}
func NewMessageHeaders(headers ...map[string][]byte) MessageHeaders {
	messageHeaders := make(MessageHeaders, 0, len(headers))
	for _, h := range headers {
		for key, value := range h {
			messageHeaders = append(messageHeaders, MessageHeader{
				Key:   key,
				Value: value,
			})
		}
	}
	return messageHeaders
}
