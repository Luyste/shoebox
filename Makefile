.PHONY: dev

dev:
	@trap 'kill 0' EXIT; \
	templ generate --watch --proxy="http://localhost:3000" & \
	tailwindcss --input static/input.css --output static/output.css --watch & \
	go run . & \
	wait
