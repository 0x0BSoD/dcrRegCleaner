package main

type Config struct {
	Registry   string `yaml:"registry"`
	Service    string `yaml:"service"`
	TagsToStay int    `yaml:"tags_to_stay"`
	ApiVersion string `yaml:"api_ver"`
}

type ToPrint []ToPrintItem

type ToPrintItem struct {
	Repo   string
	Tag    string
	Digest string
	Size   int
}

type ApiRepos struct {
	Repositories []string `json:"repositories"`
}

type ApiTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ApiManifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}
