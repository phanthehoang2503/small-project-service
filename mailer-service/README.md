## mailer-service

This service handles email notifications.

### Overview

The `mailer-service` listens for system events (like Order Created) and sends emails to users. Uses **MailHog** for local development to capture emails without sending them to the real world.

### Run locally

From repository root:

```powershell
cd mailer-service/cmd/api
go run .
```

### MailHog UI

View sent emails at: http://localhost:8025
