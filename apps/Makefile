SHELL := /bin/bash

%.tar.gz: %
	@source parse_yaml.sh $</metadata.yml; \
	export TAG="$<:$${Version}"; \
	export FILE="$<-$${Version}"; \
	echo "Building $${TAG}"; \
	docker buildx build --platform linux/arm64 -t gcr.io/hurrabuild/$${TAG} --build-arg APP_SRC=$< -o type=docker,dest=- . | gzip > $${FILE}-arm64.tar.gz; \
	docker buildx build --platform linux/amd64 -t gcr.io/hurrabuild/$${TAG} --build-arg APP_SRC=$< -o type=docker,dest=- . | gzip > $${FILE}-amd64.tar.gz; \
	ln -s $${FILE}-amd64.tar.gz $*-amd64.tar.gz; \
	ln -s $${FILE}-arm64.tar.gz $*-arm64.tar.gz; \
	ln -s $${FILE}-amd64.tar.gz $@;
	python build-app-containers.py $<

clean:
	rm *.tar.gz