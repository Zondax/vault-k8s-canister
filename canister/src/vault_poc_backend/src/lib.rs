use candid::{CandidType, Principal};
use hex;
use ic_cdk::api::{caller as caller_api, data_certificate, set_certified_data, time};
use ic_cdk::storage;
use ic_cdk_macros::*;
use serde::{Deserialize, Serialize};
use sha256;
use std::cell::RefCell;
use std::mem;
use json_canon::to_string;

const PENDING_OPERATION_ENUM_CREATE: u8 = 1;
const PENDING_OPERATION_ENUM_UPDATE: u8 = 2;
const PENDING_OPERATION_ENUM_NONE: u8 = 3;

// const PermissionTypeEnumReadAndWrite: u8 = 1;
// const PermissionTypeEnumReadOnly: u8 = 2;

#[derive(Clone, CandidType, Serialize, Deserialize)]
pub struct InitArgs {
    admin: Principal,
}

#[derive(Clone, CandidType, Serialize, Deserialize)]
pub struct CertifiedStatus {
    certificate: Option<Vec<u8>>,
    consumers: Vec<Consumer>,
    pending_consumer_reqs: Vec<Consumer>,
    pending_secret_reqs: Vec<Secret>,
    secrets: Vec<Secret>,
    data: String,
}


#[derive(Clone, CandidType, Serialize, Deserialize)]
pub struct Secret {
    create_timestamp: String,
    creator: String,
    id: String,
    kube_id: String,
    pending_type: u8,
    ttl: u32,
    update_timestamp: String,
}

#[derive(Clone, CandidType, Serialize, Deserialize)]
pub struct Consumer {
    create_timestamp: String,
    creator: String,
    id: String,
    kube_id: String,
    pending_type: u8,
    permission_type: u8,
    secret_kube_id: String,
    update_timestamp: String,
}

#[derive(Clone, CandidType, Serialize, Deserialize)]
pub struct User {
    create_timestamp: String,
    creator: String,
    id: Principal,
    update_timestamp: String,
    root: bool,
}

#[derive(Clone, CandidType, Serialize, Deserialize)]
struct CanisterState {
    consumers: Vec<Consumer>,
    pending_consumer_reqs: Vec<Consumer>,
    pending_secret_reqs: Vec<Secret>,
    secrets: Vec<Secret>,
    users: Vec<User>,
}

thread_local! {
    pub static PENDING_SECRET_REQS: RefCell<Vec<Secret>> = RefCell::new(Vec::new());
    pub static SECRETS: RefCell<Vec<Secret>> = RefCell::new(Vec::new());
    pub static PENDING_CONSUMER_REQS: RefCell<Vec<Consumer>> = RefCell::new(Vec::new());
    pub static CONSUMERS: RefCell<Vec<Consumer>> = RefCell::new(Vec::new());
    pub static USERS: RefCell<Vec<User>> = RefCell::new(Vec::new());
}

#[init]
fn init(args: InitArgs) {
    ic_cdk::setup();

    USERS.with(|users| {
        users.borrow_mut().push(User {
            root: true,
            id: args.admin,
            creator: args.admin.to_string(),
            create_timestamp: time().to_string(),
            update_timestamp: time().to_string(),
        })
    });

    update_certified_data();
}

// Upgrade Hooks

#[pre_upgrade]
/// The pre_upgrade hook determines anything your canister
/// should do before it goes offline for a code upgrade.
fn pre_upgrade() {
    SECRETS.with(|secrets| {
        PENDING_SECRET_REQS.with(|pending_secret_reqs| {
            USERS.with(|users| {
                PENDING_CONSUMER_REQS.with(|pending_consumer_reqs| {
                    CONSUMERS.with(|consumers| {
                        let old_state = CanisterState {
                            secrets: mem::take(&mut secrets.borrow_mut()),
                            pending_secret_reqs: mem::take(&mut pending_secret_reqs.borrow_mut()),
                            consumers: mem::take(&mut consumers.borrow_mut()),
                            pending_consumer_reqs: mem::take(&mut pending_consumer_reqs.borrow_mut()),
                            users: mem::take(&mut users.borrow_mut()),
                        };

                        // storage::stable_save is the API used to write canister state out.
                        // More explicit error handling *can* be useful, but if we fail to read out/in stable memory on upgrade
                        // it means the data won't be accessible to the canister in any way.
                        storage::stable_save((old_state,)).unwrap();
                    })
                })
            })
        })
    });

    update_certified_data();
}

#[post_upgrade]
/// The post_upgrade hook determines anything your canister should do after it restarts
fn post_upgrade() {
    let (old_state,): (CanisterState,) = storage::stable_restore().unwrap();
    SECRETS.with(|secrets| {
        PENDING_SECRET_REQS.with(|pending_secret_reqs| {
            USERS.with(|users| {
                PENDING_CONSUMER_REQS.with(|pending_consumer_reqs| {
                    CONSUMERS.with(|consumers| {
                        *secrets.borrow_mut() = old_state.secrets;
                        *pending_secret_reqs.borrow_mut() = old_state.pending_secret_reqs;
                        *consumers.borrow_mut() = old_state.consumers;
                        *pending_consumer_reqs.borrow_mut() = old_state.pending_consumer_reqs;
                        *users.borrow_mut() = old_state.users;
                    })
                })
            })
        })
    });

    update_certified_data();
}

// Canister methods

#[update(name = "whoami")]
fn whoami() -> String {
    caller_api().to_string()
}

#[update]
fn add_secret(kube_id: String, ttl: u32) -> Result<(), String> {
    let secret_exists =
        SECRETS.with_borrow(|secrets| secrets.iter().any(|s| s.kube_id.eq(&kube_id)));

    if secret_exists {
        return Err("use update_secret fn because secret already exists".to_string());
    }

    PENDING_SECRET_REQS.with(|pending_secret_reqs| {
        let mut writer = pending_secret_reqs.borrow_mut();

        writer.push(Secret {
            id: time().to_string(),
            kube_id,
            creator: caller().to_string(),
            pending_type: PENDING_OPERATION_ENUM_CREATE,
            create_timestamp: time().to_string(),
            update_timestamp: time().to_string(),
            ttl,
        });
    });

    update_certified_data();

    Ok(())
}


#[update]
fn add_consumer(kube_id: String, secret_kube_id: String, permission_type: u8) -> Result<(), String> {
    let secret_exists =
        SECRETS.with_borrow(|secrets| secrets.iter().any(|s| s.kube_id.eq(&secret_kube_id)));
    let consumer_exists =
        CONSUMERS.with_borrow(|consumers| consumers.iter().any(|c| c.kube_id.eq(&kube_id) && c.secret_kube_id.eq(&secret_kube_id)));

    if !secret_exists {
        return Err("create consumer first with add_consumer fn".to_string());
    }

    if consumer_exists {
        return Err("consumer/secret already exists, use update_consumer fn".to_string());
    }

    PENDING_CONSUMER_REQS.with(|pending_consumer_reqs| {
        let mut writer = pending_consumer_reqs.borrow_mut();

        writer.push(Consumer {
            id: time().to_string(),
            kube_id,
            secret_kube_id,
            permission_type,
            creator: caller().to_string(),
            pending_type: PENDING_OPERATION_ENUM_CREATE,
            create_timestamp: time().to_string(),
            update_timestamp: time().to_string(),
        });
    });

    update_certified_data();

    Ok(())
}

#[update]
fn update_consumer(
    kube_id: String, secret_kube_id: String, permission_type: u8
) -> Result<(), String> {
    let (existing_consumer_id, existing_consumer_creator) = CONSUMERS
        .with_borrow(|consumer| {
            consumer
                .iter()
                .find(|s| s.kube_id.eq(&kube_id) &&  s.secret_kube_id.eq(&secret_kube_id))
                .map(|s| (s.id.clone(), s.creator.clone())) // Extract the needed fields directly
        })
        .ok_or("consumer not found".to_string())?; // Handle the case where the consumer is not found


    PENDING_CONSUMER_REQS.with(|pending_consumer_reqs| {
        let mut writer = pending_consumer_reqs.borrow_mut();
        writer.push(Consumer {
            id: existing_consumer_id,
            secret_kube_id,
            kube_id,
            creator: existing_consumer_creator,
            permission_type,
            pending_type: PENDING_OPERATION_ENUM_UPDATE,
            create_timestamp: time().to_string(),
            update_timestamp: time().to_string(),
        });
    });

    update_certified_data();

    Ok(())
}

#[update]
fn update_secret(
    kube_id: String,
    ttl: u32,
) -> Result<(), String> {
    let (existing_secret_id, existing_secret_creator) = SECRETS
        .with_borrow(|secrets| {
            secrets
                .iter()
                .find(|s| s.kube_id.eq(&kube_id))
                .map(|s| (s.id.clone(), s.creator.clone())) // Extract the needed fields directly
        })
        .ok_or("secret not found".to_string())?; // Handle the case where the secret is not found


    PENDING_SECRET_REQS.with(|pending_secret_reqs| {
        let mut writer = pending_secret_reqs.borrow_mut();
        writer.push(Secret {
            id: existing_secret_id,
            kube_id,
            creator: existing_secret_creator,
            pending_type: PENDING_OPERATION_ENUM_UPDATE,
            create_timestamp: time().to_string(),
            update_timestamp: time().to_string(),
            ttl,
        });
    });

    update_certified_data();

    Ok(())
}

#[update(guard = "registered_users_only")]
fn approve_secret(id: String) -> Result<(), String> {
    match SECRETS.with(|secrets_ref| {
        PENDING_SECRET_REQS.with(|pending_secret_reqs_ref| {
            let mut pending_secrets_writer = pending_secret_reqs_ref.borrow_mut();
            let mut secrets_writer = secrets_ref.borrow_mut();
            for (_index, pending_secret) in pending_secrets_writer.iter().enumerate() {
                if pending_secret.id.eq(&id) {
                    let mut new_secret = pending_secret.clone();
                    new_secret.pending_type = PENDING_OPERATION_ENUM_NONE;

                    let pending_type = pending_secret.pending_type.clone();
                    let pending_kube_id = pending_secret.kube_id.clone();

                    if pending_type == PENDING_OPERATION_ENUM_UPDATE {
                        secrets_writer.retain(|s| s.kube_id.ne(&pending_kube_id));
                    }

                    secrets_writer.push(new_secret);
                    pending_secrets_writer.retain(|s| s.kube_id.ne(&pending_kube_id));
                    return Ok(());
                }
            }
            return Err("secret not found".to_string());
        })
    }) {
        Ok(_) => (),
        Err(e) => return Err(e.to_string()),
    };

    update_certified_data();

    Ok(())
}

#[update(guard = "registered_users_only")]
fn approve_consumer(id: String) -> Result<(), String> {
    match CONSUMERS.with(|consumers_ref| {
        PENDING_CONSUMER_REQS.with(|pending_consumers_reqs_ref| {
            let mut pending_consumer_reqs_writer = pending_consumers_reqs_ref.borrow_mut();
            let mut consumers_writer = consumers_ref.borrow_mut();
            for (_index, pending_consumer) in pending_consumer_reqs_writer.iter().enumerate() {
                if pending_consumer.id.eq(&id) {
                    let mut new_consumer = pending_consumer.clone();
                    new_consumer.pending_type = PENDING_OPERATION_ENUM_NONE;

                    let pending_type = pending_consumer.pending_type.clone();
                    let pending_kube_id = pending_consumer.kube_id.clone();
                    let pending_secret_kube_id = pending_consumer.secret_kube_id.clone();

                    if pending_type == PENDING_OPERATION_ENUM_UPDATE {
                        consumers_writer.retain(|s| !(s.kube_id.eq(&pending_kube_id) && s.secret_kube_id.eq(&pending_secret_kube_id)));
                    }

                    consumers_writer.push(new_consumer);
                    pending_consumer_reqs_writer.retain(|s| !(s.kube_id.eq(&pending_kube_id) && s.secret_kube_id.eq(&pending_secret_kube_id)));
                    return Ok(());
                }
            }
            return Err("consumer not found".to_string());
        })
    }) {
        Ok(_) => (),
        Err(e) => return Err(e.to_string()),
    };

    update_certified_data();

    Ok(())
}


#[update(guard = "registered_users_only")]
fn revoke_secret(id: String) -> Result<(), String> {
    SECRETS.with(|secrets_ref| {
        let mut writer = secrets_ref.borrow_mut();
        let secret_index = writer
            .iter()
            .position(|s| s.id.eq(&id))
            .ok_or("secret not found")
            .unwrap();

        writer.remove(secret_index);
    });

    update_certified_data();

    Ok(())
}


#[update(guard = "registered_users_only")]
fn revoke_consumer(id: String) -> Result<(), String> {
    CONSUMERS.with(|consumers_ref| {
        let mut writer = consumers_ref.borrow_mut();
        let consumer_index = writer
            .iter()
            .position(|s| s.id.eq(&id))
            .ok_or("consumer not found")
            .unwrap();

        writer.remove(consumer_index);
    });

    update_certified_data();

    Ok(())
}

#[update(guard = "registered_users_only")]
fn revoke_pending_consumer(id: String) -> Result<(), String> {
    PENDING_CONSUMER_REQS.with(|pending_consumer_reqs_ref| {
        let mut writer = pending_consumer_reqs_ref.borrow_mut();
        let consumer_index = writer
            .iter()
            .position(|s| s.id.eq(&id))
            .ok_or("consumer not found")
            .unwrap();

        writer.remove(consumer_index);
    });

    update_certified_data();

    Ok(())
}

#[update(guard = "registered_users_only")]
fn revoke_pending_secret(id: String) -> Result<(), String> {
    PENDING_SECRET_REQS.with(|secrets_ref| {
        let mut writer = secrets_ref.borrow_mut();
        let secret_index = writer
            .iter()
            .position(|s| s.id.eq(&id))
            .ok_or("secret not found")
            .unwrap();

        writer.remove(secret_index);
    });

    update_certified_data();

    Ok(())
}

#[update(guard = "root_only")]
fn add_privileged_user(id: Principal) -> Result<(), String> {
    USERS.with(|users_ref| {
        let mut writer = users_ref.borrow_mut();
        writer.push(User { id, root: true, creator: caller().to_string(), create_timestamp: time().to_string(), update_timestamp: time().to_string() });
    });

    Ok(())
}

#[update(guard = "registered_users_only")]
fn add_user(id: Principal) -> Result<(), String> {
    USERS.with(|users_ref| {
        let mut writer = users_ref.borrow_mut();
        writer.push(User { id, root: false, creator: caller().to_string(), create_timestamp: time().to_string(), update_timestamp: time().to_string() });
    });

    Ok(())
}

#[update(guard = "root_only")]
fn remove_privileged_user(id: Principal) -> Result<(), String> {
    USERS.with(|users_ref| {
        let mut users_writer = users_ref.borrow_mut();
        users_writer.retain(|u| !(u.id.eq(&id) && u.root == true));
    });

    Ok(())
}
#[update(guard = "registered_users_only")]
fn remove_user(id: Principal) -> Result<(), String> {
    USERS.with(|users_ref| {
        let mut users_writer = users_ref.borrow_mut();
        users_writer.retain(|u| !(u.id.eq(&id) && u.root == false));
    });

    Ok(())
}

#[query]
fn get_certified_status() -> CertifiedStatus {
    let secrets_ref = SECRETS.with(|reference| reference.borrow().clone());
    let consumers_ref = CONSUMERS.with(|reference| reference.borrow().clone());
    let pending_secret_reqs_ref = PENDING_SECRET_REQS.with(|reference| reference.borrow().clone());
    let pending_consumer_reqs_ref = PENDING_CONSUMER_REQS.with(|reference| reference.borrow().clone());

    let secrets = to_string(&secrets_ref).unwrap();
    let consumers = to_string(&consumers_ref).unwrap();
    let pending_secret_reqs = to_string(&pending_secret_reqs_ref).unwrap();
    let pending_consumer_reqs = to_string(&pending_consumer_reqs_ref).unwrap();

    let data = format!("{consumers}{pending_consumer_reqs}{pending_secret_reqs}{secrets}");
    let certificate: Option<Vec<u8>> = data_certificate();
    return CertifiedStatus {
        secrets: secrets_ref,
        consumers: consumers_ref,
        pending_secret_reqs: pending_secret_reqs_ref,
        pending_consumer_reqs: pending_consumer_reqs_ref,
        certificate,
        data,
    };
}

#[query]
fn get_users() -> Vec<User> {
    USERS.with(|users| users.borrow().clone())
}

// Helpers

fn caller() -> Principal {
    let caller = caller_api();
    // The anonymous principal is not allowed to interact with the
    // encrypted notes canister.
    /*if caller == Principal::anonymous() {
        panic!("Anonymous principal not allowed to make calls.")
    }*/
    caller
}

fn root_only() -> Result<(), String> {
    is_user_allowed(true)
}

fn registered_users_only() -> Result<(), String> {
    is_user_allowed(false)
}

fn is_user_allowed(root_required: bool) -> Result<(), String> {
    USERS.with(|users_ref| {
        for user in users_ref.borrow().iter() {
            if user.id.eq(&caller()) && (!root_required || (root_required && user.root)) {
                return Ok(());
            }
        }
        return Err("cannot access this route".to_string());
    })
}

fn update_certified_data() -> () {
    let secrets_ref = SECRETS.with(|reference| reference.borrow().clone());
    let consumers_ref = CONSUMERS.with(|reference| reference.borrow().clone());
    let pending_secrets_reqs_ref = PENDING_SECRET_REQS.with(|reference| reference.borrow().clone());
    let pending_consumers_reqs_ref = PENDING_CONSUMER_REQS.with(|reference| reference.borrow().clone());

    let secrets = to_string(&secrets_ref).unwrap();
    let consumers = to_string(&consumers_ref).unwrap();
    let pending_secrets_reqs = to_string(&pending_secrets_reqs_ref).unwrap();
    let pending_consumers_reqs = to_string(&pending_consumers_reqs_ref).unwrap();

    let result_string_consumers = sha256::digest(consumers);
    let result_string_pending_consumers = sha256::digest(pending_consumers_reqs);
    let result_string_pending_secrets = sha256::digest(pending_secrets_reqs);
    let result_string_secrets = sha256::digest(secrets);
    let result_string = sha256::digest(format!("{result_string_consumers}{result_string_pending_consumers}{result_string_pending_secrets}{result_string_secrets}"));

    let result_hex = hex::decode(result_string).unwrap();
    set_certified_data(result_hex.as_slice());
}

// Enable Candid export
ic_cdk::export_candid!();
