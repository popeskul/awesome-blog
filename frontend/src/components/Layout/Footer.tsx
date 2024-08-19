import React from 'react';
import { Box, Typography } from '@mui/material';

const Footer: React.FC = () => {
    return (
        <Box component="footer" sx={{ mt: 'auto', py: 3 }}>
            <Typography variant="body2" color="text.secondary" align="center">
                © 2024 Blog App. All rights reserved.
            </Typography>
        </Box>
    );
};

export default Footer;
