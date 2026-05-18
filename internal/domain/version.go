package domain

type VersionInfo struct {
	Version string `json:"version" yaml:"version"`
	Commit  string `json:"commit"  yaml:"commit"`
	Built   string `json:"built"   yaml:"built"`
	Go      string `json:"go"      yaml:"go"`
}
