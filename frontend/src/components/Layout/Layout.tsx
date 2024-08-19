import React from 'react';
import { Container, CssBaseline, ThemeProvider } from '@mui/material';
import Header from './Header';
import Footer from './Footer';
import theme from '../../styles/theme';

interface LayoutProps {
    children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <Header />
            <Container component="main" maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                {children}
            </Container>
            <Footer />
        </ThemeProvider>
    );
};

export default Layout;
