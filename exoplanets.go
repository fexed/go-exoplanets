package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

type exoplanet struct { // https://exoplanetarchive.ipac.caltech.edu/cgi-bin/nstedAPI/nph-nstedAPI?table=cumulative
	id     string // 0
	name   string // 2
	status string // 3
	score  string // 5
	period string // 10
}

func getData(url string) ([][]string, error) {
	fmt.Println("Downloading data")
	response, error := http.Get(url)

	fmt.Println("Data downloaded")
	if error != nil {
		return nil, fmt.Errorf("GET error: %v", error)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", error)
	}

	fmt.Println("Parsing data")
	csvLines, error := csv.NewReader(response.Body).ReadAll()
	if error != nil {
		return nil, fmt.Errorf("Error reading body: %v", error)
	}

	return csvLines, nil
}

func main() {
	fmt.Println("NASA Exoplanets explorer")

	csvLines, error := getData("https://exoplanetarchive.ipac.caltech.edu/cgi-bin/nstedAPI/nph-nstedAPI?table=cumulative")
	if error != nil {
		log.Printf("Failed to get data: %v", error)
	} else {
		var exoplanets []exoplanet
		var legend = csvLines[0]
		for i := 0; i < len(legend); i++ {
			fmt.Printf("%d. %s\n", i, legend[i])
		}
		for _, line := range csvLines {
			if line[3] == "CONFIRMED" { // only confirmed exoplanets
				exp := exoplanet{
					id:     line[0],
					name:   line[2],
					status: line[3],
					score:  line[5],
					period: line[10],
				}
				exoplanets = append(exoplanets, exp)
			}
		}
		fmt.Printf("Obtained data about %d confirmed exoplanets\n", len(exoplanets))
		for i := 0; i < 10; i++ {
			fmt.Printf("%s.\t%s\t(score: %s)\t%s days\n", exoplanets[i].id, exoplanets[i].name, exoplanets[i].score, exoplanets[i].period)
		}
	}
}
