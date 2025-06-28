import VueJwtDecode from 'vue-jwt-decode'

export function isLoggedIn (state) {
  return !!state.accessToken
}

export function authStatus (state) {
  return state.status
}

export function displayCurrentUser (state) {
  try {
    // console.log(VueJwtDecode.decode(state.accessToken))
    return VueJwtDecode.decode(state.accessToken).displayName
  } catch (err) {
    return ''
  }
}

export function roleCurrentUser (state) {
  try {
    // console.log(VueJwtDecode.decode(state.accessToken))
    return VueJwtDecode.decode(state.accessToken).role
  } catch (err) {
    return ''
  }
}
