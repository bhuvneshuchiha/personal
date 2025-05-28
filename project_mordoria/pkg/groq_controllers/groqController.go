package groq_controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	finalMessage "github.com/bhuvneshuhciha/project_mordoria/pkg/final_message"
)

type PreparedDataStruct struct {
	MessageList []string
	AiScore     string
	MainMessage string
}

var DataStruct = &PreparedDataStruct{
	MessageList: make([]string, 3),
	AiScore:     "",
	MainMessage: "",
}

type GroqMessage struct {
	Role    string
	Content string
}

type GroqRequest struct {
	Model   string
	Message []GroqMessage
}

type GroqResponse struct {
	Choices []struct {
		Message GroqMessage `json:"message"`
	} `json:"choices"`
}

func PrepareData(msgPayload *finalMessage.FinalPayload, dataStruct *PreparedDataStruct) error {
	tempList := make([]string, 5)
	averageScore := msgPayload.Ai_emot_score

	DataStruct.MessageList = tempList
	DataStruct.AiScore = averageScore
	DataStruct.MainMessage = `You are a smart assistant, your job is to take all these messages\n
							present in the array, and create a summary out of those messages.
							The emotion of the summary should be determind by the ai_score parameter.
							If the ai_score is between 0-2, the the response should be really sad,
							depressing, and dull. If the ai_score is betweeen 3-5, the the message
							should be rude and arrogant and you should act like a toxic dad. If the
							ai_score is between 5-8, then you should response with very sexual content
							and try to make the users uncomfortable.If the response is above 8 then
							you should be very very witty, sarcastic, smart, egoistic and charismatic.`

	msg := msgPayload.Payload
	for _, val := range msg {
		tempList = append(tempList, val.MessageString)
	}
	if len(tempList) > 0 && (averageScore != "0") {
	} else {
		return errors.New("Either temp list was empty or average score is not string")
	}
	return nil
}

func SendDataToGroq() (string, error) {
	PrepareData(finalMessage.MsgBody, DataStruct)
	promptBytes, err := json.MarshalIndent(DataStruct, "", "")
	if err != nil {
		log.Println("Json was not properly converted to string", err)
		return "Error:", err
	}
	prettyPrompt := string(promptBytes)

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "api key was not found", errors.New("Api key not found")
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	requestBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Message: []GroqMessage{
			{
				Role:    "user",
				Content: prettyPrompt,
			},
		},
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error while sending to groq's endpoint", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", errors.New("Didnt get any response from groq api")
	}

	var groqResponse GroqResponse
	if err := json.Unmarshal(body, &groqResponse); err != nil {
		return "", err
	}
	if len(groqResponse.Choices) > 0 {
		return groqResponse.Choices[0].Message.Content, nil
	}

	return "", errors.New("No response from groq")
}
