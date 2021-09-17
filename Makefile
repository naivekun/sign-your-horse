all: web backend

web:
	cd cloudscan/CloudScan-WEB/ && npm install && CI='' npm run build

backend:
	go build .

clean:
	rm sign-your-horse
	rm -rd cloudscan/CloudScan-WEB/build cloudscan/CloudScan-WEB/node_modules