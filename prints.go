package main

import (
	"fmt"
)

func PrintDependencies(deps []dependency) {
	for _, d := range deps {
		fmt.Println(*d.Ecosystem, d.Name, *d.Version, *d.ChangeType)
		if d.PackageURL != nil && *d.PackageURL != "" {
			fmt.Println(*d.PackageURL)
		} else {
			fmt.Println("empty package url")
		}
		if d.SourceRepository != nil && *d.SourceRepository != "" {
			fmt.Println(*d.SourceRepository)
		} else {
			fmt.Println("empty src url")
		}
		fmt.Printf("===================\n\n")
	}
}
