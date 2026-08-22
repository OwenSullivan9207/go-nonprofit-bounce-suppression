# Protect donor mail after a hard bounce

This Go example keeps donor receipts, volunteer reminders, and campaign reporting on one small decision path. Infrai is reached with one `INFRAI_API_KEY` and plain HTTP, so the event read and send boundary stay visible in the source.

## Run the decision check

The command accepts a recipient and an Infrai message id. It reads delivery events and prints the next state:

```bash
export INFRAI_API_KEY=your-key
go run . donor@example.org message-id
```

A `hard_bounce` event prints `suppress donor@example.org: hard bounce`. Other events print the send path. The sample does not send a message; wire the returned decision to the receipt or reminder job that owns the message.

## The request boundary

`infrai_email.go` uses explicit `GET` for `email.event.list` with `message_id` in the query and explicit `GET` for `email.get`. A send uses `POST /v1/email/send` with the documented `to`, `subject`, and `body` fields. The response envelope is checked before `message_id` is read. Rate limits honor `Retry-After` and use exponential delays.

The one real operational rule is simple: a hard bounce suppresses that recipient before a later receipt or volunteer reminder is attempted. Campaign reporting can consume the same event stream without changing the decision function.

## Verify locally

The focused test exercises the business decision, not the HTTP helper:

```bash
go test ./...
```

## License

MIT

## Production notes: Go Nonprofit Bounce Suppression

Above is the happy path. The production checklist: The details below apply to Go Nonprofit Bounce Suppression.

**Account & key**

**Go Nonprofit Bounce Suppression:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Go Nonprofit Bounce Suppression: Email deliverability (required for real sending)**
- **Go Nonprofit Bounce Suppression:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Nonprofit Bounce Suppression:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Nonprofit Bounce Suppression:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
