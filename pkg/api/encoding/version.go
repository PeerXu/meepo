package encoding_api

type Version struct {
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
	GitHash   string `json:"gitHash"`
	Built     string `json:"built"`
	Platform  string `json:"platform"`
}
