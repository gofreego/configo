build:
	@echo "Building flutter project..."
	@cd configo_ui && flutter build web --release && cd ..
	@cp -r configo_ui/build/web/* ./configo/internal/ui/static/
	@echo "Flutter Build completed"

setup:
	@echo "instaling dependencies..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	@echo "Dependencies installed"

doc:
	@echo "Generating API documentation..."
	@swag init 
	@echo "API documentation generated"
