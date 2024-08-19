import React, { useState } from 'react';
import { TextField, Button, Box } from '@mui/material';
import { useUser } from '../../contexts/UserContext';
import { useRedirect } from '../../hooks/useRedirect';

interface CommentFormProps {
    postId: string;
    onSubmit: (content: string) => void;
}

export const CommentForm: React.FC<CommentFormProps> = ({ postId, onSubmit }) => {
    const [content, setContent] = useState('');
    const { user } = useUser();
    const { redirectToLogin } = useRedirect();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) {
            redirectToLogin(`/posts/${postId}`);
            return;
        }
        onSubmit(content);
        setContent('');
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
            <TextField
                fullWidth
                multiline
                rows={3}
                variant="outlined"
                placeholder="Write a comment..."
                value={content}
                onChange={(e) => setContent(e.target.value)}
                sx={{ mb: 2 }}
                required
            />
            <Button type="submit" variant="contained" color="primary">
                Post Comment
            </Button>
        </Box>
    );
};
