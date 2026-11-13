
# Project 3 - Country Explorer

Name: Shabaan Saleem
Course: COMP 415 - Emerging Languages

## About the Project

For this project, I made a GUI program in Go that shows information about different countries in Europe. The program gets the country information from the REST Countries API.

The user can choose a country from a dropdown menu, and the program will display its flag and information.

## How to Run

1. Open the project in GoLand.
2. Make sure Go and a C compiler are installed.
3. Run `go mod tidy` if the dependencies are missing.
4. Run the program using `go run .` or the Run button in GoLand.
5. Enter your REST Countries API key.
6. Click "Load European Countries".
7. Choose a country from the dropdown menu.

## What the Program Shows

When a country is selected, the program displays:

- Common name
- Official name
- Official languages
- Country flag
- Whether the country is landlocked
- Population
- Capital
- Subregion

## Decisions I Made

I chose Europe as my region for this project.

I used the Fyne library to make the GUI because it has widgets like dropdown menus, buttons, labels, and images.

I used structs to organize the country information that comes from the API.

I also sorted the countries alphabetically so they are easier to find in the dropdown menu.

Instead of putting my API key directly in the code, I added a text box where the user can enter their own key.

## Testing

I tested the program with different countries, including Croatia and Bulgaria. The dropdown worked, and the program displayed their flags and country information correctly.

## Anything Not Completed

Everything I planned for this project is completed.
