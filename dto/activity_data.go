package dto

type ActivityData struct {
	Invitor      string `json:"invitor"`
	Sign         string `json:"sign"`
	ActivityTime int64  `json:"activity_time"`
	TryQQ        string `json:"try_qq"`
}
