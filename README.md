# reacjirouter

A modern replacement for [Reacji Channeler](https://reacji-channeler.builtbyslack.com/).

## Slack configuration

Feature | Category | Scope required | Description
--- | --- | ---
`chat.postMessage` | API call | `chat:write` | posts links, interacts with user for config
`reaction_added` | bot event subscription | reactions:read | provides info and trigger for link posts
`message.im` | bot event subscription | im:history | read config commands

## TODO

- receive reacji events
- post message links
  - needs emoji-channel map
    - v0: hard-coded
    - v1: user-configurable
      - v0: in memory
      - v1: persistent
    - interface
      - bot DM
        - `help`
        - `list`
        - `:emoji: #channel`
        - `add`
      - slash command (do we really want this?)
        - `/reacjirouter list`
          - pagination
        - `/reacjirouter :emoji: #channel`
          - invite self to public channels
          - warning if channel isn't joinable
        - `/reacjirouter` (help)
      - app home interface TBD
- persist reacji events to job queue
- OAuth
- deal with visibility of reactions somehow
  - Block Kit channel picker which formats a copy/paste-able invitation string
- new features
  - templatized posts, e.g. "{{User}} posted {{reaction}} in {{Channel}}: {{Link}}"
  - allow multiple channels per emoji
