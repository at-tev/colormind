package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"colormind/errors"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(modelsCmd)
}

// models is to convert models list json to it.
type modelsList struct {
	List       []string `json:"result"`
	lastUpdate time.Time
}

// show displays model values as list.
func (m *modelsList) show() {
	fmt.Println("Here is a list of currently available models:")
	for _, model := range m.List {
		fmt.Println("‣ ", model)
	}
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Get currently available models",
	Long: "This command gets list of currently available models for color schemes from" +
		" the colormind api. The list of models gets updated every day at 07:00 0+UTC",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		models.show()
	},
}

// models is list of currently available model values.
var models modelsList

// Models gets models list by getData function and converts it from json to models type.
func Models() {
	param := input{}
	modelsData, err := getData(param)
	if err != nil {
		err = fmt.Errorf("couldn't get models data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	sentry.CaptureMessage("Models list data are gotten successfully.")

	err = json.Unmarshal(modelsData, &models)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal models data: %w", err)
		errors.HandleError(errors.InternalErr, err)
		return
	}

	viper.Set("models.list", models.List)
	viper.Set("models.last_update", time.Now().UTC())
	viper.WriteConfig()
}

// getData gets scheme and models data from colormind by doing request to its api.
func getData(param input) ([]byte, error) {
	var body io.Reader
	method := "GET"
	url := modelsURL

	if param.Model != "" {
		method = "POST"
		url = schemeURL
		data, err := json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("couldn't  marshal input data: %w", err)
		}
		body = bytes.NewReader(data)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("bad %s request: %w", method, err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %w", err)
	}

	defer response.Body.Close()

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}

	return result, nil
}

// updated checks if list of model values is updated
func (m modelsList) updated() bool {
	cur := time.Now().UTC()
	timeToUpdate := time.Date(cur.Year(), cur.Month(), cur.Day(), 7, 0, 29, 999999, time.UTC)
	if cur.Hour() < 7 {
		timeToUpdate.AddDate(0, 0, -1)
	}

	return m.lastUpdate.After(timeToUpdate)
}
