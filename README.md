# Cretor Back — Backend API

Backend API for the Cretor real estate platform.

## Overview

RESTful API that handles lead submissions, sales pipeline management, user authentication, email delivery, and server-side conversion tracking. Deployed as AWS Lambda functions behind API Gateway.

## Features

- **Lead Management** — Create, read, update, delete leads with status tracking
- **Sales Pipeline** — Kanban pipeline with stages, notes, and automatic lead-to-pipeline creation
- **Authentication** — JWT-based auth with refresh tokens, role-based access (admin, manager, agent)
- **Email** — Transactional emails via AWS SES with PDF attachments
- **Meta Conversion API** — Server-side event tracking for Facebook/Instagram lead events
- **Multi-tenant** — Tenant isolation via X-Tenant-ID header
- **Rate Limiting** — Public lead endpoint limited to 5 requests/minute per IP
- **Bot Protection** — Honeypot field validation on lead forms

## Tech Stack

- **Language:** Go 1.24
- **Router:** Chi v5
- **Database:** MongoDB Atlas (encrypted at rest)
- **Auth:** JWT (golang-jwt)
- **Cloud:** AWS Lambda + API Gateway
- **Email:** AWS SES v2
- **Logging:** Uber Zap

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | /leads | Public (rate-limited) | Create lead from landing page form |
| GET | /leads | Protected | List leads with search/filter |
| GET | /leads/:id | Protected | Get lead by ID |
| PUT | /leads/:id | Protected | Update lead |
| DELETE | /leads/:id | Protected | Delete lead |
| GET | /pipeline | Protected | List pipeline entries |
| POST | /pipeline | Protected | Create pipeline entry |
| PUT | /pipeline/:id | Protected | Update pipeline entry |
| PATCH | /pipeline/:id/move | Protected | Move to pipeline stage |
| POST | /auth/login | Public | Login |
| POST | /auth/signup | Public | Register |
| POST | /auth/refresh | Public | Refresh JWT token |
| POST | /email/send | Protected | Send email with attachment |

## Related Repositories

- [cretor-front](https://github.com/ponte-tech/cretor-front) — Public website
- [cretor-front-admin](https://github.com/ponte-tech/cretor-front-admin) — CRM admin panel
- [cretor-ads-automation](https://github.com/ponte-tech/cretor-ads-automation) — Ads automation CLI
