import axios from "axios";
import { APPEAL, RESULT } from "../config/api";
import type { Result } from "../components/types/result";
import type { MainEvent } from "../components/types/event";
import type { CreateAppealRequest } from "../components/types/appeal";

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
    const raw = res.data.data.data;

    console.log("Raw events data:", raw);
    return raw
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

export async function axiosCreateAppeal(token: string, data: CreateAppealRequest) {
    const res = await axios.post(
        APPEAL.create,
        data,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data.status;

}

export async function axiosGetAppeal(token: string, appealId: string) {
    const res = await axios.get(
        APPEAL.get + `${appealId}`,
        {
            withCredentials: true,
            headers: {
                Authorization: `Bearer ${token}` // добавляем токен
            }
        }
    );
    return res.data.data;
}