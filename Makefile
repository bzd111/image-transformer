build:
	docker build -t zidy/image-transform:0.6 .
	docker push zidy/image-transform:0.6 

apply:	
	kubectl apply -f deploy

delete:	
	kubectl delete -f deploy

