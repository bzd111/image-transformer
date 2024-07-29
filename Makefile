build:
	docker build -t zidy/image-transform:0.1 .

deploy:	
	kubectl apply -f deploy

clean:	
	kubectl clean -f deploy

