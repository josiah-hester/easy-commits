package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url,omitempty"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "config":
		handleConfig()
	case "commit":
		handleCommit()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Easy Commits - AI-powered git commit message generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  easy-commits config    Configure AI provider and API key")
	fmt.Println("  easy-commits commit    Generate and create a commit with AI-generated message")
	fmt.Println("  easy-commits help      Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  easy-commits config")
	fmt.Println("  easy-commits commit")
	fmt.Println("  easy-commits commit --context \"Fixed the login bug\"")
}

func handleConfig() {
	config := Config{}

	fmt.Print("Select AI provider (openai/anthropic/ollama): ")
	fmt.Scanln(&config.Provider)

	if config.Provider == "ollama" {
		fmt.Print("Enter Ollama base URL (default: http://localhost:11434): ")
		var baseURL string
		fmt.Scanln(&baseURL)
		if baseURL == "" {
			config.BaseURL = "http://localhost:11434"
		} else {
			config.BaseURL = baseURL
		}
		fmt.Print("Enter model name (e.g., llama2, codellama): ")
		fmt.Scanln(&config.Model)
	} else {
		fmt.Print("Enter API key: ")
		fmt.Scanln(&config.APIKey)

		switch config.Provider {
		case "openai":
			config.Model = "gpt-3.5-turbo"
		case "anthropic":
			config.Model = "claude-3-haiku-20240307"
		}
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	configPath := fmt.Sprintf("%s/.easy-commits-config.json", homeDir)
	err = os.WriteFile(configPath, configData, 0600)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
}

func handleCommit() {
	// Check if we're in a git repository
	if !isGitRepo() {
		fmt.Println("Error: Not in a git repository")
		return
	}

	// Get git diff
	diff, err := getGitDiff()
	if err != nil {
		fmt.Printf("Error getting git diff: %v\n", err)
		return
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No changes to commit")
		return
	}

	// Get user context if provided
	var userContext string
	for i, arg := range os.Args {
		if arg == "--context" && i+1 < len(os.Args) {
			userContext = os.Args[i+1]
			break
		}
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Run 'easy-commits config' to set up your AI provider")
		return
	}

	// Generate commit message
	commitMessage, err := generateCommitMessage(config, diff, userContext)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		return
	}

	// Show the generated message and ask for confirmation
	fmt.Println("Generated commit message:")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println(commitMessage)
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Print("Use this commit message? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		err = createCommit(commitMessage)
		if err != nil {
			fmt.Printf("Error creating commit: %v\n", err)
			return
		}
		fmt.Println("Commit created successfully!")
	} else {
		fmt.Println("Commit cancelled")
	}
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// If no staged changes, get unstaged changes
	if strings.TrimSpace(string(output)) == "" {
		cmd = exec.Command("git", "diff")
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
	}

	return string(output), nil
}

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := fmt.Sprintf("%s/.easy-commits-config.json", homeDir)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

func generateCommitMessage(config *Config, diff, userContext string) (string, error) {
	prompt := buildPrompt(diff, userContext)

	switch config.Provider {
	case "openai":
		return callOpenAI(config, prompt)
	case "anthropic":
		return callAnthropic(config, prompt)
	case "ollama":
		return callOllama(config, prompt)
	default:
		return "", fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

func buildPrompt(diff, userContext string) string {
	prompt := `You are an expert at writing clear, concise git commit messages. Based on the git diff provided, generate a commit message that follows these guidelines:

1. Use the conventional commit format: type(scope): description
2. Types: feat, fix, docs, style, refactor, test, chore
3. Keep the first line under 50 characters
4. Use imperative mood ("add" not "added")
5. Be specific about what changed and why

Git diff:
` + diff

	if userContext != "" {
		prompt += "\n\nAdditional context from user: " + userContext
	}

	prompt += "\n\nGenerate only the commit message, no additional text or explanation."

	return prompt
}

func callOpenAI(config *Config, prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model: config.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp OpenAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
}

func callAnthropic(config *Config, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":      config.Model,
		"max_tokens": 150,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var anthropicResp map[string]interface{}
	err = json.Unmarshal(body, &anthropicResp)
	if err != nil {
		return "", err
	}

	content, ok := anthropicResp["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	textContent, ok := content[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format from Anthropic")
	}

	text, ok := textContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text in Anthropic response")
	}

	return strings.TrimSpace(text), nil
}

func callOllama(config *Config, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  config.Model,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := config.BaseURL + "/api/generate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp map[string]interface{}
	err = json.Unmarshal(body, &ollamaResp)
	if err != nil {
		return "", err
	}

	response, ok := ollamaResp["response"].(string)
	if !ok {
		return "", fmt.Errorf("no response from Ollama")
	}

	return strings.TrimSpace(response), nil
}

func createCommit(message string) error {
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stage changes: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}
