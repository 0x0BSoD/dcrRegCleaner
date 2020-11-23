package main

import (
	"encoding/json"
	"fmt"
)

func catalog(c *Client) (ApiRepos, error) {
	bData, _, err := c.Get("/_catalog")
	if err != nil {
		return ApiRepos{}, err
	}

	var res ApiRepos
	err = json.Unmarshal(bData, &res)
	if err != nil {
		return ApiRepos{}, err
	}

	return res, nil
}

func tags(c *Client, repo string) (ApiTags, error) {
	bData, _, err := c.Get(fmt.Sprintf("/%s/tags/list", repo))
	if err != nil {
		return ApiTags{}, err
	}

	var res ApiTags
	err = json.Unmarshal(bData, &res)
	if err != nil {
		return ApiTags{}, err
	}

	return res, nil
}

func shaDigestAndSize(c *Client, repo, tag string) (string, int, error) {
	bData, headers, err := c.Get(fmt.Sprintf("/%s/manifests/%s", repo, tag))
	if err != nil {
		return "", 0, err
	}

	digest := headers.Get("Docker-Content-Digest")

	var res ApiManifest
	err = json.Unmarshal(bData, &res)
	if err != nil {
		return "", 0, err
	}

	var size int
	for _, c := range res.Layers {
		size += c.Size
	}

	return digest, size + res.Config.Size, nil
}

func removeImage(c *Client, repo, tag string) error {
	_, _, err := c.Delete(fmt.Sprintf("/%s/%s", repo, tag))
	if err != nil {
		return err
	}

	return nil
}

func printRegContent(c *Client) error {
	repos, err := catalog(c)
	if err != nil {
		return err
	}

	var toPrint ToPrint

	for _, r := range repos.Repositories {
		tags, err := tags(c, r)
		if err != nil {
			return err
		}

		for _, t := range tags.Tags {
			d, s, err := shaDigestAndSize(c, r, t)
			if err != nil {
				return err
			}

			toPrint = append(toPrint, ToPrintItem{
				Repo:   r,
				Tag:    t,
				Digest: d,
				Size:   s,
			})
		}
	}

	printAsTable(toPrint)

	return nil
}
