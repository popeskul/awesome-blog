import React from 'react';
import { Typography, Box } from '@mui/material';
import { PostList } from '../components/Posts/PostList';

export const PostListPage: React.FC = () => {
    return (
        <Box>
            <Typography variant="h4" component="h1" gutterBottom>
                Blog Posts
            </Typography>
            <PostList />
        </Box>
    );
};
