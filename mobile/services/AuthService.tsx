import axios from 'axios'
import TokenManager from '../utils/TokenManager'

const API = 'http://localhost:8080/api/auth'

const AuthService = {
  login: async (email: string, password: string): Promise<void> => {
    const { data } = await axios.post(`${API}/login`, { email, password })
    const { accessToken, refreshToken } = data
    await TokenManager.save(accessToken, refreshToken)
  },

  
  refresh: async (): Promise<string> => {
    const refresh = await TokenManager.getRefresh()
    if (!refresh) throw new Error('No refresh token available')

    const { data } = await axios.post(
      `${API}/refresh`,
      {},
      { headers: { Authorization: `Bearer ${refresh}` } }
    )
    const { accessToken } = data
    await TokenManager.save(accessToken, refresh)
    return accessToken
  },

  
  logout: async (): Promise<void> => {
    const refresh = await TokenManager.getRefresh()
    try {
      await axios.post(
        `${API}/logout`,
        {},
        { headers: { Authorization: `Bearer ${refresh}` } }
      )
    } catch (err) {
      console.error('AuthService.logout API error', err)
    } finally {
      await TokenManager.clear()
    }
  },
}

export default AuthService
