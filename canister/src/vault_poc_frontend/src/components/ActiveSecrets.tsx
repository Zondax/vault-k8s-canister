import * as React from 'react'
import { useEffect, useState } from 'react'
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Alert, Button, LinearProgress } from '@mui/material'

import { _SERVICE, Secret } from '../../../declarations/vault_poc_backend/vault_poc_backend.did'
import { ActorSubclass } from '@dfinity/agent'
import { getMaskedPrincipal, getPendingOperationEnumString, getStatusEnumString, parseDateTimeFromUnixTimestamp } from '../lib/utils'

type Props = {
  actor: ActorSubclass<_SERVICE>
  data: Secret[]
  reloadDataFn: () => void
}

export const ActiveSecrets: React.FunctionComponent<Props> = ({ data, reloadDataFn, actor }) => {
  const [errorMsg, setErrorMsg] = useState<string>('')
  const [inProgress, setInProgress] = useState<boolean>(false)

  const revokeRule = async (id: string) => {
    setErrorMsg('')
    setInProgress(true)
    try {
      const resp = await actor.revoke_secret(id)
      'Err' in resp ? setErrorMsg(resp.Err) : reloadDataFn()
    } catch (e) {
      setErrorMsg(String(e))
    }
    setInProgress(false)
  }

  return (
    <div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="PendingSecrets">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Creator</TableCell>
              <TableCell>Creation Timestamp</TableCell>
              <TableCell>Update Timestamp</TableCell>
              <TableCell>TTL</TableCell>
              <TableCell>KubeID</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map(({ id, creator, create_timestamp, update_timestamp, ttl, kube_id, pending_type }) => (
              <TableRow key={id.toString()}>
                <TableCell>{id.toString()}</TableCell>
                <TableCell>{getMaskedPrincipal(creator)}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(create_timestamp).toISOString()}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(update_timestamp).toISOString()}</TableCell>
                <TableCell>{ttl.toString()}</TableCell>
                <TableCell>{kube_id}</TableCell>
                <TableCell>{getPendingOperationEnumString(pending_type)}</TableCell>
                <TableCell>
                  <Button onClick={() => revokeRule(id)}>Revoke</Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {inProgress && <LinearProgress />}
      {errorMsg && <Alert severity="error">{errorMsg}</Alert>}
    </div>
  )
}
