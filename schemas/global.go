package schemas

type Ack struct {
	Value string `json:"ack"`
}

type NoArg = interface{}

type IdValue struct {
	Value string `param:"id" json:"id"`
}

type CountModel struct {
	Count int64 `json:"count"`
}

type Amount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type KeyValue struct {
	Key   string  `json:"key"`
	Value float32 `json:"value"`
}

type H map[string]interface{}
