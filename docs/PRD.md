# Product Requirements Document (PRD)

# FreeTokensPoker

**Version:** 1.0

**Status:** Draft

**Author:** Kunal Mehra

---

# 1. Product Overview

## Vision

FreeTokensPoker is a collaborative estimation platform designed for engineering teams that increasingly rely on AI models, coding agents, and LLM-powered workflows.

Just as Planning Poker became the standard for estimating story points, FreeTokensPoker introduces a new way for teams to estimate AI-related resources before work begins.

Instead of asking:

> "How many story points is this task?"

Teams can now ask:

* How many AI tokens will this require?
* Which AI model is most suitable?
* What will the approximate AI cost be?
* How many engineering days will AI-assisted implementation require?

Every team member estimates independently, followed by a simultaneous reveal to encourage unbiased discussion and better decision-making.

---

# 2. Background

AI is rapidly becoming a fundamental engineering tool.

Developers today use AI for:

* Writing code
* Code reviews
* Documentation
* Architecture
* Debugging
* Testing
* Refactoring
* Research

However, AI usage today is largely untracked during planning.

Organizations discuss:

* Story Points
* Sprint Capacity
* Engineering Time

But rarely discuss:

* AI Budget
* Token Consumption
* Model Selection
* Agent Usage

This exists primarily because AI is currently inexpensive due to aggressive market competition.

As AI becomes core infrastructure rather than a productivity experiment, engineering organizations will require structured planning around AI resource consumption.

FreeTokensPoker is built around this future.

---

# 3. Problem Statement

Current engineering planning ignores AI resources.

Example:

A Product Manager creates a feature.

Developers immediately begin discussing:

* Claude
* GPT
* Gemini
* Cursor
* Codex

without answering:

* Which model should be used?
* Is the expensive model necessary?
* How many tokens might be consumed?
* What budget should be allocated?
* Is AI even the correct solution?

These decisions happen ad hoc.

There is no collaborative estimation process.

---

# 4. Solution

FreeTokensPoker introduces AI Planning Poker.

The workflow remains intentionally familiar.

1. Create Room
2. Invite Team
3. Create Estimation Task
4. Everyone votes privately
5. Reveal simultaneously
6. Discuss differences
7. Save final decision

This encourages independent thinking before group discussion.

---

# 5. Product Goals

Primary Goals

* Enable collaborative AI estimation
* Reduce bias during planning
* Encourage engineering discussions around AI costs
* Create historical estimation records
* Keep the experience frictionless

Secondary Goals

* Build long-term AI estimation datasets
* Improve AI budgeting
* Enable future analytics
* Become the Planning Poker standard for AI-native engineering teams

---

# 6. Target Audience

Primary

Engineering Teams

* Software Engineers
* Tech Leads
* Engineering Managers
* Product Managers

Secondary

Organizations building AI-powered software.

Future

* AI Consulting Companies
* AI Agencies
* Enterprise Engineering Organizations
* Remote Engineering Teams

---

# 7. Core Philosophy

The product should remain extremely simple.

No dashboards.

No complicated setup.

No enterprise onboarding.

A user should reach estimation in under one minute.

Simplicity is a feature.

---

# 8. User Journey

Landing Page

↓

Create Room

or

Join Room

↓

Enter Email

↓

Receive Login Link / OTP

↓

Enter Room

↓

PM Creates Task

↓

Everyone Estimates

↓

Reveal

↓

Discussion

↓

Save Results

---

# 9. Authentication

Authentication intentionally avoids passwords.

Users identify themselves only using email.

Reasons

* Lower friction
* Faster onboarding
* Better meeting experience
* No password management

Future enterprise authentication may include:

* Google Workspace
* Microsoft Entra ID
* SSO

---

# 10. Room System

A room represents a collaborative estimation session.

Room contains

* Members
* Active Task
* Previous Tasks
* Reveal State

One room can contain multiple estimation tasks.

---

# 11. Estimation Modes

Version 1 supports four estimation modes.

## Mode 1

AI Tokens

Example Cards

500K

1M

2M

5M

10M

20M

50M

---

## Mode 2

AI Cost

Example

$1

$5

$10

$25

$50

$100

$250

---

## Mode 3

Engineering Days

Example

1

2

3

5

8

13

21

---

## Mode 4

Best AI Model

Cards become model names instead of numbers.

Example

GPT

Claude

Gemini

DeepSeek

Cursor

Codex

Other

---

# 12. Estimation Flow

PM creates task.

Task becomes visible.

Everyone receives cards.

Each participant secretly selects one.

Nobody can see others' selections.

PM presses Reveal.

All estimates appear simultaneously.

Discussion begins.

PM selects Final Decision.

Task is archived.

---

# 13. History

Every estimation is saved.

Stored information

* Task Name
* Description
* Mode
* Participants
* Individual Estimates
* Final Estimate
* Timestamp
* Room

Users can revisit previous estimations.

---

# 14. Future Analytics

Historical estimations enable future features.

Examples

Team AI spending trends

Most selected models

Average token estimates

Prediction accuracy

Department comparisons

Budget forecasting

No analytics are included in Version 1.

---

# 15. Why This Product Exists

Planning Poker standardized estimation for engineering effort.

AI introduces a new planning dimension.

Engineering organizations will increasingly need to estimate AI resources before implementation.

FreeTokensPoker extends Planning Poker into the AI era.

The product is not attempting to replace Planning Poker.

It extends it.

---

# 16. Competitive Advantage (Moat)

## 1. First-Mover Position

Very few tools focus specifically on collaborative AI estimation.

Most estimation products still revolve around story points.

---

## 2. Familiar Workflow

Engineering teams already understand Planning Poker.

Minimal learning curve.

---

## 3. AI-Native Design

The product is designed specifically for:

* Tokens
* Cost
* AI Models
* Agents

rather than adapting legacy planning software.

---

## 4. Historical AI Dataset

Over time the platform can accumulate anonymized estimation data.

Potential insights include:

* Average token usage
* Preferred models
* Cost trends
* Estimation accuracy

This dataset becomes increasingly valuable.

---

## 5. Future Integrations

Possible integrations include:

GitHub

Linear

Jira

Slack

Cursor

Claude Code

OpenAI

Anthropic

Google AI

Engineering workflows become richer as AI platforms mature.

---

# 17. Why "Free"

The name reflects accessibility rather than pricing.

The goal is to make AI planning easy to adopt.

Future enterprise features may be monetized without changing the product philosophy.

---

# 18. Version 1 Scope

Included

✓ Landing page

✓ Create Room

✓ Join Room

✓ Email authentication

✓ Room management

✓ Create estimation task

✓ Four estimation modes

✓ Card selection

✓ Reveal

✓ Final decision

✓ History

Not Included

✗ Billing

✗ AI integrations

✗ Token tracking

✗ Analytics

✗ Organization management

✗ Admin dashboard

✗ SSO

✗ Notifications

---

# 19. Success Metrics

Short-Term

* Number of rooms created
* Number of completed estimation sessions
* Average session duration
* Repeat users

Long-Term

* Weekly active engineering teams
* Enterprise adoption
* Historical estimation dataset growth
* AI planning sessions per organization

---

# 20. Long-Term Vision

FreeTokensPoker aims to become the standard collaborative estimation platform for AI-native engineering teams.

As AI evolves from an optional productivity tool into a core engineering resource, teams will require structured planning around model selection, token consumption, cost, and AI-assisted execution.

The long-term opportunity extends beyond estimation into analytics, governance, budgeting, and AI resource management.

The initial product intentionally focuses on doing one thing exceptionally well:

Helping teams align on AI decisions before work begins.
