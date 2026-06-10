package response

type Body struct {
	Data any `json:"data,omitempty"`
}

func OK(data any) Body {
	return Body{Data: data}
}
