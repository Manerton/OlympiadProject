import { createContext, useContext, useEffect, useState } from "react";
import type { RegisterForm, UserAuth } from "../types/user";
import axios from "axios";
import { AUTH } from "../../config/api";
import {jwtDecode} from "jwt-decode";

interface JwtPayload {
  sub: string;
  email: string;
  role: number;
  // добавь свои поля
}

interface AuthContextType {
    user: UserAuth | null
    accessToken: string | null
    login: (email: string, password: string) => Promise<void>
    register: (data: RegisterForm) => Promise<void>
    logout: () => void
    refresh: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
    const [user, setUser] = useState<UserAuth | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null)

    const register = async (data: RegisterForm) => {
        const response = await axios.post(AUTH.register, data, {withCredentials: true})
        if (response.status == 200) {
            console.log("success regisger")
        } else {
            console.log("failed register")
        }
    }

    const login = async (email: string, password: string) => {
        const response = await axios.post(AUTH.login, {email, password}, {withCredentials: true})
        const accessToken = response.data.data.access_token
        console.log(response.data.data.access_token)
        setAccessToken(accessToken)
        const decoded = jwtDecode<JwtPayload>(accessToken);
       
        setUser({
            id: decoded.sub,
            Email: decoded.email,
            role:  Number(decoded.role),
        });
    };

    const logout = () => {
        setAccessToken(null)
        setUser(null)
        axios.post(AUTH.logout, {}, {withCredentials: true})
    }

    const refresh = async () => {
        try {
            const response = await axios.post(AUTH.refresh, {}, {withCredentials: true})
            setAccessToken(response.data.data.access_token)
            console.log(response.data.data.access_token)
            const decoded = jwtDecode<JwtPayload>(response.data.data.access_token);
            setUser({
                id: decoded.sub,
                Email: decoded.email,
                role:  Number(decoded.role),
            });
        } catch(err) {
            console.error("refresh failed", err)
            logout()
        }
    }

    useEffect(() => {
        refresh()
    }, [])

    return (
        <AuthContext.Provider value={{user, accessToken, login, register, logout, refresh}}>
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