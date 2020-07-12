# Craigslist-watcher

Tracks craigslist search queries (RSS links) and alerts to slack if new postings found.

## How to get craigslist search query RSS link

- Go to a craigslist search page; apply your filters
- Near the right bottom, there should be RSS button; right click > "copy link location"
- Add those links to `craigslistWatcher.go` . In this repo, I am using a github config repo to store those links.

## Implementation
- Based on [nightswatch](http://github.com/vvksh/nightswatch); implements `Watcher` interface
- After each interval, it fetches each craigslist RSS feed and checks for new items.
- Each new item is posted to a slack channel specified by return value of `SlackChannel()`

## Alternative solutions:
- I tried to use `RSS` app for slack which fetches RSS feeds and posts to a slack channel; however it returned error with craigslist RSS link.
- This implementation also allows adding custom logic (none yet) on the updates.

## Deploy to docker
- **Generate Dockerfile**
    - [Optional] Replace "{{GITHUB_TOKEN}}" in `Dockerfile-template` with appropriate token for github access
    - Replace "{{SLACK_WEBHOOK}}" in `Dockerfile-template` with your slack webhook url
    - `mv Dockerfile-template Dockerfile`
- **Build and deploy** (shown for swarm mode)
    ```bash
    docker-compose build
    docker-compose push
    docker stack --compose-file docker-compose.yaml  craigslist-watcher
    ```