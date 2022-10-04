package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Ladder struct {
	Data []struct {
		AdditionalData struct {
			AbilityCost1        string `json:"abilityCost_1"`
			AbilityCost2        string `json:"abilityCost_2"`
			AbilityCost3        string `json:"abilityCost_3"`
			AbilityCost4        string `json:"abilityCost_4"`
			AbilityCost5        string `json:"abilityCost_5"`
			AbilityDescription1 string `json:"abilityDescription_1"`
			AbilityDescription2 string `json:"abilityDescription_2"`
			AbilityDescription3 string `json:"abilityDescription_3"`
			AbilityDescription4 string `json:"abilityDescription_4"`
			AbilityDescription5 string `json:"abilityDescription_5"`
			AbilityID1          string `json:"abilityID_1"`
			AbilityID2          string `json:"abilityID_2"`
			AbilityID3          string `json:"abilityID_3"`
			AbilityID4          string `json:"abilityID_4"`
			AbilityID5          string `json:"abilityID_5"`
			AbilityImagePath1   string `json:"abilityImagePath_1"`
			AbilityImagePath2   string `json:"abilityImagePath_2"`
			AbilityImagePath3   string `json:"abilityImagePath_3"`
			AbilityImagePath4   string `json:"abilityImagePath_4"`
			AbilityImagePath5   string `json:"abilityImagePath_5"`
			AbilityName1        string `json:"abilityName_1"`
			AbilityName2        string `json:"abilityName_2"`
			AbilityName3        string `json:"abilityName_3"`
			AbilityName4        string `json:"abilityName_4"`
			AbilityName5        string `json:"abilityName_5"`
			AbilityTags1        string `json:"abilityTags_1"`
			AbilityTags2        string `json:"abilityTags_2"`
			AbilityTags3        string `json:"abilityTags_3"`
			AbilityTags4        string `json:"abilityTags_4"`
			AbilityTags5        string `json:"abilityTags_5"`
			CharClass           string `json:"char_class"`
			CharLvl             string `json:"char_lvl"`
			CharName            string `json:"char_name"`
			Deaths              int64  `json:"deaths"`
			Died                bool   `json:"died"`
			Exp                 string `json:"exp"`
			MaxWave             string `json:"max_wave"`
			PlayerUsername      string `json:"player_username"`
		} `json:"additionalData"`
		Point int64 `json:"point"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func CreateLeaderboardReport(apiURL string, fileName string, leaderboardName string) error {
	// Setup vars
	var jsonData Ladder
	var totalUsed = make(map[string]int)
	var skillOne = make(map[string]int)
	var skillTwo = make(map[string]int)
	var skillThree = make(map[string]int)
	var skillFour = make(map[string]int)
	var skillFive = make(map[string]int)
	var FlagOne bool
	var FlagTwo bool
	var FlagThree bool
	var FlagFour bool
	var FlagFive bool
	FlagOne = false
	FlagTwo = false
	FlagThree = false
	FlagFour = false
	FlagFive = false
	var min = 1
	var minSkillName string
	var max int
	var maxSkillName string
	// Setup  client
	httpClient := http.DefaultClient
	// Get the data
	resp, err := httpClient.Get(apiURL)
	if err != nil {
		panic(err)
	}
	// Close the data once done
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	// Read it
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//  unmarshall the data
	err = json.Unmarshal(respData, &jsonData)
	if err != nil {
		panic(err)
	}
	// LE sends "Data retrived" as the message if the API call is successful.
	if jsonData.Message != "Data retrived" {
		fmt.Println("Error: ", jsonData.Message)
		panic("API request failed")
	}
	// Iterate through top 100 of all time.
	for i := range jsonData.Data {
		// j = ability names
		for j := range skillOne {
			// if the ability name 1 is in the map, add 1 to the value
			if jsonData.Data[i].AdditionalData.AbilityName1 == j {
				skillOne[j]++
				// set flag to the true so we don't add it and reset to 1 later
				FlagOne = true
			}
		}
		// same loop, a bit repetitive code. I know.
		if FlagOne == false {
			skillOne[jsonData.Data[i].AdditionalData.AbilityName1] = 1
		}
		for j := range skillTwo {
			if jsonData.Data[i].AdditionalData.AbilityName2 == j {
				skillTwo[j]++
				FlagTwo = true
			}
		}
		if FlagTwo == false {
			skillTwo[jsonData.Data[i].AdditionalData.AbilityName2] = 1
		}
		for j := range skillThree {
			if jsonData.Data[i].AdditionalData.AbilityName3 == j {
				skillThree[j]++
				FlagThree = true
			}
		}
		if FlagThree == false {
			skillThree[jsonData.Data[i].AdditionalData.AbilityName3] = 1
		}
		for j := range skillFour {
			if jsonData.Data[i].AdditionalData.AbilityName4 == j {
				skillFour[j]++
				FlagFour = true
			}
		}
		if FlagFour == false {
			skillFour[jsonData.Data[i].AdditionalData.AbilityName4] = 1
		}
		for j := range skillFive {
			if jsonData.Data[i].AdditionalData.AbilityName5 == j {
				skillFive[j]++
				FlagFive = true
			}
		}
		if FlagFive == false {
			skillFive[jsonData.Data[i].AdditionalData.AbilityName5] = 1
		}
		// reset the flags
		FlagOne = false
		FlagTwo = false
		FlagThree = false
		FlagFour = false
		FlagFive = false
	}
	// Main loop:
	// 1. Iterate through the map
	// 2. If the value is greater than the min, set the min to the value
	// 3. If the value is less than the max, set the max to the value
	// In either cases, if skill name exists on totalused map, add 1 to the value
	// If it doesn't exist, set the value to 1.
	// This is done for each skill.
	file, err := os.Create(fileName)
	for i := range skillOne {
		if skillOne[i] > max {
			max = skillOne[i]
			maxSkillName = i
			totalUsed[i]++
		}
		if skillOne[i] < min {
			min = skillOne[i]
			minSkillName = i
			totalUsed[i]++
		}
		if skillOne[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "&" + i
			}
			totalUsed[i]++
		}
		if skillOne[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "&" + i
			}
			totalUsed[i]++
		}
	}
	_, err = fmt.Fprintln(file, "=======", leaderboardName, "=======")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "LEADERBOARD FIRST PLACE")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "Char name > ", jsonData.Data[0].AdditionalData.CharName)
	_, err = fmt.Fprintln(file, "Char level > ", jsonData.Data[0].AdditionalData.CharLvl)
	_, err = fmt.Fprintln(file, "Char class > ", jsonData.Data[0].AdditionalData.CharClass)
	_, err = fmt.Fprintln(file, "Char's maximum arena wave > ", jsonData.Data[0].AdditionalData.MaxWave)
	_, err = fmt.Fprintln(file, "Char's total amount of deaths > ", jsonData.Data[0].AdditionalData.Deaths)
	_, err = fmt.Fprintln(file, "Char's first ability > ", jsonData.Data[0].AdditionalData.AbilityName1)
	_, err = fmt.Fprintln(file, "Char's second ability > ", jsonData.Data[0].AdditionalData.AbilityName2)
	_, err = fmt.Fprintln(file, "Char's third ability > ", jsonData.Data[0].AdditionalData.AbilityName3)
	_, err = fmt.Fprintln(file, "Char's fourth ability > ", jsonData.Data[0].AdditionalData.AbilityName4)
	_, err = fmt.Fprintln(file, "Char's fifth ability > ", jsonData.Data[0].AdditionalData.AbilityName5)
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "SKILL SLOTS")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "Skill 1")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Maximum skill put into first slot:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	// Reset min, max, and minSkillName, maxSkillName.
	min = 1
	max = 1
	minSkillName = ""
	maxSkillName = ""
	// Make the same loop for ability 2.
	for i := range skillTwo {
		if skillTwo[i] > max {
			max = skillTwo[i]
			maxSkillName = i
			totalUsed[i]++
		}
		if skillTwo[i] < min {
			min = skillTwo[i]
			minSkillName = i
			totalUsed[i]++
		}
		if skillTwo[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "&" + i
			}
			totalUsed[i]++
		}
		if skillTwo[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "&" + i
			}
			totalUsed[i]++
		}
	}
	_, err = fmt.Fprintln(file, "Skill 2")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Maximum skill put into second slot:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	// Reset min, max, and minSkillName, maxSkillName.
	min = 1
	max = 1
	minSkillName = ""
	maxSkillName = ""
	// Make the same loop for ability 3.
	for i := range skillThree {
		if skillThree[i] > max {
			max = skillThree[i]
			maxSkillName = i
			totalUsed[i]++
		}
		if skillThree[i] < min {
			min = skillThree[i]
			minSkillName = i
			totalUsed[i]++
		}
		if skillThree[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "&" + i
			}
			totalUsed[i]++
		}
		if skillThree[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "&" + i
			}
			totalUsed[i]++
		}
	}
	_, err = fmt.Fprintln(file, "Skill 3")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Maximum skill put into third slot:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	// Reset min, max, and minSkillName, maxSkillName.
	min = 1
	max = 1
	minSkillName = ""
	maxSkillName = ""
	// Make the same loop for ability 4.
	for i := range skillFour {
		if skillFour[i] > max {
			max = skillFour[i]
			maxSkillName = i
			totalUsed[i]++
		}
		if skillFour[i] < min {
			min = skillFour[i]
			minSkillName = i
			totalUsed[i]++
		}
		if skillFour[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "&" + i
			}
			totalUsed[i]++
		}
		if skillFour[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "&" + i
			}
			totalUsed[i]++
		}
	}
	_, err = fmt.Fprintln(file, "Skill 4")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Maximum skill put into fourth slot:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	// Reset min, max, and minSkillName, maxSkillName.
	min = 1
	max = 1
	minSkillName = ""
	maxSkillName = ""
	// Make the same loop for ability 5.
	for i := range skillFive {
		if skillFive[i] > max {
			max = skillFour[i]
			maxSkillName = i
			totalUsed[i]++
		}
		if skillFive[i] < min {
			min = skillFour[i]
			minSkillName = i
			// If the skill is already in the map, add 1 to the value.
			totalUsed[i]++
		}
		if skillFive[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "&" + i
			}
			totalUsed[i]++
		}
		if skillFive[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "&" + i
			}
			// If the skill is already in the map, add 1 to the value.
			totalUsed[i]++
		}
	}
	_, err = fmt.Fprintln(file, "Skill 5")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Maximum skill put into fifth slot:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "HOW MANY UNIQUE PLAYERS USED EACH SKILL")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "Total number of times each skill was used:")
	if err != nil {
		panic(err)
	}
	min = 1
	max = 1
	for i := range totalUsed {
		_, err = fmt.Fprintln(file, i, ":", totalUsed[i])
		if err != nil {
			panic(err)
		}
		if totalUsed[i] > max {
			max = totalUsed[i]
			maxSkillName = i
		}
		if totalUsed[i] <= min {
			min = totalUsed[i]
			minSkillName = i
		}
		if totalUsed[i] == max {
			if strings.Contains(maxSkillName, i) {
				continue
			} else {
				maxSkillName = maxSkillName + "/" + i
			}
		}
		if totalUsed[i] == min {
			if strings.Contains(minSkillName, i) {
				continue
			} else {
				minSkillName = minSkillName + "/" + i
			}
		}
	}
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "MAXIMUM / MINIMUM UNIQUE SKILLS USED")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "===============================")
	_, err = fmt.Fprintln(file, "Maximum number of unique skills used:", maxSkillName, "with players of:", max)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, "Minimum number of unique skills used:", minSkillName, "with players of:", min)
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	err := CreateLeaderboardReport("https://leapi.lastepoch.com/api/leader-board?code=beta085softcoreallclassarenawave", "leaderboardtop100.txt", "Leaderboard Class Independent Top 100")
	if err != nil {
		return
	}
}
