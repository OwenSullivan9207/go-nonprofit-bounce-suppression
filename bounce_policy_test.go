package main

import "testing"

func TestHardBounceSuppressesMatchingRecipient(t *testing.T) {
	decision := DecideBounce(EmailEvent{Event: "hard_bounce", To: "donor@example.org"}, DonorMessage{Email: "donor@example.org", Kind: "receipt"})
	if decision.Send || !decision.AddToSuppress {
		t.Fatalf("expected suppression decision, got %+v", decision)
	}
}

func TestOtherEventRemainsSendable(t *testing.T) {
	decision := DecideBounce(EmailEvent{Event: "delivered", To: "donor@example.org"}, DonorMessage{Email: "donor@example.org", Kind: "reminder"})
	if !decision.Send || decision.AddToSuppress {
		t.Fatalf("expected send decision, got %+v", decision)
	}
}
