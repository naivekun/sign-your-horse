all: web backend

web:
	cd cloudscan/CloudScan-WEB/ && npm install && CI='' npm run build

backend:
	go build .