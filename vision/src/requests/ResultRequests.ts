import axios from "axios";
import { RESULT } from "../config/api";
import type { Result } from "../components/types/result";
import type { MainEvent } from "../components/types/event";

export async function axiosGetHistoryEventsByUser(token: string, userId: string) {
    const res = await axios.get(
        RESULT.allEventsByUserId + `${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data as MainEvent[];
}

export async function axiosResultGetByEventUser(token: string, eventId: string, userId: string) {
    const res = await axios.get(
        RESULT.allByEventIdUserId + `${eventId}/` + `${userId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data as Result[];
}