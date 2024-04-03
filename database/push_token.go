package database

type PushToken struct {
	Id            string   `json:"id"`
	Token         string   `json:"token"`
	AllowedApps   []string `json:"allowedApps,omitempty"`
	AllowedTracks []string `json:"allowedTracks,omitempty"`
}
