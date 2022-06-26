module broker-service

go 1.18

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	toolbox v0.0.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.0 // indirect
	golang.org/x/net v0.0.0-20210505024714-0287a6fb4125 // indirect
)

replace toolbox v0.0.0 => ../toolbox
