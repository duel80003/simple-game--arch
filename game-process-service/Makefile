format:
	buf format -w proto/v1

lint:
	buf lint proto/v1 --config buf.yaml

generate:
	buf generate proto/v1

build_dev:
	docker build --build-arg ENV_FILE=dev -t tg_read_service .