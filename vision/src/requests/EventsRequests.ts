import axios from "axios";
import { API_CONFIG } from "../config/api";

export async function fetchRegionalStages() {
  const res = await axios.get(API_CONFIG.REGIONAL, {
    withCredentials: true,
  });
  return res.data.data;
}

export async function fetchEvent(id: string) {
  const res = await axios.get(`${API_CONFIG.EVENT}/${id}`, {
    withCredentials: true,
  });
  return res.data.data;
}

export async function fetchStagesCount(id: string) {
  const res = await axios.get(`${API_CONFIG.CHILD}/${id}`, {
    withCredentials: true,
  });
  return res.data.data.count;
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