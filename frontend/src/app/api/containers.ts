"use client";
import axios from "axios";

const API_URL = process.env.NEXT_PUBLIC_BACKEND_API_URL || "";
const API_KEY = process.env.NEXT_PUBLIC_BACKEND_AUTH_API_KEY || "";
console.log(API_URL, API_KEY)

export interface Container {
  container_id: string;
  ip_address: string;
  name: string;
  status: string;
  ping_time: number;
  last_successful_ping: string;
}

export const fetchContainers = async (): Promise<Container[]> => {
  try {
    const response = await axios.get<Container[]>(`${API_URL}/container_status`, {
      headers: {
        Accept: "application/json",
        "X-Api-Key": API_KEY,
      },
      timeout: 5000,
    });
    return response.data;
  } catch (error) {
    console.error("Ошибка запроса контейнеров:", error);
    throw error;
  }
};
