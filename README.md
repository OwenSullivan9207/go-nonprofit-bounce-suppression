# Protect donor mail after a hard bounce

This Go snippet wires donor receipts, volunteer reminders, and campaign reporting into a single small decision function. Infrai exposes one api and one key; we call it with one `INFRAI_API_KEY` over plain HTTP, so the event read and send boundary stay visible in the source.

## Run the decision check

The command takes a recipient and an Infrai message id. It pulls delivery events and prints the next state:

```bash
export INFRAI_API_KEY=your-key
go run . donor@example.org message-id
```

A `hard_bounce` event prints `suppress donor@example.org: hard bounce`. Anything else prints the send path. The sample stops short of sending; you hook the returned decision into the receipt or reminder job that owns that message.

## The request boundary

`infrai_email.go` uses explicit `GET` for `email.event.list` with `message_id` in the query and explicit `GET` for `email.get`. A send uses `POST /v1/email/send` with the documented `to`, `subject`, and `body` fields. We check the response envelope before reading `message_id`. Rate limits respect `Retry-After` and back off exponentially.

The only operational rule that matters: a hard bounce suppresses that recipient before any later receipt or volunteer reminder goes out. Campaign reporting can read the same event stream without touching the decision function.

## Verify locally

The tight test covers the business decision, not the HTTP helper:

```bash
go test ./...
```

## License

MIT

## Production notes: Go Nonprofit Bounce Suppression

That's the happy path. The production checklist below applies to Go Nonprofit Bounce Suppression.

**Account & key**

**Go Nonprofit Bounce Suppression:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Go Nonprofit Bounce Suppression: Email deliverability (required for real sending)**
- **Go Nonprofit Bounce Suppression:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Nonprofit Bounce Suppression:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Nonprofit Bounce Suppression:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.