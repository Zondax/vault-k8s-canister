---
title: "Development"
sidebar_position: 1
---

## Pre-requisites

1. [DFX](https://internetcomputer.org/docs/current/references/cli-reference/)
1. [Make](https://www.gnu.org/software/make/) the build automation tool - most likely you will already have it but just in case.
1. [Rust](https://www.rust-lang.org/learn/get-started)

## Steps

### Running local
1. Create a new identity with `dfx identity new <name>`
1. Use this new identity with `dfx identity use <name>`
1. If needed, you can get the identity principal with `dfx identity get-principal`
1. Start local ICP replica using dfx with `make start_env`
1. On a new terminal, deploy the local internet_identity, backend and frontend canisters with `make`
1. At this point, you should be able to see the frontend and backend canister id on the terminal. You can use Chrome or Firefox to navigate the frontend
