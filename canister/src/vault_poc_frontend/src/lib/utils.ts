export const getStatusEnumString = (status: number) => {
  switch (status) {
    case 1:
      return 'Pending'
    case 2:
      return 'Approved'
    case 3:
      return 'Revoked'
    default:
      return 'Unknown'
  }
}
export const getPendingOperationEnumString = (pendingOp: number) => {
  switch (pendingOp) {
    case 1:
      return 'Create'
    case 2:
      return 'Update'
    case 3:
      return 'None'
    default:
      return 'Unknown'
  }
}

export const getShortPermissionType = (permType: number) => {
  switch (permType) {
    case 1:
      return 'RW'
    case 2:
      return 'R'
    default:
      return 'X'
  }
}

export const getMaskedPrincipal = (principal: string) => {
  return principal.slice(0, 6) + '...' + principal.slice(principal.length - 6)
}

export const parseDateTimeFromUnixTimestamp = (timestamp: string): Date => {
  return new Date(parseInt((BigInt(timestamp) / BigInt(1000000)).toString()))
}
