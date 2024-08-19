import { useState, useEffect, useCallback } from 'react';
import { api, Post } from '../services/api';
import { useAlert } from '../contexts/AlertContext';

export const usePosts = () => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const { showAlert } = useAlert();

    const fetchPosts = useCallback(async () => {
        setLoading(true);
        try {
            const response = await api.getPosts({
                page,
                limit: 10,
                sort: 'created_at_desc',
            });
            setPosts(response.data);
            setTotalPages(Math.ceil(response.pagination.total / response.pagination.limit));
            setError(null);
        } catch (err) {
            console.error('Error fetching posts:', err);
            setError('Failed to fetch posts');
            showAlert('Failed to fetch posts', 'error');
            setPosts([]);
        } finally {
            setLoading(false);
        }
    }, [page, showAlert]);

    useEffect(() => {
        fetchPosts();
    }, [fetchPosts]);

    const handleDeletePost = async (id: string) => {
        try {
            await api.deletePost(id);
            showAlert('Post deleted successfully', 'success');
            await fetchPosts();
        } catch (err) {
            console.error('Error deleting post:', err);
            showAlert('Failed to delete post', 'error');
        }
    };

    return { posts, loading, error, page, setPage, totalPages, handleDeletePost };
};