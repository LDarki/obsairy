// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { authService } from "@services/authService"

export function useAuth() {
    const qc = useQueryClient()

    const meQuery = useQuery({
        queryKey: ["auth", "me"],
        queryFn: authService.me,
        retry: false,
        refetchOnWindowFocus: false,
        staleTime: 1000 * 60,
    })

    const loginMutation = useMutation({
        mutationFn: ({ email, password }: { email: string; password: string }) =>
            authService.login(email, password),
        onSuccess: (user) => {
            qc.setQueryData(["auth", "me"], user)
        },
    })

    const logoutMutation = useMutation({
        mutationFn: authService.logout,
        onSuccess: () => {
            qc.removeQueries({ queryKey: ["auth", "me"] })
        },
    })

    return {
        user: meQuery.data,
        isLoading: meQuery.isLoading,
        isAuthenticated: !!meQuery.data,
        login: loginMutation,
        logout: logoutMutation,
    }
}