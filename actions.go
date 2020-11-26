package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

func removeImage(c *Client, repo, reference string) error {
	_, _, err := c.Delete(fmt.Sprintf("/%s/manifests/%s", repo, reference))
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

	toPrint := make(ToPrint)

	for _, r := range repos.Repositories {
		tags, err := tags(c, r)
		if err != nil {
			return err
		}

		for _, t := range tags.Tags {
			digest, created, err := shaDigestAndSize(c, r, t)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			toPrint[r] = append(toPrint[r], ToPrintItem{
				Created: created,
				Repo:    r,
				Digest:  digest,
				Size:    0,
				Tag:     t,
			})
		}
	}

	data := make([][]string, 0, len(toPrint))

	for _, i := range toPrint {
		sort.Sort(i)
		for _, c := range i {
			data = append(data, []string{c.Repo, c.Tag, c.Created.Format("2006-01-02"), c.Digest})
		}
	}

	printAsTable(data)

	return nil
}

func removeOldImages(c *Client) error {
	repos, err := catalog(c)
	if err != nil {
		return err
	}

	toPrint := make(ToPrint)

	for _, r := range repos.Repositories {
		tags, err := tags(c, r)
		if err != nil {
			return err
		}

		for _, t := range tags.Tags {
			digest, created, err := shaDigestAndSize(c, r, t)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			toPrint[r] = append(toPrint[r], ToPrintItem{
				Created: created,
				Repo:    r,
				Digest:  digest,
				Size:    0,
				Tag:     t,
			})
		}
	}

	var result [][]string
	var toRemove [][]string

	for _, i := range toPrint {
		sort.Sort(i)
		var data [][]string
		for j := len(i) - 1; j >= 0; j-- {
			if len(data) < c.TagsToStay {
				data = append(data, []string{i[j].Repo, i[j].Tag, i[j].Created.Format("2006-01-02"), i[j].Digest})
			} else {
				data = append(data, []string{"[TO REMOVE] " + i[j].Repo, i[j].Tag, i[j].Created.Format("2006-01-02"), i[j].Digest})
				toRemove = append(toRemove, []string{i[j].Digest, i[j].Repo, i[j].Tag})
			}
		}
		result = append(result, data...)
	}
	if c.RmCheck {
		printAsTable(result)
	} else {
		for _, t := range toRemove {
			fmt.Printf("[X] Removing %s:%s ", t[1], t[2])
			err := removeImage(c, t[1], t[0])
			if err != nil {
				fmt.Println("-- failed")
			} else {
				fmt.Println("-- done")
			}
		}
	}

	return nil
}
