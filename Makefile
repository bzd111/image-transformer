build:
	docker build -t zidy/image-transform:0.5 .
	docker push zidy/image-transform:0.5 

apply:	
	kubectl apply -f deploy

delete:	
	kubectl delete -f deploy

