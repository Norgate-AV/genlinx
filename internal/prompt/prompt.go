package prompt

import (
	"errors"
	"os"

	"github.com/charmbracelet/huh"
)

// Confirm presents a yes/no confirmation prompt with the given title and
// description. It returns true when the user confirms, false when they decline.
// A ctrl+c / ErrUserAborted causes the process to exit cleanly.
func Confirm(title, description string) (bool, error) {
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Value(&confirmed),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			os.Exit(0)
		}

		return false, err
	}

	return confirmed, nil
}

// SelectFiles presents an interactive multi-select prompt and returns the chosen
// files — arrow keys navigate, x toggles, enter confirms, ctrl+a selects all / none.
// At least one file must be selected.
func SelectFiles(files []string) ([]string, error) {
	options := make([]huh.Option[string], len(files))
	for i, f := range files {
		options[i] = huh.NewOption(f, f)
	}

	var selected []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select file(s):").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			os.Exit(0)
		}

		return nil, err
	}

	if len(selected) == 0 {
		return nil, errors.New("you must choose at least one file")
	}

	return selected, nil
}
