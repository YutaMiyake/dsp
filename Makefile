DOCKER_REPO=asia.gcr.io/cyberagent-224

ctr-cluster:
	gcloud container clusters get-credentials hitori-ctr-cluster2 --project cyberagent-224 --zone asia-northeast1-a

dsp-cluster:
	gcloud container clusters get-credentials hitori-dsp-cluster --project cyberagent-224 --zone asia-northeast1-a

# DOCKER TASKS
build:
	docker build -t $(APP_NAME) ./$(APP_NAME)/

build-nc:
	docker build --no-cache -t $(APP_NAME) ./$(APP_NAME)/

publish-latest: tag-latest
	@echo 'publish latest to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):latest

publish-version: tag-version
	@echo 'publish $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

tag-latest:
	@echo 'create tag latest'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):latest

tag-version:
	@echo 'create tag $(VERSION)'
	docker tag $(APP_NAME) $(DOCKER_REPO)/$(APP_NAME):$(VERSION)

gcr-login:
	docker login -u oauth2accesstoken -p `gcloud auth application-default print-access-token` https://asia.gcr.io

deploy:
	kubectl apply -f $(APP_NAME)/deployment.yaml
	kubectl apply -f $(APP_NAME)/service.yaml
	kubectl apply -f $(APP_NAME)/autoscale.yaml

delete: 
	kubectl delete -f $(APP_NAME)/deployment.yaml
	kubectl delete -f $(APP_NAME)/service.yaml
	kubectl delete -f $(APP_NAME)/autoscale.yaml

delete-all:
	kubectl delete deployment,service --all
