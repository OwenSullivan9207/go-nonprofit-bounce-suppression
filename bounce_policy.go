package main

import "strings"

type DonorMessage struct {
	Email string
	Kind  string
}

type SuppressionDecision struct {
	Send          bool
	AddToSuppress bool
	Reason        string
}

func DecideBounce(event EmailEvent, message DonorMessage) SuppressionDecision {
	if strings.EqualFold(event.Event, "hard_bounce") && strings.EqualFold(event.To, message.Email) {
		return SuppressionDecision{AddToSuppress: true, Reason: "hard bounce"}
	}
	return SuppressionDecision{Send: true, Reason: "recipient is deliverable"}
}
