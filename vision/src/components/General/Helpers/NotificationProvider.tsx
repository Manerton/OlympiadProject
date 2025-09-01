import React, { createContext, useContext, useEffect } from "react";
import { toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

const WS_URL = "ws://notification.olymp.local:8095";

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    useEffect(() => {
        const socket = new WebSocket(WS_URL);

        socket.onopen = () => console.log("✅ WebSocket connected");
        socket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                // Показываем toast уведомление
                toast.info(message.data, { position: "bottom-right" });
            } catch (err) {
                console.error("Ошибка разбора WS сообщения", err);
            }
        };
        socket.onclose = () => console.log("❌ WebSocket closed");
        socket.onerror = (err) => console.log("⚠️ WebSocket error", err);

        return () => socket.close();
    }, []);

    return <>{children}</>;
};