package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . <donor-email> <message-id>")
		os.Exit(2)
	}
	client, err := NewEmailClient()
	if err != nil {
		panic(err)
	}
	events, err := client.Events(os.Args[2])
	if err != nil {
		panic(err)
	}
	decision := SuppressionDecision{Send: true, Reason: "no hard bounce observed"}
	for _, event := range events {
		decision = DecideBounce(event, DonorMessage{Email: os.Args[1], Kind: "receipt"})
		if decision.AddToSuppress {
			break
		}
	}
	if decision.AddToSuppress {
		fmt.Printf("suppress %s: %s\n", os.Args[1], decision.Reason)
		return
	}
	fmt.Printf("send receipt or volunteer reminder to %s: %s\n", os.Args[1], decision.Reason)
}
