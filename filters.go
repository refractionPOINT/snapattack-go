package snapattack

type Filter struct {
	Operator string      `json:"op"`
	Items    []*Filter   `json:"items,omitempty"`
	Field    string      `json:"field,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}
