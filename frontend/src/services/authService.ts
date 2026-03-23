// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { http } from "@services/http"
import type { User } from "@type/user"

export const authService = {
    login(email: string, password: string) {
        return http<User>("/auth/login", {
            method: "POST",
            body: { email, password },
        })
    },

    register(email: string, password: string) {
        return http<User>("/auth/register", {
            method: "POST",
            body: { email, password },
        })
    },

    me() {
        return http<User>("/auth/sessions/me")
    },

    logout() {
        return http<void>("/auth/sessions/me", {
            method: "DELETE",
        })
    },
}