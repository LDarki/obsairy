// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import React from "react"

type variant = "primary" | "secondary" | "danger" | "success" | "warning" | "default"
type size = "sm" | "md" | "lg" | "xl" | "full"

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    children: React.ReactNode
    disabled?: boolean
    className?: string
    variant?: variant
    size?: size
}

const styles: Record<variant, { bg: string; text: string }> = {
    default: { bg: "bg-obsairy-600", text: "text-white" },
    primary: { bg: "bg-blue-500", text: "text-white" },
    secondary: { bg: "bg-gray-500", text: "text-white" },
    danger: { bg: "bg-red-500", text: "text-white" },
    success: { bg: "bg-green-500", text: "text-white" },
    warning: { bg: "bg-yellow-500", text: "text-white" },
}

const sizes: Record<size, string> = {
    sm: "px-2 py-1",
    md: "px-4 py-2",
    lg: "px-6 py-3",
    xl: "px-8 py-4",
    full: "w-full p-2",
}

export function Button({ children, onClick, disabled, className, variant, size }: Readonly<ButtonProps>) {
    const current = styles[variant || "default"]
    const currentSize = sizes[size || "md"]

    return (
        <button
            onClick={onClick}
            disabled={disabled}
            className={`${currentSize} rounded-md ${current.bg} ${current.text} hover:opacity-70 transition-all transition-duration-200 cursor-pointer ${className}`}
        >
            {children}
        </button>
    )
}