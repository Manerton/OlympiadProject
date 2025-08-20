import { createContext, useContext, useEffect, useState } from "react";
import type { UserAuth } from "../types/user";
import axios from "axios";
import { AUTH } from "../../config/api";


interface AuthContextType {
    user: UserAuth | null
    accessToken: string | null
    login: (email: string, password: string) => Promise<void>
    logout: () => void
    refresh: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
    const [user, setUser] = useState<UserAuth | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null)

    const login = async (email: string, password: string) => {
        const response = await axios.post(AUTH.login, {email, password}, {withCredentials: true})
        const {accessToken, user} = response.data
        
        setAccessToken(accessToken)
        setUser(user)
    };

    const logout = () => {
        setAccessToken(null)
        setUser(null)
        // axios.post()
    }

    const refresh = async () => {

    }

    useEffect(() => {
        refresh()
    })

    return (
        <AuthContext.Provider value={{user, accessToken, login, logout, refresh}}>
            {children}
        </AuthContext.Provider >
    )
}


export const useAuth = () => {
    const ctx = useContext(AuthContext);
    if (!ctx) 
        throw new Error("useAuth must be used within AuthProvider")
    return ctx
}