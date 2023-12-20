deploy: deploy_identity deploy_backend generate_candid deploy_frontend

deploy_identity:
	dfx deploy internet_identity


deploy_backend:
	# " RUSTFLAGS="--cfg tokio_unstable" " is added because of tokio wasm compilation process. If this is not added as rust flags, the process will fail.
	RUSTFLAGS="--cfg tokio_unstable" dfx deploy vault_poc_backend --argument '(record {admin = principal "$(ADMIN_PRINCIPAL)"})' --argument-type idl

deploy_frontend:
	# " RUSTFLAGS="--cfg tokio_unstable" " is added because of tokio wasm compilation process. If this is not added as rust flags, the process will fail.
	RUSTFLAGS="--cfg tokio_unstable" dfx deploy vault_poc_frontend

generate_candid:
	candid-extractor target/wasm32-unknown-unknown/release/vault_poc_backend.wasm > src/vault_poc_backend/vault_poc_backend.did

# https://internetcomputer.org/docs/current/developer-docs/backend/candid/generating-candid
install_candid_extractor:
	cargo install candid-extractor

create_secret:
	dfx canister call $(BACKEND_CANISTER_ID) add_secret '("secret_id_1", 11131 : nat32)'

update_secret:
	dfx canister call -v $(BACKEND_CANISTER_ID) update_secret '("secret_id_1": text, 1000 : nat32)'

create_consumers:
	dfx canister call $(BACKEND_CANISTER_ID) add_consumer '("consumer_1", "secret_id_1", 1)'
	dfx canister call $(BACKEND_CANISTER_ID) add_consumer '("consumer_2", "secret_id_1", 2)'

get_certified_status:
	dfx canister call $(BACKEND_CANISTER_ID) get_certified_status '()'

update_consumers:
	dfx canister call $(BACKEND_CANISTER_ID) add_consumer '("consumer_1", "secret_id_1", 2)'
	dfx canister call $(BACKEND_CANISTER_ID) add_consumer '("consumer_2", "secret_id_1", 1)'

add_privileged_user:
	dfx canister call $(BACKEND_CANISTER_ID) add_privileged_user '(principal "$(NEW_USER)")'

start_env:
	dfx start

update_env:
	dfx upgrade
	dfx start --clean

dfx_list_identities:
	dfx identity list

dfx_get_principal:
	dfx identity get-principal