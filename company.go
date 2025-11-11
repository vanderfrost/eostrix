package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CompanyProblem struct {
	Company        string
	Difficulty     string
	Title          string
	Frequency      string
	AcceptanceRate string
	Link           string
}

// search each company folder for the correct six month cvs file
func findSixMonthCSV(companyDir string) (string, error) {
	entries, err := os.ReadDir(companyDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(e.Name(), "3. Six Months.csv") {
			return filepath.Join(companyDir, e.Name()), nil
		}
	}
	return "", fmt.Errorf("a six month csv not found in %s", companyDir)
}

// load the leetcode problems from the six month cvs files to the CompanyProblem struct
func loadAllCompanyProblems(rootDir string) ([]CompanyProblem, error) {
	var allProblems []CompanyProblem

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		companyName := entry.Name()
		companyDir := filepath.Join(rootDir, companyName)

		csvPath, err := findSixMonthCSV(companyDir)
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", companyName, err)
			continue
		}

		f, err := os.Open(csvPath)
		if err != nil {
			fmt.Printf("Failed to open %s: %v\n", csvPath, err)
			continue
		}

		r := csv.NewReader(f)

		if _, err := r.Read(); err != nil {
			f.Close()
			fmt.Printf("Failed to read header of %s: %v\n", csvPath, err)
			continue
		}

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("error reading %s: %v\n", csvPath, err)
				continue
			}

			if len(record) < 5 {
				fmt.Printf("skipping bad row in %s: %v\n", csvPath, record)
				continue
			}

			allProblems = append(allProblems, CompanyProblem{
				Company:        companyName,
				Difficulty:     record[0],
				Title:          record[1],
				Frequency:      record[2],
				AcceptanceRate: record[3],
				Link:           record[4],
			})
		}

		f.Close()
	}

	fmt.Println(allProblems)

	return allProblems, nil
}
