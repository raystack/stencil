package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type Prompter interface {
	Select(string, string, []string) (int, error)
	MultiSelect(string, string, []string) (int, error)
	Input(string, string) (string, error)
	Confirm(string, bool) (bool, error)
}

func New() Prompter {
	return &surveyPrompter{}
}

type surveyPrompter struct {
}

func (p *surveyPrompter) ask(q survey.Prompt, response interface{}) error {
	err := survey.AskOne(q, response)
	if err == nil {
		return nil
	}
	return fmt.Errorf("could not prompt: %w", err)
}

func (p *surveyPrompter) Select(message, defaultValue string, options []string) (result int, err error) {
	q := &survey.Select{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}

	err = p.ask(q, &result)

	return
}

func (p *surveyPrompter) MultiSelect(message, defaultValue string, options []string) (result int, err error) {
	q := &survey.MultiSelect{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}

	err = p.ask(q, &result)

	return
}

func (p *surveyPrompter) Input(prompt, defaultValue string) (result string, err error) {
	err = p.ask(&survey.Input{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}

func (p *surveyPrompter) Confirm(prompt string, defaultValue bool) (result bool, err error) {
	err = p.ask(&survey.Confirm{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}
