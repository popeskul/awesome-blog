import React from 'react';
import { Link } from 'react-router-dom';
import { Typography, Button, Box } from '@mui/material';

const HomePage: React.FC = () => {
    return (
        <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            minHeight="80vh"
        >
            <Typography variant="h2" component="h1" gutterBottom>
                Welcome to Our Blog
            </Typography>
            <Typography variant="h5" component="h2" gutterBottom>
                Explore our latest posts and join the conversation
            </Typography>
            <Button
                component={Link}
                to="/posts"
                variant="contained"
                color="primary"
                size="large"
                sx={{ mt: 4 }}
            >
                View All Posts
            </Button>
        </Box>
    );
};

export default HomePage;
