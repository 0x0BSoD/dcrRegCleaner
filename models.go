package main

import "time"

type Config struct {
	Registry   string `yaml:"registry"`
	Service    string `yaml:"service"`
	TagsToStay int    `yaml:"tags_to_stay"`
	ApiVersion string `yaml:"api_ver"`
}

type ToPrint []ToPrintItem

func (t ToPrint) Len() int {
	return len(t)
}

func (t ToPrint) Less(i, j int) bool {
	return t[i].Created.Before(t[j].Created)
}

func (t ToPrint) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type ToPrintItem struct {
	Created time.Time
	Repo    string
	Digest  string
	Size    int
	Tag     string
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
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
	Signatures []struct {
		Header struct {
			Jwk struct {
				Crv string `json:"crv"`
				Kid string `json:"kid"`
				Kty string `json:"kty"`
				X   string `json:"x"`
				Y   string `json:"y"`
			} `json:"jwk"`
			Alg string `json:"alg"`
		} `json:"header"`
		Signature string `json:"signature"`
		Protected string `json:"protected"`
	} `json:"signatures"`
}

type ApiManifestContainerConfig struct {
	ID              string    `json:"id"`
	Parent          string    `json:"parent"`
	Created         time.Time `json:"created"`
	ContainerConfig struct {
		Cmd []string `json:"Cmd"`
	} `json:"container_config"`
}

type ApiManifestCommon struct {
	Architecture string `json:"architecture"`
	Config       struct {
		Hostname     string `json:"Hostname"`
		Domainname   string `json:"Domainname"`
		User         string `json:"User"`
		AttachStdin  bool   `json:"AttachStdin"`
		AttachStdout bool   `json:"AttachStdout"`
		AttachStderr bool   `json:"AttachStderr"`
		ExposedPorts struct {
			Eight0TCP struct {
			} `json:"80/tcp"`
		} `json:"ExposedPorts"`
		Tty        bool        `json:"Tty"`
		OpenStdin  bool        `json:"OpenStdin"`
		StdinOnce  bool        `json:"StdinOnce"`
		Env        []string    `json:"Env"`
		Cmd        []string    `json:"Cmd"`
		Image      string      `json:"Image"`
		Volumes    interface{} `json:"Volumes"`
		WorkingDir string      `json:"WorkingDir"`
		Entrypoint []string    `json:"Entrypoint"`
		OnBuild    interface{} `json:"OnBuild"`
		Labels     struct {
			Maintainer string `json:"maintainer"`
		} `json:"Labels"`
		StopSignal string `json:"StopSignal"`
	} `json:"config"`
	Container       string `json:"container"`
	ContainerConfig struct {
		Hostname     string `json:"Hostname"`
		Domainname   string `json:"Domainname"`
		User         string `json:"User"`
		AttachStdin  bool   `json:"AttachStdin"`
		AttachStdout bool   `json:"AttachStdout"`
		AttachStderr bool   `json:"AttachStderr"`
		ExposedPorts struct {
			Eight0TCP struct {
			} `json:"80/tcp"`
		} `json:"ExposedPorts"`
		Tty        bool        `json:"Tty"`
		OpenStdin  bool        `json:"OpenStdin"`
		StdinOnce  bool        `json:"StdinOnce"`
		Env        []string    `json:"Env"`
		Cmd        []string    `json:"Cmd"`
		Image      string      `json:"Image"`
		Volumes    interface{} `json:"Volumes"`
		WorkingDir string      `json:"WorkingDir"`
		Entrypoint []string    `json:"Entrypoint"`
		OnBuild    interface{} `json:"OnBuild"`
		Labels     struct {
			Maintainer string `json:"maintainer"`
		} `json:"Labels"`
		StopSignal string `json:"StopSignal"`
	} `json:"container_config"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	ID            string    `json:"id"`
	Os            string    `json:"os"`
	Parent        string    `json:"parent"`
	Throwaway     bool      `json:"throwaway"`
}
