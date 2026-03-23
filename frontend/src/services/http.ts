// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================


const API_URL = import.meta.env.VITE_API_URL

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE"

type RequestOptions<TBody = unknown> = {
    method?: HttpMethod
    body?: TBody
    headers?: HeadersInit
}

/**
 * Cliente HTTP base
 */
export async function http<TResponse>(
    endpoint: string,
    { method = "GET", body, headers }: RequestOptions = {}
): Promise<TResponse> {
    const res = await fetch(`${API_URL}${endpoint}`, {
        method,
        credentials: "include", // cookies httpOnly
        headers: {
            "Content-Type": "application/json",
            ...headers,
        },
        body: body ? JSON.stringify(body) : undefined,
    })

    if (!res.ok) {
        const error = await parseError(res)
        throw error
    }

    if (res.status === 204) {
        return null as TResponse
    }

    return res.json()
}

/**
 * Intenta parsear errores del backend
 */
async function parseError(res: Response): Promise<Error> {
    try {
        const data = await res.json()
        return new Error(data.message || "HTTP_ERROR")
    } catch {
        return new Error(res.statusText || "HTTP_ERROR")
    }
}