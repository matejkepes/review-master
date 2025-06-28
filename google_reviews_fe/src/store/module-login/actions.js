import { api } from 'boot/axios'

export function login ({ commit }, user) {
  return new Promise((resolve, reject) => {
    commit('authRequest')
    // console.log('user = %O', user)
    // api.post('/login', user)
    api({ url: '/login', data: user, method: 'POST' })
      .then(resp => {
        const accessToken = resp.data.token
        localStorage.setItem('access_token', accessToken)
        // api.defaults.headers.common['Authorization'] =
        api.defaults.headers.common.Authorization =
          'Bearer ' + accessToken
        commit('authSuccess', accessToken)
        resolve(resp)
      })
      .catch(err => {
        commit('authError')
        localStorage.removeItem('access_token')
        reject(err)
      })
  })
}

export function logout ({ commit }) {
  // return new Promise((resolve, reject) => {
  return new Promise(resolve => {
    commit('logout')
    localStorage.removeItem('access_token')
    // delete api.defaults.headers.common['Authorization']
    delete api.defaults.headers.common.Authorization
    resolve()
  })
}
