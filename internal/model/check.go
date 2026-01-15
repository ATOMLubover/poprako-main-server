package model

type CheckVersionReply struct {
	LatestVersion string `json:"latest_version"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	AllowUsage    bool   `json:"allow_usage"`
}
