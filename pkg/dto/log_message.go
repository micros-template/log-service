package dto

type LogMessage struct {
	Type    string `json:"type"`
	Service string `json:"service"`
	Msg     string `json:"msg"`
}
