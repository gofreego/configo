build:
	@echo "Building flutter project..."
	@cd web && flutter build web --release && cd ..
	@cp -r web/build/web/* ./configo/static/
	@echo "Flutter Build completed"

