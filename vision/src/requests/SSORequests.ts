import axios from "axios";
import { AUTH, USER } from "../config/api";
import type { RegisterForm } from "../components/types/user";

export async function axiosSSORegister(data: RegisterForm) {
    const res = await axios.post(AUTH.register, data, { withCredentials: true });
    return res.status;
}

export async function axiosSSOLogin(email: string, password: string): Promise<string> {
    const res = await axios.post(AUTH.login, { email, password }, { withCredentials: true });
    const token = res.data.data.access_token;
    return token
}

export async function axiosSSOLogout(token: string) {
    const res = await axios.post(
        AUTH.logout,
        {}, // тело запроса
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.status
}

export async function axiosSSORefresh() {
    const res = await axios.post(AUTH.refresh, {}, { withCredentials: true });
    const token = res.data.data.access_token;
    return token
}

export async function axiosSSOUserInfo(token: string, userId: string) {
    const res = await axios.get(
        USER.info + `/${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data
}

export async function axiosSSOUserParticipantInfo(token: string, userId: string) {
    const res = await axios.get(
        USER.info + `/${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data
}