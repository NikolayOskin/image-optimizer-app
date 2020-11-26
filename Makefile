build:
	docker build -f prod.Dockerfile --rm -t leakso86/tinyimg .

push:
	docker push leakso86/tinyimg

clean:
	docker image prune --filter label=stage=builder