// ===============================================================
// Obsairy
//
// https://github.com/LDarki/obsairy
//
// Copyright (c) 2026 LDarki
//
// Licensed under the Apache 2.0 License. See LICENSE file in the project root for full license information.
// ===============================================================

import { useState } from "react";

interface InputProps {
    label: string;
    type?: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    required?: boolean;
}

export const Input = ({ label, type = "text", value, onChange, placeholder = " ", required = false }: InputProps) => {
    const [focused, setFocused] = useState(false);

    const isFloating = focused || !!value;

    return (
        <div className="relative w-full">
            <input
                type={type}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                onFocus={() => setFocused(true)}
                onBlur={() => setFocused(false)}
                placeholder={placeholder}
                required={required}
                className="peer w-full px-4 py-3 border rounded-md border-obsairy-600 focus:outline-none focus:ring-2 focus:ring-obsairy-600 bg-transparent placeholder-transparent"
            />
            <label
                className={`absolute left-4 transition-all duration-200
                    pointer-events-none
                    ${isFloating
                        ? "top-0 text-obsairy-600 text-xs"
                        : "top-3 text-gray-400 text-base"
                    }`}
            >
                {label}
            </label>
        </div>
    );
};