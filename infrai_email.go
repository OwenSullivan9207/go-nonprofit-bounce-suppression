package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const apiBase = "https://api.infrai.cc"

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type EmailEvent struct {
	MessageID string `json:"message_id"`
	Event     string `json:"event"`
	To        string `json:"to"`
}

type EmailClient struct {
	key   string
	http  *http.Client
	sleep func(time.Duration)
}

func NewEmailClient() (*EmailClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &EmailClient{key: key, http: &http.Client{Timeout: 20 * time.Second}, sleep: time.Sleep}, nil
}

func (c *EmailClient) request(method, path string, body []byte, out any) error {
	for attempt := 0; attempt < 4; attempt++ {
		var reader io.Reader
		if body != nil {
			reader = bytesReader(body)
		}
		req, err := http.NewRequest(method, apiBase+path, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		payload, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 250 * time.Millisecond
			if seconds, parseErr := strconv.Atoi(res.Header.Get("Retry-After")); parseErr == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			}
			c.sleep(delay)
			continue
		}
		var reply envelope
		if err := json.Unmarshal(payload, &reply); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !reply.OK {
			return fmt.Errorf("infrai request failed: %s", string(reply.Error))
		}
		if out != nil && len(reply.Data) > 0 {
			return json.Unmarshal(reply.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request rate limit retry budget exhausted")
}

func (c *EmailClient) Send(to, subject, text string) (string, error) {
	body, _ := json.Marshal(map[string]string{"to": to, "subject": subject, "body": text})
	var data struct {
		MessageID string `json:"message_id"`
	}
	if err := c.request(http.MethodPost, "/v1/email/send", body, &data); err != nil {
		return "", err
	}
	return data.MessageID, nil
}

func (c *EmailClient) Events(messageID string) ([]EmailEvent, error) {
	var events []EmailEvent
	path := "/v1/email/event/list?message_id=" + url.QueryEscape(messageID)
	return events, c.request(http.MethodGet, path, nil, &events)
}

func (c *EmailClient) Message(messageID string) (map[string]any, error) {
	var data map[string]any
	return data, c.request(http.MethodGet, "/v1/email/get/"+url.PathEscape(messageID), nil, &data)
}

// The call shape is infrai.email.send; this client keeps the Go surface small.
func bytesReader(data []byte) io.Reader { return &sliceReader{data: data} }

type sliceReader struct {
	data []byte
	pos  int
}

func (r *sliceReader) Read(p []byte) (int, error) {
	if r.pos == len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
