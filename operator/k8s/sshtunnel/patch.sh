#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Generate a keypair
yes | ssh-keygen -q -P "" -t ed25519 -f "${SCRIPT_DIR}/dev_ed25519" >/dev/null 2>&1
chmod 600 "${SCRIPT_DIR}/dev_ed25519"

# Prepare config map with SSH keys and apply
public_key=$(sed 's/^/    /' "${SCRIPT_DIR}/dev_ed25519.pub")
sed "s|KEY_GOES_HERE|$public_key|" "${SCRIPT_DIR}/configmap-template.yaml" > "${SCRIPT_DIR}/configmap.yaml"
kubectl apply -n development -f "${SCRIPT_DIR}/configmap.yaml"

## Now patch the deployment
kubectl patch -n development deployment tororu-operator-golem --type=strategic --patch-file="${SCRIPT_DIR}/patch.yaml"
