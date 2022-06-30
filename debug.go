package main

import (
	"fmt"
	"path"

	"github.com/aidenwang9867/scorecard-bigquery-auth/app/query"
)

func PrintDependencyToStdOut(d query.Dependency) {
	result := ""
	if d.IsDirect {
		result += fmt.Sprintf("name: %v \nstatus: %v \nversion: %v \necosys: %v \n",
			d.Name, d.ChangeType, d.Version, d.Ecosystem)
	} else {
		result += fmt.Sprintf("name: %v \nstatus: %v \nversion: %v \necosys: %v \n",
			d.Name, "indirect", d.Version, d.Ecosystem)
	}
	if d.PackageURL != "" {
		result += fmt.Sprintf("pkg_url: %v\n", d.PackageURL)
	}
	if d.SrcRepoURL != "" {
		result += fmt.Sprintf("src_url: %v\n", d.SrcRepoURL)
	}
	if len(d.Vulnerabilities) != 0 {
		result += fmt.Sprintln("vulnerabilities: ")
		for _, v := range d.Vulnerabilities {
			result += fmt.Sprintf("\tvuln_id: %v\n", v.ID)
			result += fmt.Sprintf("\tvuln_url: %v\n\n", v.SourceURL)
			if v.Title != "" {
				result += fmt.Sprintf("\tvuln_summary: %v\n", v.Title)
			}
			if v.GitHubSeverity != "" {
				result += fmt.Sprintf("\tseverity: %v\n", v.GitHubSeverity)
			}
			if v.Score.Valid {
				result += fmt.Sprintf("\tCVSS3Score: %v\n", v.Score.Float64)
			}
			if !v.DisclosedTime.IsZero() {
				result += fmt.Sprintf("\ttime_disclosed: %v\n", v.DisclosedTime)
			}
		}
	}
	fmt.Println(result)
}

func PrintDependencyChangeInfo(deps []query.Dependency) {
	removed, added := map[string]query.Dependency{}, map[string]query.Dependency{}
	for _, d := range deps {
		switch d.ChangeType {
		case query.Added:
			added[d.Name] = d
		case query.Removed:
			removed[d.Name] = d
		}
	}

	results := ""
	updated := map[string]query.Dependency{}
	for dName, d := range added {
		if _, ok := removed[dName]; ok {
			updated[dName] = d
		} else {
			current := fmt.Sprintf("**`" + "added" + "`** ")
			if len(d.Vulnerabilities) != 1 {
				current += fmt.Sprintf(createVulnAlert(d))
			}
			current += fmt.Sprintf(
				"%s: %s @ %s",
				d.Ecosystem, d.Name, d.Version,
			)
			results += current + "\n\n"
		}
	}
	for dName, d := range updated {
		current := fmt.Sprintf(
			"**`" + "updated" + "`**")
		if len(d.Vulnerabilities) != 1 {
			current += fmt.Sprintf(createVulnAlert(d))
		}
		current += fmt.Sprintf(
			" %s: %s @ %s (**old**) :arrow_right: %s @ %s @ %s (**new**)",
			added[dName].Ecosystem, added[dName].Name, added[dName].Version,
			d.Ecosystem, d.Name, d.Version,
		)
		results += current + "\n\n"
	}
	for dName, d := range removed {
		if _, ok := added[dName]; !ok {
			current := fmt.Sprintf(
				"~~**`"+"removed"+"`**~~"+" %s: %s @ %s",
				d.Ecosystem, d.Name, d.Version,
			)
			results += current + "\n\n"
		}
	}

	if results == "" {
		fmt.Println("No dependency changes")
	} else {
		fmt.Println(results)
	}
}

func createVulnAlert(d query.Dependency) string {
	system, name, version := "", d.Name, d.Version
	if d.Ecosystem == "gomod" {
		system = "go"
	} else if d.Ecosystem == "pip" {
		system = "pypi"
	}
	result := fmt.Sprintf(
		"[**`"+"vulnerable"+"`**](%s) ",
		"https://deps.dev/"+path.Join(system, name, version),
	)
	return result
}
