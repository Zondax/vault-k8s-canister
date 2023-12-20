import * as React from 'react'
import { useEffect, useState } from 'react'
import { Container, Button, FormLabel, Divider, Chip, Alert } from '@mui/material'
import { styled } from '@mui/material/styles'
import { ActorSubclass, Identity } from '@dfinity/agent'

import { vault_poc_backend } from '../../../declarations/vault_poc_backend'
import { _SERVICE, CertifiedStatus, User } from '../../../declarations/vault_poc_backend/vault_poc_backend.did'
import { ActiveSecrets } from './ActiveSecrets'
import { PendingSecrets } from './PendingSecrets'
import { Login } from './Login'
import { PendingConsumers } from './PendingConsumers'
import { ActiveConsumers } from './ActiveConsumers'
import { validateCertificate } from '../lib/certificate'
import { Users } from './Users'

const Root = styled('div')(({ theme }) => ({
  width: '100%',
  ...theme.typography.body2,
  '& > :not(style) + :not(style)': {
    marginTop: theme.spacing(2),
  },
}))

export const App = () => {
  const [actor, setActor] = useState<ActorSubclass<_SERVICE>>(vault_poc_backend)
  const [identity, setIdentity] = useState<Identity | null>(null)
  const [status, setStatus] = useState<CertifiedStatus>({ secrets: [], consumers: [], pending_consumer_reqs: [], pending_secret_reqs: [], certificate: [], data: '' })
  const [users, setUsers] = useState<User[]>([])
  const [statusValidated, setStatusValidated] = useState<boolean | null>(null)

  useEffect(() => {
    loadSecretsData()
  }, [])

  const loadSecretsData = async () => {
    const certifiedStatusResp = await actor.get_certified_status()
    const users = await actor.get_users()
    const result = await validateCertificate(certifiedStatusResp)

    setStatus(certifiedStatusResp)
    setUsers(users)
    setStatusValidated(result)
  }

  const logout = () => {
    setIdentity(null)
    setActor(vault_poc_backend)
  }

  return (
    <Root>
      <Container maxWidth="lg">
        <div style={{ alignContent: 'flex-end' }}>
          <FormLabel>{!!identity ? identity.getPrincipal().toString() : 'Anonymous'}</FormLabel>
          <br></br>
          {!identity ? <Login setActor={setActor} setIdentity={setIdentity}></Login> : <Button onClick={logout}>Logout</Button>}
          <Button onClick={() => loadSecretsData()}>Reload</Button>
        </div>
        <br></br>
        {statusValidated === null ? null : statusValidated ? <Alert severity="success">Rules verified</Alert> : <Alert severity="error">Rules are corrupted</Alert>}
        <br></br>
        <Divider>
          <Chip label="Users"></Chip>
        </Divider>
        <Users actor={actor} data={users} reloadDataFn={loadSecretsData}></Users>
        <br></br>
        <br></br>
        <br></br>
        <Divider>
          <Chip label="Secrets: Active"></Chip>
        </Divider>
        <ActiveSecrets actor={actor} data={status.secrets} reloadDataFn={loadSecretsData}></ActiveSecrets>
        <br></br>
        <Divider>
          <Chip label="Secrets: Pending"></Chip>
        </Divider>
        <PendingSecrets actor={actor} rules={status.pending_secret_reqs} reloadDataFn={loadSecretsData}></PendingSecrets>
        <br></br>
        <br></br>
        <br></br>
        <Divider>
          <Chip label="Consumers: Active"></Chip>
        </Divider>
        <ActiveConsumers actor={actor} data={status.consumers} reloadDataFn={loadSecretsData}></ActiveConsumers>
        <br></br>
        <Divider>
          <Chip label="Consumers: Pending"></Chip>
        </Divider>
        <PendingConsumers actor={actor} rules={status.pending_consumer_reqs} reloadDataFn={loadSecretsData}></PendingConsumers>
      </Container>
    </Root>
  )
}
