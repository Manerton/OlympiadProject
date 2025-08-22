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
    initialized: boolean | null
    login: (email: string, password: string) => Promise<void>
    register: (data: RegisterForm) => Promise<void>
    logout: () => void
    refresh: () => Promise<UserAuth | null>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
    const [user, setUser] = useState<UserAuth | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [loading, setLoading] = useState(true); // флаг загрузки авторизации
    const [initialized, setInitialized] = useState(false);

    const register = async (data: RegisterForm) => {
        const response = await axios.post(AUTH.register, data, {withCredentials: true});
        if (response.status === 200) {
            console.log("success register");
        } else {
            console.log("failed register");
        }
    };

    const login = async (email: string, password: string) => {
        const response = await axios.post(AUTH.login, {email, password}, {withCredentials: true});
        const token = response.data.data.access_token;
        setAccessToken(token);
        const decoded = jwtDecode<JwtPayload>(token);
        setUser({
            id: decoded.sub,
            Email: decoded.email,
            role: Number(decoded.role),
        });
    };

    const logout = async () => {
        const token = accessToken; // берем текущий accessToken из состояния
        console.log(token)
        setAccessToken(null);
        setUser(null);

        try {
            await axios.post(
                AUTH.logout,
                {}, // тело запроса
                {
                    withCredentials: true,
                    headers: {
                        Authorization: `Bearer ${token}` // добавляем токен
                    }
                }
            );
        } catch (err) {
            console.error("Failed to logout:", err);
        }
    };
    const refresh = async (): Promise<UserAuth | null> => {
        try {
            const response = await axios.post(AUTH.refresh, {}, { withCredentials: true });
            const token = response.data.data.access_token;
            setAccessToken(token);
            const decoded = jwtDecode<JwtPayload>(token);
            const newUser = { id: decoded.sub, Email: decoded.email, role: Number(decoded.role) };
            setUser(newUser);
            return newUser;
        } catch (err) {
            setUser(null);
            return null;
        } finally {
            setInitialized(true);
        }
    };

    // при монтировании вызываем refresh
    useEffect(() => {
        refresh();
    }, []);

    // // пока идет проверка токена показываем "Загрузка..."
    // if (loading) {
    //     return <div>Загрузка...</div>;
    // }

    return (
        <AuthContext.Provider value={{user, accessToken, login, register, logout, refresh, initialized}}>
            {children}
        </AuthContext.Provider>
    );
};


export const useAuth = () => {
    const ctx = useContext(AuthContext);
    if (!ctx) 
        throw new Error("useAuth must be used within AuthProvider")
    return ctx
}