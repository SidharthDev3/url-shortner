# URL Shortener — System Design

## What is this project?
A full-stack web app that converts long URLs into short, shareable links, built with a Go backend, React frontend, and SQLite database.

## Why we want this project?
Long URLs are hard to share, remember, and use across social media, messages, or print. A shortener makes links compact, clean, and easy to distribute — without losing the original destination.

## What it will do?
- Convert any long URL into a unique short code
- Instantly redirect visitors from the short link to the original URL
- Show a dashboard of all previously created links
- Let users copy a short link with one click
- Detect duplicate URLs and reuse the existing short code instead of creating a new one

## How the system works?
1. User submits a long URL from the frontend.
2. Backend checks if the URL already exists — if yes, returns the existing short code.
3. If not, backend generates a new unique short code and saves the mapping (`short_code → original_url`) in SQLite.
4. The short link is shown to the user with a copy button, and added to the dashboard.
5. When someone visits the short link, the backend looks up the code and redirects them to the original URL.