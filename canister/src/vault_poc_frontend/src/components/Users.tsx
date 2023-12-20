import * as React from 'react'
import { useEffect, useState } from 'react'
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Alert, Button, LinearProgress, Input, Checkbox } from '@mui/material'

import { _SERVICE, Secret, User } from '../../../declarations/vault_poc_backend/vault_poc_backend.did'
import { ActorSubclass } from '@dfinity/agent'
import { getMaskedPrincipal, getPendingOperationEnumString, getStatusEnumString, parseDateTimeFromUnixTimestamp } from '../lib/utils'
import { Principal } from '@dfinity/principal'

type Props = {
  actor: ActorSubclass<_SERVICE>
  data: User[]
  reloadDataFn: () => void
}

export const Users: React.FunctionComponent<Props> = ({ data, reloadDataFn, actor }) => {
  const [errorMsg, setErrorMsg] = useState<string>('')
  const [inProgress, setInProgress] = useState<boolean>(false)
  const [newUserId, setNewUserId] = useState<string>('')
  const [newUserRoot, setNewUserRoot] = useState<boolean>(false)

  const removeUser = async (id: Principal, root: boolean) => {
    setErrorMsg('')
    setInProgress(true)
    try {
      const resp = await (root ? actor.remove_privileged_user(id) : actor.remove_user(id))
      'Err' in resp ? setErrorMsg(resp.Err) : reloadDataFn()
    } catch (e) {
      setErrorMsg(String(e))
    }
    setInProgress(false)
  }

  const createUser = async (id: string, root: boolean) => {
    setErrorMsg('')
    setInProgress(true)
    try {
      const resp = await (root ? actor.add_privileged_user(Principal.fromText(id)) : actor.add_user(Principal.fromText(id)))
      'Err' in resp ? setErrorMsg(resp.Err) : reloadDataFn()
    } catch (e) {
      setErrorMsg(String(e))
    }
    setInProgress(false)
    setNewUserId('')
    setNewUserRoot(false)
  }

  return (
    <div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="Users">
          <TableHead>
            <TableRow>
              <TableCell>Principal</TableCell>
              <TableCell>Creator</TableCell>
              <TableCell>Creation Timestamp</TableCell>
              <TableCell>Update Timestamp</TableCell>
              <TableCell>Root</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map(({ id, root, create_timestamp, update_timestamp, creator }) => (
              <TableRow key={id.toString()}>
                <TableCell>{getMaskedPrincipal(id.toString())}</TableCell>
                <TableCell>{getMaskedPrincipal(creator)}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(create_timestamp).toISOString()}</TableCell>
                <TableCell>{parseDateTimeFromUnixTimestamp(update_timestamp).toISOString()}</TableCell>
                <TableCell>{String(root)}</TableCell>
                <TableCell>
                  <Button onClick={() => removeUser(id, root)}>Remove</Button>
                </TableCell>
              </TableRow>
            ))}
            <TableRow key="new">
              <TableCell>
                <Input value={newUserId} onChange={e => setNewUserId(e.target.value)} />
              </TableCell>
              <TableCell>-</TableCell>
              <TableCell>-</TableCell>
              <TableCell>-</TableCell>
              <TableCell>
                <Checkbox checked={newUserRoot} onChange={e => setNewUserRoot(e.target.checked)}></Checkbox>
              </TableCell>
              <TableCell>
                <Button onClick={() => createUser(newUserId, newUserRoot)}>Add</Button>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </TableContainer>
      {inProgress && <LinearProgress />}
      {errorMsg && <Alert severity="error">{errorMsg}</Alert>}
    </div>
  )
}
