package ymusic

type InvocationInfo struct {
	ExecDurationMillis any    `json:"exec-duration-millis"`
	Hostname           string `json:"hostname"`
	ReqID              string `json:"req-id"`
}
