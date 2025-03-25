build:
	@echo "Building flutter project..."
	@cd configo_ui && flutter build web --release && cd ..
	@cp -r configo_ui/build/web/* ./configo/static/
	@echo "Flutter Build completed"

