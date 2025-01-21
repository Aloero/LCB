package LCB

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"mime/multipart"
	"path/filepath"
)

func (b *Bot) AnswerCallbackQuery(callbackQueryID, text, show_alert string) {
	payload := map[string]string{
        "callback_query_id": callbackQueryID,
        "text":             text,
        "show_alert":       show_alert,
    }

    jsonPayload, _ := json.Marshal(payload)
    http.Post("https://api.telegram.org/bot"+ b.Token +"/answerCallbackQuery", "application/json", bytes.NewBuffer(jsonPayload))
}

func (b *Bot) SendPhoto(chatID int64, photoPathOrFileID string, caption string, utils *Utils) int {
    var requestBody io.Reader
    var err error
    var writer *multipart.Writer

    if isFileID(photoPathOrFileID) {
        message := map[string]interface{}{
            "chat_id": chatID,
            "photo":   photoPathOrFileID,
			"parse_mode": "HTML",
        }

        if caption != "" {
            message["caption"] = caption
        }

		if utils.Reply != nil {
			message["reply_markup"] = utils.Reply
		}
		if utils.Inline != nil {
			message["reply_markup"] = utils.Inline
		}

        messageJSON, err := json.Marshal(message)
        if err != nil {
            log.Println("Error marshalling message:", err)
            return 0
        }
        requestBody = bytes.NewBuffer(messageJSON)
    } else {
        file, err := os.Open(photoPathOrFileID)
        if err != nil {
            log.Println("Error opening file:", err)
            return 0
        }
        defer file.Close()

        var buffer bytes.Buffer
        writer = multipart.NewWriter(&buffer)

        photoPart, err := writer.CreateFormFile("photo", filepath.Base(photoPathOrFileID))
        if err != nil {
            log.Println("Error creating form file:", err)
            return 0
        }
        _, err = io.Copy(photoPart, file)
        if err != nil {
            log.Println("Error copying file:", err)
            return 0
        }

        err = writer.WriteField("chat_id", fmt.Sprintf("%d", chatID))
        if err != nil {
            log.Println("Error writing field:", err)
            return 0
        }
		err = writer.WriteField("parse_mode", "HTML")
        if err != nil {
            log.Println("Error writing field:", err)
            return 0
        }
        if caption != "" {
            err = writer.WriteField("caption", caption)
            if err != nil {
                log.Println("Error writing caption:", err)
                return 0
            }
        }

		if utils.Reply != nil {
			err = writer.WriteField("reply_markup", serializeKeyboard(utils.Reply))
		} else if utils.Inline != nil {
			err = writer.WriteField("reply_markup", serializeKeyboard(utils.Inline))
		}
		if err != nil {
			log.Println("Error writing keyboard:", err)
			return 0
		}


        writer.Close()
        requestBody = &buffer
    }

    url := "https://api.telegram.org/bot" + b.Token + "/sendPhoto"
    req, err := http.NewRequest("POST", url, requestBody)
    if err != nil {
        log.Println("Error creating request:", err)
        return 0
    }

    if isFileID(photoPathOrFileID) {
        req.Header.Set("Content-Type", "application/json")
    } else {
        req.Header.Set("Content-Type", writer.FormDataContentType())
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error sending request:", err)
        return 0
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
        return 0
    }

    var response ResponsePostMessage
    err = json.Unmarshal(body, &response)
    if err != nil {
        log.Fatal(err)
        return 0
    }
    if !response.Ok {
        return 0
    }

    return response.Result.MessageID
}

func isFileID(pathOrID string) bool {
    return len(pathOrID) > 0 && pathOrID[0] == 'A'
}

func serializeKeyboard(keyboard interface{}) string {
    keyboardJSON, err := json.Marshal(keyboard)
    if err != nil {
        log.Println("Error marshalling keyboard:", err)
        return ""
    }
    return string(keyboardJSON)
}

func (b *Bot) DeleteMessage(chatID int64, messageID int64) {
	message := map[string]interface{}{
		"chat_id": chatID,
		"message_id": messageID,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	url := "https://api.telegram.org/bot" + b.Token + "/deleteMessage"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from Telegram: (Status Code: %d)\n", resp.StatusCode)
		return
	}
}
 
func (b *Bot) SendDice(chatID int64, emoji string, utils Utils) int {
	message := map[string]interface{}{
		"chat_id": chatID,
		"emoji":   emoji,
	}

	if utils.ReplyMessage != nil {
		message["reply_to_message_id"] = *utils.ReplyMessage
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return 0
	}

	url := "https://api.telegram.org/bot" + b.Token + "/sendDice"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error creating request:", err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	var response ResponsePostMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	if !response.Ok {
		return 0
	}

	return response.Result.MessageID
}

func (b *Bot) EditMessage(chatID int64, messageID int64, text string, utils Utils) int {
	if len(text) > 1000 {
		text = text[:1000] + "..."
	}

	message := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
		"parse_mode": "HTML",
	}

	if utils.Reply != nil {
		message["reply_markup"] = utils.Reply
	}
	if utils.Inline != nil {
		message["reply_markup"] = utils.Inline
	}

	if utils.ReplyMessage != nil {
		message["reply_to_message_id"] = *utils.ReplyMessage
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return 0
	}

	url := "https://api.telegram.org/bot" + b.Token + "/editMessageText"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error creating request:", err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	var response ResponsePostMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	if !response.Ok {
		return 0
	}

	return response.Result.MessageID
}
	
func (b *Bot) SendMessage(chatID int64, text string, utils Utils) int {
	if len(text) > 10000 {
		text = text[:10000] + "..."
	}

	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
		"parse_mode": "HTML",
	}

	if utils.Reply != nil {
		message["reply_markup"] = utils.Reply
	}
	if utils.Inline != nil {
		message["reply_markup"] = utils.Inline
	}
	if utils.ReplyMessage != nil {
		message["reply_to_message_id"] = *utils.ReplyMessage
	}
	if utils.MessageThreadID != nil {
		message["message_thread_id"] = *utils.MessageThreadID
	}
	
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return 0
	}

	url := "https://api.telegram.org/bot" + b.Token + "/sendMessage"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error creating request:", err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	var response ResponsePostMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	if !response.Ok {
		return 0
	}

	return response.Result.MessageID
}

func (b *Bot) DownloadFile(fileID string) (io.Reader, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", b.Token, fileID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var fileResponse FileResponse
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return nil, err
	}

	filePath := fileResponse.Result.FilePath
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.Token, filePath)

	resp2, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	return resp2.Body, nil
}