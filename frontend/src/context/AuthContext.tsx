// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { createContext, useContext } from "react"
import { useAuth } from "@hooks/useAuth"

const AuthContext = createContext<ReturnType<typeof useAuth> | null>(null)

export function AuthProvider({ children }: Readonly<{ children: React.ReactNode }>) {
    const auth = useAuth()
    return <AuthContext.Provider value={auth} > {children} </AuthContext.Provider>
}

export function useAuthContext() {
    const ctx = useContext(AuthContext)
    if (!ctx) throw new Error("useAuthContext must be used inside AuthProvider")
    return ctx
}