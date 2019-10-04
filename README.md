# reacjirouter

A modern replacement for [Reacji Channeler](https://reacji-channeler.builtbyslack.com/).

## Slack configuration

Feature | Category | Scope required | Description
--- | --- | --- | ---
`auth.test` | API call | - | get team URL and user ID
`chat.postMessage` | API call | `chat:write` | posts links, interacts with user for config
`conversations.info` | API call | `channels:read`, `groups:read`, etc. | check bot's membership
`reaction_added` | bot event subscription | `reactions:read` | provides info and trigger for link posts

## TODO

- add actual persistence
  - tokens
  - routes
- defer reaction handling to JQ
- make sure to ignore reacji from private channels and/or DMs
  - or maybe this is safe?
- new features
  - templatized posts, e.g. "{{User}} posted {{reaction}} in {{Channel}}: {{Link}}"
  - allow multiple channels per emoji
