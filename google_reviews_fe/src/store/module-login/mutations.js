export function authRequest (state) {
  state.status = 'loading'
}

export function authSuccess (state, accessToken) {
  state.status = 'success'
  state.accessToken = accessToken
}

export function authError (state) {
  state.status = 'error'
}

export function logout (state) {
  state.accessToken = ''
  state.status = ''
}
