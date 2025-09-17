import axios from "axios";
import { AUTH, PARTICIPANT, SCHOOLS, USER } from "../config/api";
import type { ChangePasswordForm, ForgotPasswordForm, RegisterForm, User, UserParticipant } from "../components/types/user";
import type { School } from "../components/types/schools";

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
        USER.info + `${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data as User;
}

export async function axiosSSOUserParticipantInfo(token: string, userId: string) {
    const res = await axios.get(
        PARTICIPANT.info + `${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    const data = res.data.data;

    return {
        User: {
            email: data.email,
            firstname: data.firstname,
            surname: data.surname,
            patronymic: data.patronymic,
            phone_number: data.phone_number,
            birthdate: data.birthdate,
            gender: data.gender,
            role: data.role,
        },
        school: data.school_id,
        disability: Number(data.disability),
        classnumber: Number(data.class_number),
        citezenship: Number(data.citizenship),
    };
}

export async function axiosSSOAllSchools(): Promise<School[]> {
    const res = await axios.get(SCHOOLS.all);
    return res.data.data as School[]; // data → []SchoolResponseDTO
}

export async function axiosSSOChangePassword(token: string, data: ChangePasswordForm) { 
    const res = await axios.post(
        USER.changePassword + data.user_id,
        data, {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.status
}

export async function axiosSSOForgotPassword(data: ForgotPasswordForm) { 
    const res = await axios.post(
        AUTH.forgotPassword,
        data
    );
    return res.status
}