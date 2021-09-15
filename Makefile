all: web backend

web:
	cd cloudscan/CloudScan-WEB/ && npm install && npm run build

backend:
	go build .