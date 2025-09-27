export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const API_ENDPOINTS = {
    CONFIG: '/config',
    LOGS: '/logs',
    HEALTH: '/health',
    PROFILE: '/profile',
    WEBHOOK: '/webhook',
    DEBUG: '/debug'
} as const 