package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Configuration structure
type Config struct {
	APIKeys map[string]string `json:"api_keys"`
	DefaultModel string        `json:"default_model"`
}

// Model represents the current state of the application
type Model struct {
	config       Config
	currentModel string
	textarea     textarea.Model
	textinput    textinput.Model
	messages     []Message
	viewport     []string
	width        int
	height       int
	inputMode    bool
	configMode   bool
	err          error
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// API response structures
type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

type GrokResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	modelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	userStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	assistantStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)
)

func main() {
	// Load or create config
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the model
	model := initialModel(config)

	// Start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(config Config) Model {
	// Initialize textarea for multi-line input
	ta := textarea.New()
	ta.Placeholder = "Type your message here..."
	ta.Focus()
	ta.CharLimit = 4000
	ta.SetWidth(80)
	ta.SetHeight(3)

	// Initialize textinput for API keys
	ti := textinput.New()
	ti.Placeholder = "Enter API key..."
	ti.CharLimit = 200
	ti.Width = 60

	currentModel := config.DefaultModel
	if currentModel == "" {
		currentModel = "grok"
	}

	return Model{
		config:       config,
		currentModel: currentModel,
		textarea:     ta,
		textinput:    ti,
		messages:     []Message{},
		viewport:     []string{},
		inputMode:    true,
		configMode:   false,
		width:        80,
		height:       24,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.configMode {
				m.configMode = false
				return m, nil
			}
			return m, tea.Quit

		case "ctrl+m": // Switch model
			m.currentModel = switchModel(m.currentModel)
			return m, nil

		case "ctrl+k": // Configure API keys
			m.configMode = !m.configMode
			if m.configMode {
				m.textinput.Focus()
			} else {
				m.textarea.Focus()
			}
			return m, nil

		case "ctrl+r": // Reset conversation
			m.messages = []Message{}
			m.viewport = []string{}
			return m, nil

		case "enter":
			if m.configMode {
				// Handle API key input
				key := strings.TrimSpace(m.textinput.Value())
				if key != "" {
					m.config.APIKeys[m.currentModel] = key
					saveConfig(m.config)
					m.textinput.SetValue("")
					m.configMode = false
					m.textarea.Focus()
				}
				return m, nil
			}

			// Handle message input
			userInput := strings.TrimSpace(m.textarea.Value())
			if userInput == "" {
				return m, nil
			}

			// Add user message
			m.messages = append(m.messages, Message{Role: "user", Content: userInput})
			m.viewport = append(m.viewport, userStyle.Render("You: ")+userInput)

			// Get tmux context
			context := getTmuxContext()
			if context != "" {
				contextMsg := fmt.Sprintf("Current tmux session context:\n%s\n\nUser question: %s", context, userInput)
				m.messages[len(m.messages)-1].Content = contextMsg
			}

			// Clear input
			m.textarea.SetValue("")

			// Send to AI
			return m, m.sendToAI()

		default:
			if m.configMode {
				var cmd tea.Cmd
				m.textinput, cmd = m.textinput.Update(msg)
				return m, cmd
			} else {
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}
		}

	case aiResponseMsg:
		response := string(msg)
		m.messages = append(m.messages, Message{Role: "assistant", Content: response})
		m.viewport = append(m.viewport, assistantStyle.Render("AI: ")+response)
		return m, nil

	case errMsg:
		m.err = error(msg)
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("AI CLI Assistant") + " " +
		modelStyle.Render(fmt.Sprintf("Model: %s", m.currentModel))
	b.WriteString(header + "\n\n")

	// Show configuration mode
	if m.configMode {
		b.WriteString("Configure API Key for " + m.currentModel + ":\n")
		b.WriteString(m.textinput.View() + "\n")
		b.WriteString("Press Enter to save, Ctrl+K to cancel\n\n")
	}

	// Chat history
	start := 0
	if len(m.viewport) > m.height-10 {
		start = len(m.viewport) - (m.height - 10)
	}

	for i := start; i < len(m.viewport); i++ {
		b.WriteString(m.viewport[i] + "\n")
	}

	// Input area
	if !m.configMode {
		b.WriteString("\n" + m.textarea.View())
	}

	// Help text
	help := "\nControls: Ctrl+M (switch model) | Ctrl+K (API keys) | Ctrl+R (reset) | Ctrl+C (quit)"
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(help))

	// Error display
	if m.err != nil {
		b.WriteString("\n" + errorStyle.Render("Error: "+m.err.Error()))
	}

	return b.String()
}

// Custom message types
type aiResponseMsg string
type errMsg error

func (m Model) sendToAI() tea.Cmd {
	return func() tea.Msg {
		apiKey, exists := m.config.APIKeys[m.currentModel]
		if !exists || apiKey == "" {
			return errMsg(fmt.Errorf("API key not configured for %s", m.currentModel))
		}

		response, err := callAI(m.currentModel, apiKey, m.messages)
		if err != nil {
			return errMsg(err)
		}

		return aiResponseMsg(response)
	}
}

func callAI(model, apiKey string, messages []Message) (string, error) {
	switch model {
	case "grok":
		return callGrok("", messages)
	case "claude":
		return callClaude(apiKey, messages)
	case "gpt":
		return callOpenAI(apiKey, messages)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}
}

func callGrok(apiKey string, messages []Message) (string, error) {
	reqBody := map[string]any{
		"messages":    messages,
		"model":       "grok-beta",
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

	var grokResp GrokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", err
	}

	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Grok API")
	}

	return grokResp.Choices[0].Message.Content, nil
}

func callClaude(apiKey string, messages []Message) (string, error) {
	// Convert messages to Claude format
	claudeMessages := make([]map[string]string, 0)
	for _, msg := range messages {
		if msg.Role != "system" {
			claudeMessages = append(claudeMessages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}
	}

	reqBody := map[string]any{
		"model":      "claude-3-sonnet-20240229",
		"max_tokens": 1000,
		"messages":   claudeMessages,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
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

	var claudeResp AnthropicResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no response from Claude API")
	}

	return claudeResp.Content[0].Text, nil
}

func callOpenAI(apiKey string, messages []Message) (string, error) {
	reqBody := map[string]any{
		"model":       "gpt-3.5-turbo",
		"messages":    messages,
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", err
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI API")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func getTmuxContext() string {
	// Check if we're in a tmux session
	if os.Getenv("TMUX") == "" {
		return ""
	}

	var context strings.Builder

	// Get current tmux session info
	cmd := exec.Command("tmux", "display-message", "-p", "Session: #{session_name}, Window: #{window_name}, Pane: #{pane_index}")
	if output, err := cmd.Output(); err == nil {
		context.WriteString(string(output))
	}

	// Get pane content (last 50 lines)
	cmd = exec.Command("tmux", "capture-pane", "-p", "-S", "-50")
	if output, err := cmd.Output(); err == nil {
		context.WriteString("\nCurrent pane content:\n")
		context.WriteString(string(output))
	}

	// Get current working directory
	cmd = exec.Command("tmux", "display-message", "-p", "#{pane_current_path}")
	if output, err := cmd.Output(); err == nil {
		pwd := strings.TrimSpace(string(output))
		context.WriteString(fmt.Sprintf("\nCurrent directory: %s\n", pwd))

		// List files in current directory
		if files, err := os.ReadDir(pwd); err == nil {
			context.WriteString("Files in current directory:\n")
			for _, file := range files {
				if !strings.HasPrefix(file.Name(), ".") {
					context.WriteString(fmt.Sprintf("- %s\n", file.Name()))
				}
			}
		}
	}

	return context.String()
}

func switchModel(current string) string {
	models := []string{"grok", "claude", "gpt"}
	for i, model := range models {
		if model == current {
			return models[(i+1)%len(models)]
		}
	}
	return "grok"
}

func loadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(homeDir, ".ai-cli-config.json")

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := Config{
			APIKeys:      make(map[string]string),
			DefaultModel: "grok",
		}
		saveConfig(defaultConfig)
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if config.APIKeys == nil {
		config.APIKeys = make(map[string]string)
	}

	return config, err
}

func saveConfig(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".ai-cli-config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
