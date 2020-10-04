SHELL := /bin/bash

%.tar.gz: %
	@source parse_yaml.sh $</metadata.yml; \
	export TAG="$<:$${Version}"; \
	export FILE="$<-$${Version}"; \
	echo "Building $${TAG}"; \
	docker buildx build --platform linux/arm64,linux/amd64  -t gcr.io/hurrabuild/$${TAG} --push --build-arg APP_SRC=$< .; \
	docker pull --platform linux/arm64 gcr.io/hurrabuild/$${TAG}; \
	docker tag gcr.io/hurrabuild/$${TAG} $${TAG}; \
	docker save $${TAG} | gzip > $${FILE}-arm64.tar.gz; \
	docker rmi $${TAG}; \
	docker pull --platform linux/amd64 gcr.io/hurrabuild/$${TAG}; \
	docker tag gcr.io/hurrabuild/$${TAG} $${TAG}; \
	docker save $${TAG} | gzip > $${FILE}-amd64.tar.gz; \
	docker rmi $${TAG}; \
	ln -s $${FILE}-amd64.tar.gz $*-amd64.tar.gz; \
	ln -s $${FILE}-arm64.tar.gz $*-arm64.tar.gz; \
	ln -s $${FILE}-amd64.tar.gz $@;
	python build-app-containers.py $<

clean:
	rm *.tar.gz