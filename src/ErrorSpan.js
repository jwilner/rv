import React from "react";

export function ErrorSpan({message}) {
    if (!message) {
        return <span/>
    }
    return <span className="label alert">{message}</span>
}