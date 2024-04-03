package database

type Application struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Logo                 string `json:"logo"`
	Slug                 string `json:"slug"`
	HomepageTemplate     string `json:"homepageTemplate"`
	TrackpageTemplate    string `json:"trackpageTemplate"`
	AdditionalCss        string `json:"additionalCss,omitempty"`
	AdditionalJavaScript string `json:"additionalJavaScript,omitempty"`
}
