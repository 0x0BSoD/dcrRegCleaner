package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func parseConfig(path string) (Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return Config{}, err
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	return config, nil
}


func printAsTable(toPrint ToPrint) {

	data := make([][]string, 0, len(toPrint))

	for _, c := range toPrint {
		//t := strconv.Itoa(c.Size / 1024 )
		data = append(data, []string{c.Repo, c.Tag, c.Digest})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repo", "Tag", "Digest"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoMergeCells(true)
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}