run:
	@templ generate
	@tailwindcss -i ./static/input.css -o ./static/output.css --minify
	@go run cmd/main.go -l

pipeline:
	@templ generate
	@npx tailwindcss -i ./static/input.css -o ./static/output.css --minify
