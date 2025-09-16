import axios from "axios";
import { API_CONFIG } from "../config/api";
import { MyEvent } from "../components/types/event";

export async function fetchRegionalStages() {
  const res = await axios.get(API_CONFIG.REGIONAL, {
    withCredentials: true,
  });
  return res.data.data;
}

export async function fetchEvent(id: string): Promise<MyEvent> {
  const res = await axios.get(`${API_CONFIG.EVENT}/${id}`, {
    withCredentials: true,
  });
  return res.data.data;
}

export async function fetchStagesCount(id: string) {
  const res = await axios.get(`${API_CONFIG.CHILD}/${id}`, {
    withCredentials: true,
  });
  return res.data.metadata;
}

export interface FetchOlympiadsParams {
  id: string;
  page: number;
  limit: number;
  order: string;
  search?: string;
  selectedDate?: Date | null;
}

export async function fetchOlympiads({
  id,
  page,
  limit,
  order,
  search,
  selectedDate,
}: FetchOlympiadsParams) {
  let url = `${API_CONFIG.CHILD}/${id}?page=${page}&limit=${limit}&order=${order}`;

  if (search) {
    url += `&name=${encodeURIComponent(search)}`;
  }

  if (selectedDate) {
    const dateStr = selectedDate.toISOString().split("T")[0]; // YYYY-MM-DD
    url += `&date=${dateStr}`;
  }

  const headers: Record<string, string> = {};

  const res = await axios.get(url, {
    headers,
    withCredentials: true,
  });

  return res.data; // { data, metadata }
}

export async function fetchOlympiadChildren(id: string) {
  const res = await axios.get(`${API_CONFIG.CHILD}/${id}`, {
    withCredentials: true,
  });
  return res.data.data; // массив событий (с class_number)
}

export async function fetchOlympiadStages(id: string) {
  const res = await axios.get(`${API_CONFIG.STAGES}/${id}`, {
    withCredentials: true,
  });
  return res.data.data; // массив событий (с class_number)
}

export async function axiosUpdateEvent(token: string, event: MyEvent) {
  const res = await axios.put(
    `${API_CONFIG.EVENT}/${event.id}`,
    event, {
    headers: {
      Authorization: `Bearer ${token}` // добавляем токен
    }
  })

}
