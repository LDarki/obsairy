// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { Input } from "@components/ui/Input"
import { Button } from "@components/ui/Button"
import { useState } from "react"
import { Eye, EyeOff } from "lucide-react"
import { motion } from "framer-motion"
import { useAuth } from "@hooks/useAuth"

export function LoginForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [showPassword, setShowPassword] = useState(false)

    const { login } = useAuth()

    const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
        e.preventDefault()
        login.mutate({ email, password })
    }

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-sm mx-auto">
            <h1 className="text-2xl font-bold mb-4 text-center">Log in</h1>
            <Input label="Email" type="email" value={email} onChange={setEmail} required />
            <div className="relative">
                <Input label="Password" type={showPassword ? "text" : "password"} value={password} onChange={setPassword} required />
                <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-2 top-1/2 transform -translate-y-1/2 cursor-pointer w-4 h-4 flex items-center justify-center"
                >
                    <motion.div
                        animate={{ opacity: showPassword ? 0 : 1, scale: showPassword ? 0.5 : 1 }}
                        transition={{ duration: 0.2 }}
                        className="absolute inset-0 flex items-center justify-center"
                    >
                        <Eye />
                    </motion.div>

                    <motion.div
                        animate={{ opacity: showPassword ? 1 : 0, scale: showPassword ? 1 : 0.5 }}
                        transition={{ duration: 0.2 }}
                        className="absolute inset-0 flex items-center justify-center"
                    >
                        <EyeOff />
                    </motion.div>
                </button>
            </div>
            <Button type="submit" size="full">
                Continue
            </Button>
        </form>
    )
}