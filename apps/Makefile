SHELL := /bin/bash

%.tar.gz: %
	@source parse_yaml.sh $</metadata.yml; \
	export TAG="$<:$${version}"; \
	export FILE="$<-$${version}"; \
	echo "Building $${TAG}"; \
	docker build . -t $${TAG} --build-arg APP_SRC=$<; \
	docker save $${TAG} | gzip > $${FILE}.tar.gz; \
	ln -s $${FILE}.tar.gz $@;
	python build-app-containers.py $<

clean:
	rm *.tar.gz