package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"colormind/errors"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(suggestCmd)
	suggestCmd.PersistentFlags().VarP(&model, "model", "m", "a model for color scheme.")
	suggestCmd.PersistentFlags().BoolP("hexadecimal", "H", false, "enter and show colors in hexadecimal characters")
}

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Get suggestion for color scheme",
	Long:  "This command suggests colors for your color scheme from colormind api.",
	Args:  cobra.RangeArgs(1, 4),
	Run: func(cmd *cobra.Command, args []string) {
		isHex, err := cmd.Flags().GetBool("hexadecimal")
		if err != nil {
			fmt.Println(err)
			return
		}
		Suggestion(args, model.String(), isHex)
	},
}

// Suggestion gets suggestion for color scheme dependong on model flag's value
// and partial scheme input from terminal arguments by getData and converts it from json
// to scheme type.
func Suggestion(halfScheme []string, model string, isHex bool) {
	var colors [][3]int
	var err error

	if isHex {
		colors, err = toHexScheme(halfScheme)
	} else {
		colors, err = toScheme(halfScheme)
	}

	if err != nil {
		err = fmt.Errorf("couldn't identify color: %w", err)
		errors.HandleError("An error occurred:", err)
		return
	}

	param := input{Model: model, Scheme: colors}
	schemeData, err := getData(param)
	if err != nil {
		err = fmt.Errorf("couldn't get scheme suggestion data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	sentry.CaptureMessage("Suggested scheme data are gotten successfully.")

	var result scheme
	err = json.Unmarshal(schemeData, &result)
	if err != nil {
		err = fmt.Errorf("couldn't  unmarshal scheme suggestion data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	result.show(isHex)
}

// formatScheme converts string terminal arguments to slice of colors in RGB(decimal) format.
func toScheme(scheme []string) ([][3]int, error) {
	var colors [][3]int

	for i := range scheme {
		rgb := strings.FieldsFunc(scheme[i], func(r rune) bool {
			return !unicode.IsDigit(r)
		})
		if len(rgb) == 0 {
			return nil, fmt.Errorf("there is no rgb color code in '%s'", scheme[i])
		}

		// Formatting string color to RGB color code with 3 color channels.
		// There are zeroes instead of empty color channels.
		var decimalCode [3]int
		for n, c := range rgb {
			code, err := strconv.Atoi(c)
			if err != nil {
				return nil, err
			}
			decimalCode[n] = code

			// Cutting off a rgb color code that has more than 3 color channels.
			if n == 2 {
				break
			}
		}

		if !isValidRGB(decimalCode) {
			return nil, fmt.Errorf("color '%s' is not valid rgb color code", scheme[i])
		}

		colors = append(colors, decimalCode)
	}

	return colors, nil
}

// formatScheme converts string terminal arguments to slice of colors in hexadecimal format.

func toHexScheme(scheme []string) ([][3]int, error) {
	var colors [][3]int

	for i := range scheme {
		hex := strings.TrimFunc(scheme[i], func(r rune) bool {
			return !(unicode.IsDigit(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f'))
		})
		if len(hex) == 0 {
			return nil, fmt.Errorf("there is no hexadecimal color code in '%s'", scheme[i])
		}

		if !isValidHex(hex) {
			return nil, fmt.Errorf("color '%s' is not valid hexadecimal color code", scheme[i])
		}

		// Formatting string color to hex color code.
		// There are zeroes instead of missig hex values.
		var decimalCode [3]int
		var m int
		for n := 0; n < len(hex); n += 2 {
			fmt.Sscanf(hex[n:n+2], "%X", &decimalCode[m])
			// Cutting off a hex color code that has more than 6 digits.
			if m == 2 {
				break
			}
			m++
		}

		colors = append(colors, decimalCode)
	}

	return colors, nil
}

// isValidRGB checks if color is valid RGB color code by checking color channels values.
func isValidRGB(rgb [3]int) bool {
	rgbChanMin := 0
	rgbChanMax := 255 // Every RGB color channel has to be value from 0 to 255.
	for i := range rgb {
		if rgb[i] > rgbChanMax || rgb[i] < rgbChanMin {
			return false
		}
	}
	return true
}

// isValidHex checks if color is valid hexadecimal value of even hexadecimal characters.
func isValidHex(color string) bool {
	matched, _ := regexp.Match(`^[A-Fa-f0-9]+$`, []byte(color))
	return len(color)%2 == 0 && matched
}
