// import { createContext, useContext, useEffect, useState } from "react";

// // Создаем контекст для роли и ID и name
// const RoleContext = createContext<{ role: string | null; id: number | null; name: string | null } | null>(null);

// interface RoleProviderProps {
//     children: React.ReactNode;
// }

// // Провайдер для роли и ID
// export function RoleProvider({ children }: RoleProviderProps) {
//     const [role, setRole] = useState<string | null>(null);
//     const [id, setID] = useState<number | null>(null);
//     const [name, setName] = useState<string | null>(null);

//     useEffect(() => {
//         const fetchUserInfo = async () => {
//             const response = await fetch("http://localhost:8081/my-info", {
//                 method: "GET",
//                 credentials: "include", // Для отправки cookie
//             });

            
//             if (response.ok) {
//                 const result = await response.json();
      
//                 setRole(result.role || null);
//                 setID(result.id || null);
//                 setName(result.name || null);
//             }
//         };

//         fetchUserInfo();
//     }, []);

//     return (
//         <RoleContext.Provider value={{ role, id, name }}>
//             {children}
//         </RoleContext.Provider>
//     );
// }

// // Хук для использования контекста
// export function useRole() {
//     const context = useContext(RoleContext);

//     if (!context) {
//         throw new Error("useRole must be used within a RoleProvider");
//     }

//     return context; // Возвращает объект { role, id, name }
// }
import React, { createContext, useState, useContext, ReactNode } from "react";

// Определение типа контекста
interface RoleContextType {
  role: string | null;
  id: string | null;
  name: string | null;
  setRoleData: (id: string, role: string, name: string) => void;
  clearRoleData: () => void;
}

// Создание контекста
const RoleContext = createContext<RoleContextType | undefined>(undefined);

// Провайдер контекста
export const RoleProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [role, setRole] = useState<string | null>(null);
  const [id, setId] = useState<string | null>(null);
  const [name, setName] = useState<string | null>(null);

  // Функция для установки данных роли
  const setRoleData = (newId: string, newRole: string, newName: string) => {
    setId(newId);
    setRole(newRole);
    setName(newName);
  };

  // Функция для очистки данных (при logout)
  const clearRoleData = () => {
    setId(null);
    setRole(null);
    setName(null);
  };

  return (
    <RoleContext.Provider value={{ role, id, name, setRoleData, clearRoleData }}>
      {children}
    </RoleContext.Provider>
  );
};

// Хук для удобного использования контекста
export const useRole = () => {
  const context = useContext(RoleContext);
  if (!context) {
    throw new Error("useRole must be used within a RoleProvider");
  }
  return context;
};
