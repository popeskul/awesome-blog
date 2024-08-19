import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { TextField, Button, Box, Typography, CircularProgress } from '@mui/material';
import { api } from '../../services/api';
import { useUser } from '../../contexts/UserContext';
import { useAlert } from '../../contexts/AlertContext';

export const EditPostForm: React.FC = () => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(true);
    const { user } = useUser();
    const { showAlert } = useAlert();
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();

    useEffect(() => {
        const fetchPost = async () => {
            try {
                const post = await api.getPost(id!);
                setTitle(post.title);
                setContent(post.content);
            } catch (error) {
                console.error('Failed to fetch post:', error);
                showAlert('Failed to load post', 'error');
            } finally {
                setLoading(false);
            }
        };

        fetchPost();
    }, [id, showAlert]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) {
            showAlert('You must be logged in to edit a post', 'error');
            return;
        }
        try {
            await api.updatePost(id!, { title, content });
            showAlert('Post updated successfully', 'success');
            navigate(`/posts/${id}`);
        } catch (error) {
            console.error('Failed to update post:', error);
            showAlert('Failed to update post', 'error');
        }
    };

    if (loading) {
        return <CircularProgress />;
    }

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
            <Typography variant="h4" component="h1" gutterBottom>
                Edit Post
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
                Update Post
            </Button>
        </Box>
    );
};
