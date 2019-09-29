# reacjirouter

A modern replacement for [Reacji Channeler](https://reacji-channeler.builtbyslack.com/).

## TODO

- receive reacji events
- post message links
- deferred queue for channeler
- deal with visibility somehow
  - Block Kit channel picker which formats a copy/paste-able invitation string
- bot interface
  - `help`
  - `list`
  - `:emoji: #channel`
  - `add`
- slash command interface (do we really want this?)
  - `/reacjirouter list`
    - pagination
  - `/reacjirouter :emoji: #channel`
    - invite self to public channels
    - warning if channel isn't joinable
  - `/reacjirouter` (help)
- app home interface
