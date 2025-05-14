run:
	@templ generate
	@tailwindcss -i ./static/input.css -o ./static/output.css --minify
	@go run main.go -l

pipeline:
	@templ generate
	@npx tailwindcss -i ./static/input.css -o ./static/output.css --minify
