package dto

type VersionDTO struct {
	Name        string   `json:"name"`
	HomepageURL string   `json:"homepage_url"`
	Version     string   `json:"version"`
	PrevVersion *string  `json:"last_version"`
	NextVersion *string  `json:"next_version"`
	RemoteDate  *DateDTO `json:"remote_date,omitempty"`
	UpdateDate  *DateDTO `json:"update_date"`
	Link        *string  `json:"link,omitempty"`
	Digest      *string  `json:"digest,omitempty"`
	Picture     *string  `json:"picture,omitempty"`
}
