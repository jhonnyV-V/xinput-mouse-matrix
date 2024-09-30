package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var Quit bool

func listItems() ([]list.Item, error) {
	items := []list.Item{}
	var buff bytes.Buffer
	cmd := exec.Command("xinput", "list", "--name-only")
	cmd.Stdout = &buff
	err := cmd.Run()

	rawOutput := buff.String()

	devices := strings.Split(rawOutput, "\n")

	for _, device := range devices {
		items = append(items, item(device))
	}

	return items, err
}

func getCurrentValue(deviceName string) (float32, float32, float32, error) {
	var x, y, acceleration float32
	var values []string

	cmd := exec.Command("xinput", "--list-props", deviceName)
	var buff bytes.Buffer
	cmd.Stdout = &buff
	err := cmd.Run()

	b := buff.String()
	props := strings.Split(b, "\n")[1:]

	for _, prop := range props {
		clean := strings.TrimSpace(prop)
		if !strings.Contains(clean, "Coordinate Transformation Matrix") {
			continue
		}

		values = strings.Split(
			strings.TrimSpace(strings.Split(clean, ":")[1]),
			",",
		)
	}

	temp, err := strconv.ParseFloat(
		strings.TrimSpace(values[0]),
		32,
	)
	if err != nil {
		return x, y, acceleration, err
	}
	x = float32(temp)

	temp, err = strconv.ParseFloat(
		strings.TrimSpace(values[4]),
		32,
	)
	if err != nil {
		return x, y, acceleration, err
	}
	y = float32(temp)

	temp, err = strconv.ParseFloat(
		strings.TrimSpace(values[8]),
		32,
	)
	if err != nil {
		return x, y, acceleration, err
	}
	acceleration = float32(temp)

	return x, y, acceleration, err
}

func setValue(x, y, acceleration float32, deviceName string) error {
	cmd := exec.Command(
		"xinput", "set-prop",
		deviceName,
		"Coordinate Transformation Matrix",
		fmt.Sprintf("%f,", x), "0,", "0,",
		"0,", fmt.Sprintf("%f,", y), "0,",
		"0,", "0,", fmt.Sprintf("%f,", acceleration),
	)
	err := cmd.Run()
	return err
}

func isValidFloat32(s string) (float32, bool) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err == nil
}

func main() {
	items, err := listItems()

	if err != nil {
		fmt.Printf("error %s\n", err)
		return
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if Quit {
		os.Exit(0)
	}

	x, y, acceleration, err := getCurrentValue(Choice)
	if err != nil {
		fmt.Printf("error getting current values %s\n", err)
		return
	}

	inputs := initialInputModels(x, y, acceleration)

	if _, err := tea.NewProgram(inputs).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	if Quit {
		os.Exit(0)
	}

	validationError := false
	var f float32
	var ok bool
	input := inputs.inputs[0]
	v := input.Value()
	if v != "" {
		f, ok = isValidFloat32(v)
		if ok {
			x = f
		} else {
			validationError = true
			fmt.Printf("Invalid Value \"%s\" for x\n", v)
		}
	}

	input = inputs.inputs[1]
	v = input.Value()
	if v != "" {
		f, ok = isValidFloat32(v)
		if ok {
			y = f
		} else {
			validationError = true
			fmt.Printf("Invalid Value \"%s\" for y\n", v)
		}
	}

	input = inputs.inputs[2]
	v = input.Value()
	if v != "" {
		f, ok = isValidFloat32(v)
		if ok {
			acceleration = f
		} else {
			validationError = true
			fmt.Printf("Invalid Value \"%s\" for acceleration\n", v)
		}
	}

	if validationError {
		os.Exit(2)
	}

	err = setValue(x, y, acceleration, Choice)
	if err != nil {
		fmt.Printf("error setting new values for device: %s\n", err)
	}
}
