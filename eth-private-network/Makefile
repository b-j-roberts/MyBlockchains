build:
	cd go-ethereum && make geth

docker-build:
	docker build -f docker/Dockerfile.miner . -t brandonjroberts/eth-private-miner:latest
	docker build -f docker/Dockerfile.rpc . -t brandonjroberts/eth-private-rpc:latest

docker-push:
	docker push brandonjroberts/eth-private-miner:latest
	docker push brandonjroberts/eth-private-rpc:latest

docker-run-miner:
	docker run -it -p 30303:30303 -p 8545:8545 -p 8551:8551 -p 30303:30303/udp \
		--network host --rm --name eth-private-miner brandonjroberts/eth-private-miner:latest

docker-run-rpc:
	mkdir -p build/
	-rm -f build/genesis.json
	docker cp eth-private-miner:/l1-miner-data/genesis.json build/
	docker run -it -p 30306:30306 -p 8548:8548 -p 8550:8550 -p 30306:30306/udp \
		--mount type=bind,source="$(shell pwd)"/build/genesis.json,target=/genesis.json,readonly \
		--network host --rm --name eth-private-rpc brandonjroberts/eth-private-rpc:latest

connect-docker:
	./scripts/connect-docker-peers.sh -P /l1-miner-data/geth.ipc -I /l1-rpc-data/geth.ipc

generate-account:
	./scripts/generate-account.sh -d ${HOME}/l1-miner-data/ -x

launch-miner-local:
	./scripts/generate-account.sh -d ${HOME}/l1-miner-data/ -x
	./scripts/launch-miner.sh -d ~/l1-miner-data/ -x

launch-miner-local-daemon:
	./scripts/generate-account.sh -d ${HOME}/l1-miner-data/ -x
	./scripts/launch-miner.sh -d ~/l1-miner-data/ -x -o ${OUTPUT_FILE}

launch-miner-local-noclean:
	./scripts/launch-miner.sh -d ~/l1-miner-data/

launch-rpc-local:
	./scripts/launch-rpc.sh -d ~/l1-rpc-data/ -x -g ~/l1-miner-data/genesis.json -m 30306 -r 8550 -H 8548

launch-rpc-local-daemon:
	./scripts/launch-rpc.sh -d ~/l1-rpc-data/ -x -g ~/l1-miner-data/genesis.json -m 30306 -r 8550 -H 8548 -o ${OUTPUT_FILE}

launch-rpc-local-noclean:
	./scripts/launch-rpc.sh -d ~/l1-rpc-data/ -g ~/l1-miner-data/genesis.json -m 30306 -r 8550 -H 8548

connect-local:
	./scripts/connect-peers.sh -P ~/l1-rpc-data/geth.ipc -H ~/l1-miner-data/geth.ipc

kube-deploy-miner:
	kubectl apply -f kube/miner-node-deploy-sts.yaml

kube-deploy-rpc:
	kubectl cp eth-private-miner-app-0:/l1-miner-data/genesis.json build/genesis.json
	kubectl create configmap genesisconfigmap --from-file=build/genesis.json
	kubectl apply -f kube/rpc-node-deploy-sts.yaml

kube-clean:
	-kubectl delete statefulset.apps/eth-private-miner-app
	-kubectl delete statefulset.apps/eth-private-rpc-app
	-kubectl delete configmap genesisconfigmap

kube-clean-all:
	-kubectl delete statefulset.apps/eth-private-miner-app
	-kubectl delete statefulset.apps/eth-private-rpc-app
	-kubectl delete configmap genesisconfigmap
	-kubectl delete service/eth-private-miner-app-service
	-kubectl delete service/eth-private-rpc-app-service
	-kubectl delete pvc/data-eth-private-miner-app-0
	-kubectl delete pvc/data-eth-private-rpc-app-0

kube-reset:
	-kubectl delete pod/eth-private-miner-app-0
	-kubectl delete pod/eth-private-rpc-app-0

kube-reset-all:
	-kubectl delete pvc/data-eth-private-miner-app-0
	-kubectl delete pvc/data-eth-private-rpc-app-0
	-kubectl delete pod/eth-private-miner-app-0
	-kubectl delete pod/eth-private-rpc-app-0

clean:
	rm -rf build/
	cd go-ethereum && make clean

#TODO: connect-kube, kubectl port-forward service/eth-private-miner-app-service 8545 30303
