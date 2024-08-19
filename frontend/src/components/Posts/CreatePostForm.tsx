import React, { useState } from 'react';
import { TextField, Button, Box, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { api } from '../../services/api';
import { useUser } from '../../contexts/UserContext';
import { useAlert } from '../../contexts/AlertContext';

export const CreatePostForm: React.FC = () => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const { user } = useUser();
    const { showAlert } = useAlert();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) {
            showAlert('You must be logged in to create a post', 'error');
            return;
        }
        try {
            await api.createPost({ title, content, authorId: user.id });
            showAlert('Post created successfully', 'success');
            navigate('/posts');
        } catch (error) {
            console.error('Failed to create post:', error);
            showAlert('Failed to create post', 'error');
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                Create New Post
            </Typography>
            <TextField
                margin="normal"
                required
                fullWidth
                id="title"
                label="Title"
                name="title"
                autoFocus
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <TextField
                margin="normal"
                required
                fullWidth
                name="content"
                label="Content"
                id="content"
                multiline
                rows={4}
                value={content}
                onChange={(e) => setContent(e.target.value)}
            />
            <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
            >
                Create Post
            </Button>
        </Box>
    );
};
