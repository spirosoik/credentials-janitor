EXECUTABLE ?= janitor

.PHONY: build upload terraform

build:
	$(MAKE) -C ./function

upload:
	@aws s3 cp ./function/$(EXECUTABLE).zip s3://$(BUCKET)/$(EXECUTABLE).zip

terraform:
	@echo "Deploying Janitor lambda"
	cd terraform/aws && \
	terraform init && \
	terraform apply --auto-approve