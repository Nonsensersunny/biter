package model

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	ClientIp  string `json:"client_ip"`
}

type ActionResponse struct {
	Res      string      `json:"res"`
	Error    string      `json:"error"`
	Ecode    interface{} `json:"ecode"`
	ErrorMsg string      `json:"error_msg"`
	ClientIp string      `json:"client_ip"`
}