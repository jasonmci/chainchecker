package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LanePairs []string
}

func generateMarkdownTable(headers []string, rows [][]string) string {
	table := "| "
	for _, header := range headers {
		table += header + " | "
	}
	table += "\n"

	table += "|"
	for range headers {
		table += " --- |"
	}
	table += "\n"

	for _, row := range rows {
		table += "| "
		for _, col := range row {
			table += col + " | "
		}
		// Fill in empty cells if row has fewer columns than headers
		for i := len(row); i < len(headers); i++ {
			table += " | "
		}
		table += "\n"
	}

	return table
}

func main() {
	config := Config{}
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	fmt.Println("Enter lane pairs (comma-separated, e.g., 'LaneA-LaneB,LaneX-LaneY'):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var lanePairs []string
	if input == "" {
		lanePairs = config.LanePairs // Use default lane pairs from config
	} else {
		lanePairs = strings.Split(input, ",")
	}

	rows := [][]string{}
	for _, pair := range lanePairs {
		lanes := strings.Split(pair, "-")
		if len(lanes) == 2 {
			rows = append(rows, []string{pair, lanes[0], lanes[1]})
			rows = append(rows, []string{pair, lanes[1], lanes[0]})
		}
	}

	headers := []string{"Lane Combo", "Source", "Destination", "Scenario", "Status", "Transactions"}
	tableMarkdown := generateMarkdownTable(headers, rows)
	fmt.Println(tableMarkdown)

	err = os.WriteFile("table.md", []byte(tableMarkdown), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
