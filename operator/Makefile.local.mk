helm-install:
	helm upgrade \
	--install $(APP_NAME) oci://registry.zondax.ch/zondax/golem \
	-f k8s/sshtunnel/values.yaml \
	--create-namespace -n development
	sleep 2
	./k8s/sshtunnel/patch.sh

helm-uninstall:
	helm uninstall $(APP_NAME) -n development

NGROK_URL = $(shell curl -s localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url')
RNDHASH = $(shell head -c 100 /dev/urandom | md5sum | cut -d' ' -f1)

tunnel-adm-controller:
	@ngrok http 8282

tunnel-icp:
	@ngrok http 4943

k8s-mock:
	@echo "NGROK_URL: $(NGROK_URL)"
	@echo "RNDHASH: $(RNDHASH)"
	@kubectl get namespace development > /dev/null 2>&1 || kubectl create namespace development
	@yq eval -i '.webhooks.[].clientConfig.url="$(NGROK_URL)"' k8s/mock1/mutating-webhook.yaml
	@yq eval -i '.metadata.annotations."tororu.io/hash"="$(RNDHASH)"' k8s/mock1/mutating-webhook.yaml
	@yq eval -i '.spec.template.metadata.annotations."tororu.io/hash"="$(RNDHASH)"' k8s/mock1/deployment.yaml
#	kubectl apply -f k8s/mock1/mutating-webhook.yaml -n development
#	kubectl apply -f k8s/mock1/deployment.yaml -n development
#	kubectl apply -f k8s/mock1/secret.yaml -n development
	kubectl apply -f k8s/mock1 -n development

build-images:
	make earthly

# https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/
generate-tls-files:
	# generate CA
	mkdir -p tmp
	openssl genrsa -des3 -out tmp/myCA.key 2048
	openssl req -x509 -new -nodes -key tmp/myCA.key -sha256 -days 1825 -out tmp/myCA.pem
	# generate CA-signed certificates
	echo "authorityKeyIdentifier=keyid,issuer \nbasicConstraints=CA:FALSE \nkeyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment \n subjectAltName = @alt_names \n\n[alt_names] \nDNS.1 = tororu-operator\nDNS.2 = tororu-operator.default\nDNS.3 = tororu-operator.default.svc\nDNS.4 = tororu-operator.default.svc.cluster\nDNS.5 = tororu-operator.default.svc.cluster.local\nDNS.7 = tororu-operator.default.svc.cluster.local\nDNS.8 = 127.0.0.1" > tmp/tororu.operator.ext
	openssl genrsa -out tmp/tororu.operator.key 2048
	openssl req -new -key tmp/tororu.operator.key -out tmp/tororu.operator.csr
	openssl x509 -req -in tmp/tororu.operator.csr -CA tmp/myCA.pem -CAkey tmp/myCA.key -CAcreateserial -out tmp/tororu.operator.crt -days 825 -sha256 -extfile tmp/tororu.operator.ext
	openssl base64 -A -in tmp/myCA.pem -out tmp/myCA.base64.pem
	openssl base64 -A -in tmp/tororu.operator.crt -out tmp/tororu.operator.base64.crt
	openssl base64 -A -in tmp/tororu.operator.key -out tmp/tororu.operator.base64.key


install-operator-chart:
	 helm uninstall tororu-operator || true
	 helm install tororu-operator charts/tororu-operator

install-crds-chart:
	 helm uninstall postgres-user-prateek || true
	 helm uninstall postgres-user-juan || true
	 helm uninstall tororu-crds || true
	 helm install tororu-crds charts/tororu-crd


template-operator-chart:
	 helm template charts/tororu-operator > tmp/operator-chart-output.yaml

template-crds-chart:
	 helm template charts/tororu-crd > tmp/operator-chart-output.yaml

start-local:
	CANISTER_ID="d6g4o-amaaa-aaaaa-qaaoq-cai" \
	ICP_NODE_URL="http://127.0.0.1:4943" \
	ADM_CONTROLLER_CERT_BASE64="LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVVVENDQXptZ0F3SUJBZ0lVTVBjUkZoamJoVENwam5LdTYraHRpMUtyZWprd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TXpFeU1UVXhORFEyTURaYUZ3MHlOakF6Ck1Ua3hORFEyTURaYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0VpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ3J6SndHQ2Y1RmlSNjZoOGRvNVZlbVNHd0dhWlc1cHZTdWhNUlRqbmpmCjk3Um53TFRjM1RmaGVrcUZESXRmT3hTeG9leXk4Z3M5djV1My94ZndZSElWNlZWV1pLQWRoSGZ5cm5VQk5iaWgKQ3Frb3R1Lzh0RjJienYwREtmUjlQRjRFdVk0N0M2N0NrU1J4QlMwWWtVb0JuYXZUTEJhSVkzemJyQUI2VlVsRAp6KzZRUkVFeW0vU25QeDZzWHJxTncyQi80RG9sT0V2djU2NmpIYk5LNzVyalN5MkU5ZWtuVXdVcGJVWXEyT29lCllyN1k3bjRuQ2l0ZGN6aTRWZll1SkZtSXF2YkxuTlNoOHpFa0QwNUtkV1lQTjV5alZJa2JTWm9HZjhhN1lwdnQKbUw2VmZIUll4c3M1VW9JMTdHb2FSM05LTzFQVHBMQ3owL1JCMThaRmxWNmRBZ01CQUFHamdnRTNNSUlCTXpBZgpCZ05WSFNNRUdEQVdnQlNIL1pYU0xFRk5DVTFCUkV0bzRzNWtEOFlYV0RBSkJnTlZIUk1FQWpBQU1Bc0dBMVVkCkR3UUVBd0lFOERDQjJBWURWUjBSQklIUU1JSE5nZzkwYjNKdmNuVXRiM0JsY21GMGIzS0NGM1J2Y205eWRTMXYKY0dWeVlYUnZjaTVrWldaaGRXeDBnaHQwYjNKdmNuVXRiM0JsY21GMGIzSXVaR1ZtWVhWc2RDNXpkbU9DSTNSdgpjbTl5ZFMxdmNHVnlZWFJ2Y2k1a1pXWmhkV3gwTG5OMll5NWpiSFZ6ZEdWeWdpbDBiM0p2Y25VdGIzQmxjbUYwCmIzSXVaR1ZtWVhWc2RDNXpkbU11WTJ4MWMzUmxjaTVzYjJOaGJJSXBkRzl5YjNKMUxXOXdaWEpoZEc5eUxtUmwKWm1GMWJIUXVjM1pqTG1Oc2RYTjBaWEl1Ykc5allXeUNDVEV5Tnk0d0xqQXVNVEFkQmdOVkhRNEVGZ1FVMWUzZApxVUdPZWdWcUY0b3d6aVk0a2ZnNG1kTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBRzliWU91bTJ6M1oreHFhClhCRmF3YnhOYURLWElGNTYzUmV3UTZjYTJmUDlGM0V1NHUwSWQxU1hnMTNmbHdEV1Rpd3F5VzhKTW9HY1NGc04KTVhoNFZBNWJaVjNVWE9mMDF6Zk5HeVUzQzRQUGFqUUpTM09rNC9Sb0FldnZFZ1BuV2NFSG1nSHMzVDhmSWh6cwoydW1Ncjl1U3UvV0JCOVM0ektzMVBabEwwWXVMQ1U3UU5uR2owbDM4NXBxUEMvMjJKNjVBK0xhaXdTTS9MdkxwCldBNFRZdmRpdnBJdWE0SXJlL3JkYXNITi95Unp0VU1jRW1TM2ZqTy83Q2ZIS1B1aHJpMFo5K2tqT0RkWDQ2SlAKbUNEQlBsQVBwVnFPeTRsMm5tYkZnQnBKN1lpM1ZKSEJQdkFPcm5ER0tDa29mN1lGWllBbjQ4TFZTdlJZVm9KegpwK0l5TFdVPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==" \
	ADM_CONTROLLER_KEY_BASE64="LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRQ3J6SndHQ2Y1RmlSNjYKaDhkbzVWZW1TR3dHYVpXNXB2U3VoTVJUam5qZjk3Um53TFRjM1RmaGVrcUZESXRmT3hTeG9leXk4Z3M5djV1MwoveGZ3WUhJVjZWVldaS0FkaEhmeXJuVUJOYmloQ3Frb3R1Lzh0RjJienYwREtmUjlQRjRFdVk0N0M2N0NrU1J4CkJTMFlrVW9CbmF2VExCYUlZM3pickFCNlZVbER6KzZRUkVFeW0vU25QeDZzWHJxTncyQi80RG9sT0V2djU2NmoKSGJOSzc1cmpTeTJFOWVrblV3VXBiVVlxMk9vZVlyN1k3bjRuQ2l0ZGN6aTRWZll1SkZtSXF2YkxuTlNoOHpFawpEMDVLZFdZUE41eWpWSWtiU1pvR2Y4YTdZcHZ0bUw2VmZIUll4c3M1VW9JMTdHb2FSM05LTzFQVHBMQ3owL1JCCjE4WkZsVjZkQWdNQkFBRUNnZ0VBRk9zNUpGTWJMd1JmUlg4Ni9MN1FTV01RSkVlKy8zZ2cydzgzaUtVVWV0RUMKbW8rUWRrUkpoWjhLYStEM0o0VmVJN0wveTFwRm5DTTBwdGJjNTF3WENDdjlSQ1BFaTFPUjkyN2V1R0wrTkQzRQpFejBUUThZQ2ovSklSSlpiT3RTYTdpQlovVDZTN1FZWFZkdTNmZ0pTN0pkeVVLaFJwaEhYSmpodlpuWDBFZG1RCmV3WHBIN3YrZTZkNjF2cUdMV3QyWmZFQXpHQ2NyNFBlZUpFdFcvdmhEZHN2aWJFQkZjNVhEUDJPVjUva3FqRGkKUU0yUm10TEkyL1pYeWYwcVdnWHFnWHgxb1V3OTZsVmRBQkxadG9tSzdEa290VCt3aysvQnE3VjA0TktFaTZ3WgpEMmZ4ak1EQUNiay9VM1Z6aDk1WUE3cGRORGdIcVorQ1kxVWQ4ZFJ1MFFLQmdRRG1VS1pVN1BPbXRwNmtUZUtQCkRFSkwrcnpXNnMzeGMrUUJkOSs3dy9JUldxMHpmbEV2SFhoT0FKRWVzODd2cWZuZ0Y0L3VKeFdSWjBTeUlNRngKd1hJdUJmRUZ5aGoxaUlyWUJRalAyRmZ3THVUTFpFT3BDTHBUMFRGcHdWZDJXK3FPL3BBR0FFSWVrVm92eUhmZApzekFkMUJOa2gyRHEveC9BaHhNRDZiNE15d0tCZ1FDKzlWN3BweVUzd0hHUGwycHNzZFFjQzhRVjcwMXllZkdFCi9aSTJQU3FMWng0ZjJqb2hvYTBpVkRhRGJ1SVBzR29Dd21iVGUrYWY1NVFwN3hCMjJUOHRkR001bEVFS3BVem8KcW9QT1dGaXlwTWV3RC9XYzJDbTRqeWllT0FxM3lOWjBnNFBScE11Y1I3QXNDWGZxemZHOVZiNHJKQStNcmF3Kwo3ajI1aWRUOU53S0JnUUM3cnZFOHQ3TitJY2Q1b1RhRTE3cVc3QWRESkNrYklCT24xcVh0L3ltZVZzUlorQS8wClV4R2NqdjJ6aFZlWEdtN1QzSitmdFIzd1ZiVTNhMVg0ZTJtdWM0MEw2THNhSzJEcDFJQnZ6NThwelMwSlNmV1IKSTltalFCQUNYRm9IeTdPRFA1TGlNUWV3blVaZk5mL29ISU9UYXlVNmdNL0w4SWRSZjBGUnFRTUVyUUtCZ0FzOQp3cEthcGxRNzNmT0lCRm5WdGhqWWtIaUNGOXNQVnFwdml2WHFiK0M0OTBzRXU3dFRHekFVS1FsZnM2c2N4WURZCkZObUtSNjlPSUtpL1RBYlREeWNMM1BOOHlMOXByN2Rhb2x1NVU2OWdoK2pUWjdBT0FaYTl4clJadERmUmVONXYKQjRtRjIvNmRNYi9GNXV0SnFGdHUrcnpyYUlidGltQkNBaHcwQXZmTkFvR0JBSzNjZ3VPWHQ3aDZGdjJXVWltdgpmT0J5MzBjc3lFbCtFcmUzdVplc3paS09ZSEZON01yWnJwblhSWnovWUdJcjE2a2M0OGF2Z0xHci9taDJEamRLClNkbFdLSVVkbVNaeGN0WjBRZ2Z6SURieitqOS9PYmJXSUN4MXFNL3YwWDVhUXI5dW03a1pFZFVQZllIbEZ3NTYKdWdhTnFjbzFsaGNKK0tKNUEvQmhrZE9OCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K" \
	make run

start-chart:
	make build-images
	make install-chart

install-goic:
	go install github.com/aviate-labs/agent-go/cmd/goic@v0.3.3

generate-client:
	goic generate did common/icp/backend.did client --output=common/icp/backend.go --packageName=icp