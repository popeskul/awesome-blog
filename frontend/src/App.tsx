import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Layout from './components/Layout/Layout';
import HomePage from './pages/HomePage';
import { PostListPage } from './pages/PostListPage';
import { PostDetail } from './components/Posts/PostDetail';
import { CreatePostForm } from './components/Posts/CreatePostForm';
import { EditPostForm } from './components/Posts/EditPostForm';
import { LoginForm } from './components/Auth/LoginForm';
import { RegisterForm } from './components/Auth/RegisterForm';
import { UserProvider } from './contexts/UserContext';
import { AlertProvider } from './contexts/AlertContext';

const App: React.FC = () => {
    return (
        <Router>
            <UserProvider>
                <AlertProvider>
                    <Layout>
                        <Routes>
                            <Route path="/" element={<HomePage />} />
                            <Route path="/posts" element={<PostListPage />} />
                            <Route path="/posts/:id" element={<PostDetail />} />
                            <Route path="/create-post" element={<CreatePostForm />} />
                            <Route path="/edit-post/:id" element={<EditPostForm />} />
                            <Route path="/login" element={<LoginForm />} />
                            <Route path="/register" element={<RegisterForm />} />
                        </Routes>
                    </Layout>
                </AlertProvider>
            </UserProvider>
        </Router>
    );
};

export default App;
