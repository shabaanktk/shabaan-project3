package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Professor's REST Countries API, version 5.
const apiURL = "https://api.restcountries.com/countries/v5"

// Country represents the fields needed for this assignment.
type Country struct {
	Names struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"names"`

	Languages  []json.RawMessage `json:"languages"`
	Landlocked bool              `json:"landlocked"`
	Population int64             `json:"population"`
	Subregion  string            `json:"subregion"`

	Capitals []struct {
		Name string `json:"name"`
	} `json:"capitals"`

	Flag struct {
		URLPNG string `json:"url_png"`
	} `json:"flag"`
}

// The v5 API places country records inside data.objects.
type APIResponse struct {
	Data struct {
		Objects []Country `json:"objects"`
		Meta    struct {
			More bool `json:"more"`
		} `json:"meta"`
	} `json:"data"`

	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// Fetch all pages of European countries from the API.
func getCountries(apiKey string) ([]Country, error) {
	client := &http.Client{Timeout: 25 * time.Second}
	var countries []Country

	for offset := 0; ; offset += 100 {
		params := url.Values{}
		params.Set("region", "Europe")
		params.Set("limit", "100")
		params.Set("offset", fmt.Sprint(offset))

		endpoint := apiURL + "?" + params.Encode()

		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		var result APIResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("invalid API response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			message := resp.Status
			if len(result.Errors) > 0 {
				message = result.Errors[0].Message
			}
			return nil, fmt.Errorf("API error: %s", message)
		}

		countries = append(countries, result.Data.Objects...)

		if !result.Data.Meta.More {
			break
		}
	}

	if len(countries) == 0 {
		return nil, fmt.Errorf("no European countries were returned")
	}

	// Alphabetize the country dropdown.
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Names.Common < countries[j].Names.Common
	})

	return countries, nil
}

// Convert the language records into readable names.
func getLanguages(languages []json.RawMessage) string {
	var names []string

	for _, raw := range languages {
		var language map[string]interface{}

		if json.Unmarshal(raw, &language) != nil {
			continue
		}

		// Support the language name fields returned by the API.
		name := ""
		for _, field := range []string{"name", "english", "name_en"} {
			if value, ok := language[field].(string); ok && value != "" {
				name = value
				break
			}
		}

		// Some language records have nested names.
		if name == "" {
			if nested, ok := language["names"].(map[string]interface{}); ok {
				if value, ok := nested["english"].(string); ok {
					name = value
				}
			}
		}

		if name != "" {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return "Not available"
	}

	sort.Strings(names)
	return strings.Join(names, ", ")
}

// Join the country's capital cities.
func getCapitals(country Country) string {
	var capitals []string

	for _, capital := range country.Capitals {
		if capital.Name != "" {
			capitals = append(capitals, capital.Name)
		}
	}

	if len(capitals) == 0 {
		return "Not available"
	}

	return strings.Join(capitals, ", ")
}

// Download a PNG flag and decode it into an image.
func downloadFlag(flagURL string) (image.Image, error) {
	if flagURL == "" {
		return nil, fmt.Errorf("flag not available")
	}

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(flagURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flag download failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Shabaan - Project 3 Country Explorer")
	window.Resize(fyne.NewSize(950, 760))

	title := widget.NewLabelWithStyle(
		"European Country Explorer",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	instructions := widget.NewLabel(
		"Enter your REST Countries API key, load countries, then choose one.",
	)
	instructions.Wrapping = fyne.TextWrapWord

	// The key is entered at runtime, not saved in source code.
	keyInput := widget.NewPasswordEntry()
	keyInput.SetPlaceHolder("Paste your REST Countries API key here")

	status := widget.NewLabel("Ready to load countries.")
	status.Wrapping = fyne.TextWrapWord

	commonName := widget.NewLabel("Common Name: ")
	officialName := widget.NewLabel("Official Name: ")
	languages := widget.NewLabel("Official Languages: ")
	landlocked := widget.NewLabel("Landlocked: ")
	population := widget.NewLabel("Population: ")
	capital := widget.NewLabel("Capital: ")
	subregion := widget.NewLabel("Subregion: ")

	officialName.Wrapping = fyne.TextWrapWord
	languages.Wrapping = fyne.TextWrapWord

	// Reserve a visible space for the country's flag.
	flag := canvas.NewImageFromImage(nil)
	flag.FillMode = canvas.ImageFillContain
	flag.SetMinSize(fyne.NewSize(360, 190))

	flagTitle := widget.NewLabelWithStyle(
		"Country Flag",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	countryMap := make(map[string]Country)
	selectedVersion := 0

	// Display information when a country is selected.
	selectCountry := widget.NewSelect([]string{}, func(name string) {
		country, exists := countryMap[name]
		if !exists {
			return
		}

		selectedVersion++
		version := selectedVersion

		commonName.SetText("Common Name: " + country.Names.Common)
		officialName.SetText("Official Name: " + country.Names.Official)
		languages.SetText("Official Languages: " + getLanguages(country.Languages))

		if country.Landlocked {
			landlocked.SetText("Landlocked: Yes")
		} else {
			landlocked.SetText("Landlocked: No")
		}

		population.SetText(fmt.Sprintf("Population: %d", country.Population))
		capital.SetText("Capital: " + getCapitals(country))
		subregion.SetText("Subregion: " + country.Subregion)

		status.SetText("Selected: " + name)

		flag.Image = nil
		flag.Refresh()

		// Download the flag without freezing the window.
		go func() {
			img, err := downloadFlag(country.Flag.URLPNG)

			// Fyne GUI updates must run on its main thread.
			fyne.Do(func() {
				if version != selectedVersion {
					return
				}

				if err != nil {
					status.SetText("Country loaded. Flag unavailable: " + err.Error())
					return
				}

				flag.Image = img
				flag.Refresh()
			})
		}()
	})

	selectCountry.PlaceHolder = "Choose a European country"

	// Button for retrieving the country data.
	loadButton := widget.NewButton("Load European Countries", nil)

	loadButton.OnTapped = func() {
		key := strings.TrimSpace(keyInput.Text)

		if key == "" {
			status.SetText("Please enter your API key first.")
			return
		}

		loadButton.Disable()
		status.SetText("Loading European countries...")

		go func() {
			countries, err := getCountries(key)

			fyne.Do(func() {
				loadButton.Enable()

				if err != nil {
					status.SetText("Error loading countries: " + err.Error())
					return
				}

				countryMap = make(map[string]Country)
				var names []string

				for _, country := range countries {
					name := country.Names.Common
					if name == "" {
						continue
					}

					names = append(names, name)
					countryMap[name] = country
				}

				selectCountry.Options = names
				selectCountry.ClearSelected()
				selectCountry.Refresh()

				status.SetText(fmt.Sprintf(
					"Successfully loaded %d European countries. Choose one!",
					len(names),
				))
			})
		}()
	}

	details := container.NewVBox(
		commonName,
		officialName,
		languages,
		landlocked,
		population,
		capital,
		subregion,
	)

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		instructions,
		keyInput,
		loadButton,
		selectCountry,
		status,
		widget.NewSeparator(),
		flagTitle,
		flag,
		widget.NewSeparator(),
		details,
	)

	window.SetContent(
		container.NewPadded(
			container.NewVScroll(content),
		),
	)

	window.ShowAndRun()
}
