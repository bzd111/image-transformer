build:
	docker build -t zidy/image-transform:0.3 .

apply:	
	kubectl apply -f deploy

delete:	
	kubectl delete -f deploy

