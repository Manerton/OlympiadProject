import axios from "axios";
import { NOTIFY } from "../config/api";

export async function axiosSendCode(token: string, email: string) {
    const res = await axios.get(
        NOTIFY.sendCode,
        {
            withCredentials: true,
            headers: {
                email: email,
                // TODO!! KEKA
                requestToken: "1234567890", 
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
}