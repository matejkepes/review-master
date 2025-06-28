import type { AxiosInstance } from 'axios'
import { api } from 'src/boot/axios'
import { z } from 'zod'
import type { Router } from 'vue-router'
import { useRouter } from 'vue-router'

// Define the schemas using zod
const LoginResponseSchema = z.object({
  token: z.string(),
})

type LoginResponse = z.infer<typeof LoginResponseSchema>

const ClientSchema = z.object({
  id: z.number(),
  name: z.string(),
})

export type Client = z.infer<typeof ClientSchema>

const ClientsResponseSchema = z.object({
  clients: z.array(ClientSchema).nullable(),
})

type ClientsResponse = z.infer<typeof ClientsResponseSchema>

const UserStatsResponseSchema = z.object({
  success: z.boolean(),
  err: z.string(),
  stats: z
    .array(
      z.object({
        client_id: z.number(),
        client_name: z.string(),
        sent: z.number(),
        requested: z.number(),
        group_period: z.string(),
      }),
    )
    .nullable(),
})

export type UserStatsResponse = z.infer<typeof UserStatsResponseSchema>

const ReviewsResponseSchema = z.object({
  locations: z
    .array(
      z.object({
        location_name: z.string(),
        postal_code: z.string(),
        review_ratings: z.object({
          unspecified: z.number(),
          one: z.number(),
          two: z.number(),
          three: z.number(),
          four: z.number(),
          five: z.number(),
        }),
        insights: z.object({
          number_of_business_profile_call_button_clicked: z.number(),
          number_of_business_profile_website_clicked: z.number(),
        }),
      }),
    )
    .nullable()
    .default([]),
})

export type ReviewsResponse = z.infer<typeof ReviewsResponseSchema>

const ReportsResponseSchema = z.object({
  reports: z
    .array(
      z.object({
        report_id: z.number(),
        period_start: z.string(),
        period_end: z.string(),
        generated_at: z.string(),
        client_name: z.string(),
      })
    )
    .nullable()
    .default([]),
})

export type ReportsResponse = z.infer<typeof ReportsResponseSchema>

// Define a generic error class
class ApiError extends Error {
  statusCode: number
  body?: any

  constructor(statusCode: number, message: string, body?: any) {
    super(message)
    this.statusCode = statusCode
    this.body = body
  }
}

class ApiService {
  api: AxiosInstance
  router?: Router

  constructor(api: AxiosInstance) {
    this.api = api

    // Add request interceptor for auth token
    this.api.interceptors.request.use((config) => {
      const user = localStorage.getItem('USER')
      if (user) {
        const { token } = JSON.parse(user)
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Add response interceptor for 401 handling
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('USER')
          this.router?.push('/login')
        }
        return Promise.reject(error)
      },
    )
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    try {
      const response = await this.api.post('/login', { email, password })
      return LoginResponseSchema.parse(response.data)
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }

  async getClients(): Promise<ClientsResponse> {
    try {
      const response = await this.api.get('/auth/clients')
      return ClientsResponseSchema.parse(response.data)
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }

  async getUserStats(
    startDay: string,
    endDay: string,
    timeGrouping: string,
  ): Promise<UserStatsResponse> {
    try {
      const response = await this.api.get('/auth/userstats', {
        params: { start_day: startDay, end_day: endDay, time_grouping: timeGrouping },
      })
      return UserStatsResponseSchema.parse(response.data)
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }

  async getReviews(startTime: string, endTime: string): Promise<ReviewsResponse> {
    try {
      const response = await this.api.get('/auth/reviews', {
        params: { start_time: startTime, end_time: endTime },
      })
      return ReviewsResponseSchema.parse(response.data)
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }

  async getReports(clientId?: number): Promise<ReportsResponse> {
    try {
      const params = clientId ? { client_id: clientId } : {}
      const response = await this.api.get('/auth/reports', { params })
      return ReportsResponseSchema.parse(response.data)
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }

  async getReportHTML(reportId: string): Promise<string> {
    try {
      const response = await this.api.get(`/auth/reports/${reportId}/html`)
      return response.data as string
    } catch (error: any) {
      throw new ApiError(error.response?.status, error.message, error.response?.data)
    }
  }
}

const apiService = new ApiService(api)

// create a composable
export function useApiService() {
  const router = useRouter()
  apiService.router = router
  return { apiService }
}

// export the service
export { apiService, ApiError }
