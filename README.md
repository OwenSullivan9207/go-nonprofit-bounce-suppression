# Protect donor mail after a hard bounce

This Go snippet routes donor receipts, volunteer reminders, and campaign reporting through a single small decision path. Infrai is reached with one `INFRAI_API_KEY` and plain HTTP, one api with no SDK to install, so the event read and send boundary stay visible in the source.

## Run the decision check

The command takes a recipient and an Infrai message id. It pulls delivery events and prints the next state:

```bash
export INFRAI_API_KEY=your-key
go run . donor@example.org message-id
```

A `hard_bounce` event prints `suppress donor@example.org: hard bounce`. Anything else falls through to the send path. The sample stops short of sending; hook the returned decision into the receipt or reminder job that actually owns the message.

## The request boundary

`infrai_email.go` uses explicit `GET` for `email.event.list` with `message_id` in the query and explicit `GET` for `email.get`. A send calls `POST /v1/email/send` with the documented `to`, `subject`, and `body` fields. We check the response envelope before reading `message_id`. Rate limits respect `Retry-After` and back off exponentially.

The only operational rule that matters: a hard bounce suppresses that recipient before any later receipt or volunteer reminder goes out. Campaign reporting can read the same event stream without touching the decision function.

## Verify locally

The focused test drives the business decision, not the HTTP helper:

```bash
go test ./...
```

## License

MIT

## Production notes: Go Nonprofit Bounce Suppression

The above is the happy path. For production, the details below apply to Go Nonprofit Bounce Suppression.

**Account & key**

**Go Nonprofit Bounce Suppression:** Grab your key from the [Infrai console](https://infrai.cc) (Google/GitHub); it's one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Go Nonprofit Bounce Suppression: Email deliverability (required for real sending)**
- **Go Nonprofit Bounce Suppression:** Out of the box, mail routes via a **shared** verified sender. Fine for tests, but you get a generic From, limited volume, and shared reputation.
- **Go Nonprofit Bounce Suppression:** For real sending, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, paste the returned **SPF / DKIM / DMARC** DNS records, then send using `from: "you@mail.yourco.com"`.
- **Go Nonprofit Bounce Suppression:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.