import React, { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import { api, User } from '../services/api';

interface UserContextProps {
    user: User | null;
    isAuthenticated: boolean;
    login: (username: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    register: (username: string, email: string, password: string) => Promise<void>;
}

const UserContext = createContext<UserContextProps | undefined>(undefined);

export const UserProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        const initializeAuth = async () => {
            const token = localStorage.getItem('token');
            if (token) {
                try {
                    api.setToken(token);
                    const userData = await api.getCurrentUser();
                    setUser(userData);
                } catch (error) {
                    console.error('Failed to initialize authentication:', error);
                    localStorage.removeItem('token');
                    api.setToken('');
                }
            }
        };

        initializeAuth();
    }, []);

    const login = async (username: string, password: string) => {
        const { token } = await api.login({ username, password });
        localStorage.setItem('token', token);
        api.setToken(token);
        const userData = await api.getCurrentUser();
        setUser(userData);
    };

    const logout = async () => {
        try {
            await api.logout();
        } finally {
            localStorage.removeItem('token');
            api.setToken('');
            setUser(null);
        }
    };

    const register = async (username: string, email: string, password: string) => {
        const userData = await api.register({ username, email, password });
        const { token } = await api.login({ username, password });
        localStorage.setItem('token', token);
        api.setToken(token);
        setUser(userData);
    };

    return (
        <UserContext.Provider value={{
            user,
            isAuthenticated: user !== null,
            login,
            logout,
            register
        }}>
            {children}
        </UserContext.Provider>
    );
};

export const useUser = () => {
    const context = useContext(UserContext);
    if (context === undefined) {
        throw new Error('useUser must be used within a UserProvider');
    }
    return context;
};
