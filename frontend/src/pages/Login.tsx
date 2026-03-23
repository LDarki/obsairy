// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { LoginForm } from "@components/LoginForm"
import Logo from "/Logo.svg"
import { useAuth } from "@hooks/useAuth"
import { Navigate } from "react-router-dom"

export function Login() {

    const { isAuthenticated } = useAuth()

    if (isAuthenticated) {
        return <Navigate to="/" />
    }

    return (
        <div className="flex justify-center items-center h-screen">
            <div className="bg-obsairy-800 w-full max-w-10/12 md:max-w-md p-6 rounded-md shadow-md">
                <img src={Logo} alt="Logo" className="w-32 h-32 mx-auto" />
                <LoginForm />
            </div>
        </div>
    )
}