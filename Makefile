.PHONY: dev-server ngrok

dev-server:
	reflex -r '\.go$$' -s -- sh -c "go build && ./reacjirouter"

ngrok:
	ngrok http -subdomain=`cat ./config.json | jq -r '.ngrokSubdomain'` 1234
