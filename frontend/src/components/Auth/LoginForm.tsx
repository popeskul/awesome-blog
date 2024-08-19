import React, { useState } from 'react';
import { TextField, Button, Box, Typography } from '@mui/material';
import { useUser } from '../../contexts/UserContext';
import { useAlert } from '../../contexts/AlertContext';
import { useRedirect } from '../../hooks/useRedirect';
import {AxiosError} from "axios";

export const LoginForm: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const { login } = useUser();
    const { showAlert } = useAlert();
    const { redirectAfterLogin } = useRedirect();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login(username, password);
            showAlert('Logged in successfully', 'success');
            redirectAfterLogin();
        } catch (error) {
            let errorMessage = 'Failed to log in';

            if (error instanceof AxiosError) {
                if (error.response?.data?.error) {
                    errorMessage = `Failed to log in: ${error.response.data.error}`;
                } else if (error.message) {
                    errorMessage = `Failed to log in: ${error.message}`;
                }
            }

            showAlert(errorMessage, 'error');
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
            <Typography variant="h5" component="h1" gutterBottom>
                Login
            </Typography>
            <TextField
                margin="normal"
                required
                fullWidth
                id="username"
                label="Username"
                name="username"
                autoComplete="username"
                autoFocus
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
            >
                Sign In
            </Button>
        </Box>
    );
};
