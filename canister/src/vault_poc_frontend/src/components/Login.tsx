import { AuthClient } from '@dfinity/auth-client'
import { ActorSubclass, HttpAgent, Identity } from '@dfinity/agent'
import { Button } from '@mui/material'
import * as React from 'react'

import { createActor } from '../../../declarations/vault_poc_backend'
import { _SERVICE } from '../../../declarations/vault_poc_backend/vault_poc_backend.did'

type Props = {
  setActor: (actor: ActorSubclass<_SERVICE>) => void
  setIdentity: (identity: Identity) => void
}

export const Login: React.FunctionComponent<Props> = ({ setActor, setIdentity }) => {
  const doLogin = async () => {
    // create an auth client
    let authClient = await AuthClient.create()

    // start the login process and wait for it to finish
    await new Promise<void>(resolve => {
      authClient.login({
        identityProvider: process.env.II_URL,
        onSuccess: resolve,
      })
    })

    // At this point we're authenticated, and we can get the identity from the auth client:
    const identity = authClient.getIdentity()
    // Using the identity obtained from the auth client, we can create an agent to interact with the IC.
    const agent = new HttpAgent({ identity })
    // Using the interface description of our webapp, we create an actor that we use to call the service methods.
    let actor = createActor(process.env.VAULT_POC_BACKEND_CANISTER_ID, {
      agent,
    })

    setActor(actor)
    setIdentity(identity)
  }

  return <Button onClick={doLogin}>Login</Button>
}
