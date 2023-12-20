import * as React from 'react'
import { useEffect, useState } from 'react'
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Alert, Button, LinearProgress } from '@mui/material'

import { _SERVICE, Consumer } from '../../../declarations/vault_poc_backend/vault_poc_backend.did'
import { ActorSubclass } from '@dfinity/agent'
import { getMaskedPrincipal, getPendingOperationEnumString, getShortPermissionType, getStatusEnumString, parseDateTimeFromUnixTimestamp } from '../lib/utils'

type Props = {
  actor: ActorSubclass<_SERVICE>
  data: Consumer[]
  reloadDataFn: () => void
}

export const ActiveConsumers: React.FunctionComponent<Props> = ({ data, reloadDataFn, actor }) => {
  const [errorMsg, setErrorMsg] = useState<string>('')
  const [inProgress, setInProgress] = useState<boolean>(false)

  const revokeRule = async (id: string) => {
    setErrorMsg('')
    setInProgress(true)
    try {
      const resp = await actor.revoke_consumer(id)
      'Err' in resp ? setErrorMsg(resp.Err) : reloadDataFn()
    } catch (e) {
      setErrorMsg(String(e))
    }
    setInProgress(false)
  }

  return (
    <div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="PendingConsumers">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Creator</TableCell>
              <TableCell>Creation Timestamp</TableCell>
              <TableCell>Update Timestamp</TableCell>
              <TableCell>KubeID</TableCell>
              <TableCell>Secret ID</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Permission Type</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map(({ id, creator, create_timestamp, update_timestamp, kube_id, pending_type, permission_type, secret_kube_id }) => (
              <TableRow key={id.toString()}>
                <TableCell>{id.toString()}</TableCell>
                <TableCell>{getMaskedPrincipal(creator)}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(create_timestamp).toISOString()}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(update_timestamp).toISOString()}</TableCell>
                <TableCell>{kube_id}</TableCell>
                <TableCell>{secret_kube_id}</TableCell>
                <TableCell>{getPendingOperationEnumString(pending_type)}</TableCell>
                <TableCell>{getShortPermissionType(permission_type)}</TableCell>
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
