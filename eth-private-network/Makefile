build:
	cd go-ethereum && make geth

clean:
	rm -rf build/
	cd go-ethereum && make clean

docker-build:
	docker build -f docker/Dockerfile.miner . -t brandonjroberts/eth-private-miner:latest
	docker build -f docker/Dockerfile.rpc . -t brandonjroberts/eth-private-rpc:latest
	docker build -f docker/Dockerfile.node-setup . -t brandonjroberts/eth-private-node-setup:latest

docker-push:
	docker push brandonjroberts/eth-private-miner:latest
	docker push brandonjroberts/eth-private-rpc:latest
	docker push brandonjroberts/eth-private-node-setup:latest

generate-account:
	./scripts/generate-account.sh -d ${HOME}/l1-miner-data/ -x

run-miner: generate-account
	./scripts/run-miner.sh -d ${HOME}/l1-miner-data/ -x

run-miner-daemon: generate-account
	./scripts/run-miner.sh -d ${HOME}/l1-miner-data/ -x -o ${OUTPUT_FILE}

run-miner-noclean:
	./scripts/run-miner.sh -d ${HOME}/l1-miner-data/

docker-generate-account:
	mkdir -p ${HOME}/docker-l1-miner-data/
	./scripts/generate-account.sh -d ${HOME}/docker-l1-miner-data/ -x

docker-run-miner: docker-generate-account
	docker run -it -p 30303:30303 -p 8545:8545 -p 8551:8551 -p 30303:30303/udp -p 6060:6060 \
		-v ${HOME}/docker-l1-miner-data/:/l1-miner-data/ \
		-v ${HOME}/.eth-accounts/:/root/.eth-accounts/ \
		--network host --rm --name eth-private-miner brandonjroberts/eth-private-miner:latest

docker-run-miner-noclean:
	docker run -it -p 30303:30303 -p 8545:8545 -p 8551:8551 -p 30303:30303/udp -p 6060:6060 \
		-v ${HOME}/docker-l1-miner-data/:/l1-miner-data/ \
		-v ${HOME}/.eth-accounts/:/root/.eth-accounts/ \
		--network host --rm --name eth-private-miner brandonjroberts/eth-private-miner:latest

setup-rpc-node:
	-rm -rf ${HOME}/l1-rpc-data/
	mkdir -p ${HOME}/l1-rpc-data/
	./scripts/setup-rpc.sh -d ${HOME}/l1-rpc-data/ -g ${HOME}/l1-miner-data/genesis.json

run-rpc: setup-rpc-node
	./scripts/run-rpc.sh -d ${HOME}/l1-rpc-data/ -m 30306 -r 8550 -H 8548

run-rpc-daemon: setup-rpc-node
	./scripts/run-rpc.sh -d ${HOME}/l1-rpc-data/ -m 30306 -r 8550 -H 8548 -o ${OUTPUT_FILE}

run-rpc-noclean:
	./scripts/run-rpc.sh -d ${HOME}/l1-rpc-data/ -m 30306 -r 8550 -H 8548

docker-setup-rpc-node:
	mkdir -p ${HOME}/docker-l1-rpc-data/
	docker cp eth-private-miner:/l1-miner-data/genesis.json ${HOME}/docker-l1-rpc-data/genesis.json
	./scripts/setup-rpc.sh -d ${HOME}/docker-l1-rpc-data/ -g ${HOME}/docker-l1-rpc-data/genesis.json

#TODO: connect rpc with different account easily
docker-run-rpc: docker-setup-rpc-node
	docker run -it -p 30306:30306 -p 8548:8548 -p 8550:8550 -p 30306:30306/udp -p 6061:6061 \
		-v ${HOME}/docker-l1-rpc-data/:/l1-rpc-data/ --network host \
		-v ${HOME}/.eth-accounts/:/root/.eth-accounts/ \
		--rm --name eth-private-rpc brandonjroberts/eth-private-rpc:latest

docker-run-rpc-noclean:
	docker run -it -p 30306:30306 -p 8548:8548 -p 8550:8550 -p 30306:30306/udp -p 6061:6061 \
		-v ${HOME}/docker-l1-rpc-data/:/l1-rpc-data/ --network host \
		-v ${HOME}/.eth-accounts/:/root/.eth-accounts/ \
		--rm --name eth-private-rpc brandonjroberts/eth-private-rpc:latest

connect-peers:
	./scripts/connect-peers.sh -P http://localhost:8548 -H http://localhost:8545

kube-connect-peers:
	kubectl port-forward service/eth-private-miner-app-service 8545 30303 6060 &
	kubectl port-forward service/eth-private-rpc-app-service 8548 30306 6061 &
	sleep 5
	./scripts/connect-peers.sh -P http://localhost:8548 -H http://localhost:8545 -i $(shell kubectl get pod/eth-private-rpc-app-0 --template '{{.status.podIP}}')

kube-deploy-miner:
	kubectl apply -f kube/miner-node-deploy-sts.yaml

kube-deploy-rpc:
	kubectl cp eth-private-miner-app-0:/l1-miner-data/genesis.json ${HOME}/genesis.json
	kubectl create configmap genesisconfigmap --from-file=${HOME}/genesis.json
	kubectl apply -f kube/rpc-node-deploy-sts.yaml

kube-clean:
	-kubectl delete statefulset.apps/eth-private-miner-app
	-kubectl delete statefulset.apps/eth-private-rpc-app
	-kubectl delete configmap genesisconfigmap
	-ps aux | grep kubectl\ port-forward | awk '{print $$2}' | xargs kill -9

kube-clean-all:
	-kubectl delete statefulset.apps/eth-private-miner-app
	-kubectl delete statefulset.apps/eth-private-rpc-app
	-kubectl delete configmap genesisconfigmap
	-kubectl delete service/eth-private-miner-app-service
	-kubectl delete service/eth-private-rpc-app-service
	-kubectl delete pvc/data-eth-private-miner-app-0
	-kubectl delete pvc/account-eth-private-miner-app-0
	-kubectl delete pvc/data-eth-private-rpc-app-0
	-kubectl delete pvc/account-eth-private-rpc-app-0
	-ps aux | grep kubectl\ port-forward | awk '{print $$2}' | xargs kill -9

kube-reset:
	-ps aux | grep kubectl\ port-forward | awk '{print $$2}' | xargs kill -9
	-kubectl delete pod/eth-private-miner-app-0
	-kubectl delete pod/eth-private-rpc-app-0

kube-reset-all:
	-ps aux | grep kubectl\ port-forward | awk '{print $$2}' | xargs kill -9
	-kubectl delete pvc/data-eth-private-miner-app-0
	-kubectl delete pvc/account-eth-private-miner-app-0
	-kubectl delete pvc/data-eth-private-rpc-app-0
	-kubectl delete pvc/account-eth-private-rpc-app-0
	-kubectl delete pod/eth-private-miner-app-0
	-kubectl delete pod/eth-private-rpc-app-0
