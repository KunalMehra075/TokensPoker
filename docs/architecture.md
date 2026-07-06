# FreeTokensPoker

# Technical Architecture Document

Version: 1.0

Author: Kunal Mehra

Status: Initial Architecture

---

# 1. Overview

This document describes the complete technical architecture for FreeTokensPoker.

The architecture is designed around five principles:

* Simplicity
* Performance
* Maintainability
* Scalability
* Developer Experience

The MVP intentionally avoids unnecessary complexity while leaving room for enterprise growth.

---

# 2. Technology Stack

## Frontend

* React 19
* Vite
* TypeScript
* Tailwind CSS
* shadcn/ui
* React Router
* TanStack Query
* Zustand
* React Hook Form
* Zod
* Axios
* Socket.IO Client

---

## Backend

* Golang 1.24+
* Gin Framework
* MongoDB
* Socket.IO (Go implementation)
* JWT Authentication
* Redis (Future)
* Docker

---

## Database

MongoDB

Reason:

The application revolves around Rooms, Tasks and embedded voting documents.

MongoDB provides flexibility while allowing embedded documents for estimation history.

---

## Realtime

Socket.IO

Reason:

* Room based architecture
* Automatic reconnection
* Simpler than raw WebSockets
* Well supported in React

Realtime Events

* Join Room
* Leave Room
* User Connected
* User Disconnected
* Vote Submitted
* Vote Changed
* Reveal Votes
* New Task
* Final Decision Selected

---

## Infrastructure

Frontend

Vercel

Backend

Railway / Fly.io / Render

Production

DigitalOcean / Hetzner VPS

Database

MongoDB Atlas

Email

Resend

Monitoring

BetterStack

Analytics

Plausible

Logging

Structured JSON Logs

---

# 3. System Architecture

Browser

↓

React SPA

↓

REST API

↓

Go Backend

↓

MongoDB

Realtime

↓

Socket.IO

↓

Browser

The REST API manages CRUD operations.

Socket.IO handles live collaboration.

---

# 4. Design Principles

The frontend should feel like a professional engineering tool.

Not flashy.

Not startup-themed.

Think

GitHub

Linear

Vercel Dashboard

Stripe Dashboard

The interface should communicate trust and clarity.

---

# 5. Design Language

Primary Font

Manrope

Weights

400

500

600

700

---

Primary Colors

Background

White

Secondary Background

Gray-50

Borders

Gray-200

Text

Gray-900

Secondary Text

Gray-600

Primary Accent

Blue-600

Primary Hover

Blue-700

Success

Green-600

Warning

Amber-500

Danger

Red-600

Cards

White

Buttons

Blue

Hover

Dark Blue

No gradients.

No glassmorphism.

No heavy animations.

---

# 6. Layout

Maximum Width

1280px

Content Centered

Large whitespace

Rounded corners

Medium shadows

Minimal icons

Corporate aesthetic

---

# 7. Frontend Folder Structure

```text
src/

    assets/

    components/

        ui/

        room/

        task/

        voting/

        layout/

    hooks/

    pages/

        Landing

        Login

        Room

        History

    services/

        api.ts

        socket.ts

    store/

    lib/

    types/

    utils/

    constants/

    routes/

```

---

# 8. Backend Folder Structure

```text
cmd/

internal/

    api/

    handlers/

    middleware/

    websocket/

    auth/

    services/

    repositories/

    models/

    dto/

    validators/

    config/

pkg/

main.go
```

Clean Architecture

Handler

↓

Service

↓

Repository

↓

MongoDB

Business logic never lives inside handlers.

---

# 9. Database Collections

## Users

```json
{
    "_id": "",
    "email": "",
    "name": "",
    "createdAt": "",
    "lastLogin": ""
}
```

---

## Rooms

```json
{
    "_id": "",
    "roomCode": "",
    "name": "",
    "ownerId": "",
    "members": [],
    "createdAt": ""
}
```

---

## Tasks

```json
{
    "_id": "",
    "roomId": "",
    "title": "",
    "description": "",
    "mode": "TOKENS",
    "status": "ACTIVE",
    "revealed": false,
    "createdBy": "",
    "createdAt": ""
}
```

---

## Votes

```json
{
    "_id": "",
    "taskId": "",
    "userId": "",
    "selectedCard": "",
    "createdAt": ""
}
```

---

## Final Decisions

```json
{
    "_id": "",
    "taskId": "",
    "finalValue": "",
    "selectedBy": "",
    "createdAt": ""
}
```

---

# 10. Authentication

MVP

Email Login

↓

OTP

↓

JWT

Flow

Enter Email

↓

Receive OTP

↓

Verify OTP

↓

Issue JWT

↓

Store JWT

No passwords.

No OAuth.

Future

Google Workspace

Microsoft

SSO

---

# 11. User Roles

Owner

Can

Create Room

Delete Room

Create Tasks

Reveal Votes

Choose Final Estimate

---

Member

Can

Join Room

Vote

View History

---

Future

Admin

Observer

Guest

---

# 12. API Design

Authentication

POST

/api/auth/request-otp

POST

/api/auth/verify

---

Rooms

POST

/api/rooms

GET

/api/rooms/:id

POST

/api/rooms/join

DELETE

/api/rooms/:id

---

Tasks

POST

/api/tasks

GET

/api/tasks/:id

GET

/api/rooms/:id/tasks

PATCH

/api/tasks/:id/reveal

PATCH

/api/tasks/:id/final

---

Votes

POST

/api/votes

PATCH

/api/votes

GET

/api/tasks/:id/votes

---

History

GET

/api/history

---

# 13. Socket Events

Client

JoinRoom

LeaveRoom

SubmitVote

UpdateVote

RevealVotes

CreateTask

SelectFinalDecision

Server

RoomUpdated

MemberJoined

MemberLeft

VoteReceived

VotesRevealed

TaskCreated

TaskClosed

---

# 14. Estimation Modes

Every estimation mode shares one schema.

```go
type EstimationMode string

const (

TOKENS

COST

DAYS

MODEL

)
```

Cards should NOT be hardcoded.

Every mode stores

Name

Card Values

Display Order

Future versions can support

Complexity

GPU Hours

Confidence

Risk

Latency

without changing code.

---

# 15. State Management

Global

Authentication

Current Room

Current User

Socket

Current Task

Local

Forms

Dialogs

Selected Card

Filters

TanStack Query

Server State

Zustand

Client State

---

# 16. Security

JWT Expiration

Rate Limiting

OTP Expiration

Input Validation

Zod

Go Validation

Helmet Headers

CORS

Room Code Validation

Email Validation

Sanitize Inputs

---

# 17. Performance

React Lazy Loading

Component Memoization

Query Caching

Optimistic Updates

Socket Events

instead of polling

Mongo Indexes

roomCode

email

roomId

taskId

userId

---

# 18. Error Handling

Backend

Consistent JSON

```json
{
    "success": false,
    "message": "...",
    "errorCode": "..."
}
```

Frontend

Toast Notifications

Retry Buttons

Friendly Errors

---

# 19. UI Components

Landing

Navbar

Hero

Features

Footer

---

Room

Members

Task List

Task Header

Voting Cards

Reveal Button

Final Decision

---

Dialogs

Create Room

Join Room

Create Task

OTP

Delete Confirmation

---

Reusable Components

Button

Card

Input

Badge

Dialog

Dropdown

Avatar

Toast

Table

Tabs

Tooltip

Alert

Skeleton

---

# 20. Accessibility

Keyboard Navigation

Visible Focus

Proper Contrast

ARIA Labels

Semantic HTML

Large Click Targets

---

# 21. Responsive Design

Desktop First

Supports

Desktop

Laptop

Tablet

Mobile

Planning meetings primarily occur on desktop.

---

# 22. Logging

Every important action should be logged.

Examples

User Login

Room Created

Task Created

Vote Submitted

Reveal

Final Decision

Logs should be structured JSON.

---

# 23. Future Architecture

Future services

Notification Service

Analytics Service

Organization Service

Billing Service

AI Integration Service

Audit Service

These remain separate modules.

The MVP remains a single Go service.

---

# 24. Coding Standards

Frontend

Functional Components

Strict TypeScript

No Inline Styles

Composition over Inheritance

Small Components

Backend

Context Everywhere

Dependency Injection

Interfaces

Repository Pattern

No Business Logic in Handlers

No Global Variables

---

# 25. Deployment Pipeline

GitHub

↓

Pull Request

↓

CI

Lint

Tests

Build

↓

Merge

↓

Automatic Deploy

↓

Production

---

# 26. Non-Functional Requirements

Authentication Response

< 300ms

API Response

< 200ms

Socket Latency

< 100ms

Cold Start

< 2 sec

Lighthouse

90+

Accessibility

95+

---

# 27. MVP Scope

Included

✓ Email OTP Login

✓ Create Room

✓ Join Room

✓ Room Members

✓ Create Tasks

✓ Four Estimation Modes

✓ Voting

✓ Vote Changes

✓ Reveal

✓ Final Decision

✓ Task History

✓ Realtime Updates

Not Included

✗ Organizations

✗ Analytics

✗ Billing

✗ AI APIs

✗ Slack

✗ Jira

✗ GitHub

✗ Teams

✗ SSO

---

# 28. Engineering Philosophy

FreeTokensPoker should feel like an internal engineering tool that developers immediately trust.

Every interaction should prioritize clarity over decoration.

The interface should be fast, predictable, and unobtrusive, allowing teams to focus on discussion rather than software.

As the product evolves, new capabilities should be added without increasing cognitive load. Features should be modular, optional, and discoverable, ensuring that the simplicity of the MVP remains a defining characteristic even as the platform grows into enterprise AI planning.

Also we should treat this as an event-driven collaboration app. Every action—joining a room, casting a vote, changing a vote, revealing estimates, or selecting a final decision—should be represented as a socket event. That will keep the UI feeling instantaneous and make future features like live cursors, timers, or discussion indicators much easier to add.