// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { Navigate } from "react-router-dom"
import { useAuthContext } from "@/context/AuthContext"

export function PrivateRoute({ children }: Readonly<{ children: React.ReactNode }>) {
    const { isAuthenticated, isLoading } = useAuthContext()

    if (isLoading) return <div>Cargando...</div>
    if (!isAuthenticated) return <Navigate to="/login" replace />

    return children
}