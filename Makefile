all: web backend

web:
	cd cloudscan/CloudScan-WEB/ && npm install && CI='' npm run build

backend:
	go build .

clean:
	rm -f sign-your-horse
	rm -rf cloudscan/CloudScan-WEB/build cloudscan/CloudScan-WEB/node_modules