package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"colormind/errors"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

// modelFlag implements flagSet interface to create flag set of model values to limit
// number of available model values to modelsList.
type modelFlag string

var model modelFlag = "default"

func (m *modelFlag) String() string {
	return string(*m)
}

func (m *modelFlag) Set(s string) error {
	for _, val := range models.List {
		if s == val {
			*m = modelFlag(s)
			return nil
		}
	}
	models := strings.Join(models.List, ", ")
	err := fmt.Sprintf("must be one of %s", models)
	return fmt.Errorf("unavailable model value: %s", err)
}

func (m *modelFlag) Type() string {
	return "modelFlag"
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.PersistentFlags().VarP(&model, "model", "m", "a model for color scheme.")
	randomCmd.PersistentFlags().BoolP("hexadecimal", "H", false, "entering and showing colors in hexadecimal")
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random color scheme",
	Long:  `This command fetches a random color scheme from the colormind api`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		isHex, err := cmd.Flags().GetBool("hexadecimal")
		if err != nil {
			errors.HandleError(errors.InternalErr, err)
			return
		}
		RandomScheme(model.String(), isHex)
	},
}

// RandomScheme gets random color scheme depending on model flag's value
// by getData function and converts it from json to scheme type.
func RandomScheme(model string, isHex bool) {
	param := input{Model: model}
	schemeData, err := getData(param)
	if err != nil {
		err = fmt.Errorf("couldn't get random scheme data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	sentry.CaptureMessage("Random scheme data are gotten successfully.")

	var result scheme
	err = json.Unmarshal(schemeData, &result)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal random scheme data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	result.show(isHex)
}

// input is to convert input data from it to json in all commands to do request
// to colormind api.
type input struct {
	Model  string   `json:"model"`
	Scheme [][3]int `json:"input,omitempty"`
}

// scheme is to convert json color scheme to it in random and suggest commands after
// getting response from colormind api.
type scheme struct {
	Result [5][3]int `json:"result"`
}

// show displays color scheme as list in hexadecimal or RGB(decimal) format.
func (s *scheme) show(isHex bool) {
	for _, rgb := range s.Result {
		var color string
		if isHex {
			color = fmt.Sprintf("#%X%X%X", rgb[0], rgb[1], rgb[2])
		} else {
			color = fmt.Sprintf("%d %d %d", rgb[0], rgb[1], rgb[2])
		}
		fmt.Println(color)
	}
}
