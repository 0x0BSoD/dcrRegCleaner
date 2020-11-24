package main

import (
	"encoding/json"
	"fmt"
	"time"
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

func shaDigestAndSize(c *Client, repo, tag string) (string, time.Time, error) {
	bData, headers, err := c.Get(fmt.Sprintf("/%s/manifests/%s", repo, tag))
	if err != nil {
		return "", time.Time{}, err
	}

	digest := headers.Get("Docker-Content-Digest")

	var res ApiManifest
	err = json.Unmarshal(bData, &res)
	if err != nil {
		return "", time.Time{}, err
	}

	var containerInfo ApiManifestCommon
	err = json.Unmarshal([]byte(res.History[0].V1Compatibility), &containerInfo)
	if err != nil {
		return "", time.Time{}, err
	}

	// container layers
	//for _, v := range res.History[1:] {
	//	var cc ApiManifestContainerConfig
	//	err = json.Unmarshal([]byte(v.V1Compatibility), &cc)
	//	if err != nil {
	//		return "", 0, err
	//	}
	//
	//	fmt.Println(v.V1Compatibility)
	//	fmt.Println()
	//}

	return digest, containerInfo.Created, nil
}

func removeImage(c *Client, repo, tag string) error {
	_, _, err := c.Delete(fmt.Sprintf("/%s/%s", repo, tag))
	if err != nil {
		return err
	}

	return nil
}

//func removeOldImages(c *)

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
			digest, created, err := shaDigestAndSize(c, r, t)
			if err != nil {
				return err
			}

			toPrint = append(toPrint, ToPrintItem{
				Created: created,
				Repo:    r,
				Digest:  digest,
				Size:    0,
				Tag:     t,
			})
		}
	}

	printAsTable(toPrint)

	return nil
}
