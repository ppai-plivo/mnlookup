package api

type Record struct {
	Prefix string `json:"prefix"`
	MCC    string `json:"mcc"`
	MNC    string `json:"mnc"`
}
