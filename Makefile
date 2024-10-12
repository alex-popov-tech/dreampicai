build: build-templ build-tailwind build-go
build-templ:
	@echo '--------------------------------------'
	templ generate view/
	@echo '--------------------------------------'
build-tailwind:
	@echo '--------------------------------------'
	npx tailwindcss -i ./view/index.css -o ./public/styles.css --minify
	@echo '--------------------------------------'
build-go:
	@echo '--------------------------------------'
	go build -tags dev -o ./tmp/main .
	@echo '--------------------------------------'

live:
	make -j4 live/templ live/tailwind live/server
live/templ:
	@echo '--------------------------------------'
	templ generate -watch -proxy="http://localhost:3000" -proxyport=3001 -open-browser=false ./view
	@echo '--------------------------------------'
live/tailwind:
	@echo '--------------------------------------'
	npx tailwindcss -i ./view/index.css -o ./public/styles.css --watch
	@echo '--------------------------------------'
live/server:
	@echo '--------------------------------------'
	go run github.com/air-verse/air@latest \
		--build.cmd "go build -tags dev -o ./tmp/bin/main ." \
	  --build.bin "tmp/bin/main" --build.delay "20" \
		--build.include_dir "handler,model,view,utils,pkg" \
		--build.include_file "main.go" \
		--build.log "build-errors.log" \
		--build.stop_on_error false \
		--misc.clean_on_exit true
	@echo '--------------------------------------'
