import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import {
    Typography,
    Box,
    Button,
    CircularProgress,
    Divider,
    Pagination,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle
} from '@mui/material';
import {api, Post, Comment} from '../../services/api';
import { useUser } from '../../contexts/UserContext';
import { useAlert } from '../../contexts/AlertContext';
import { CommentList } from '../Comments/CommentList';
import { CommentForm } from '../Comments/CommentForm';

export const PostDetail: React.FC = () => {
    const [post, setPost] = useState<Post | null>(null);
    const [comments, setComments] = useState<Comment[]>([]);
    const [loading, setLoading] = useState(true);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
    const { id } = useParams<{ id: string }>();
    const { user } = useUser();
    const { showAlert } = useAlert();
    const navigate = useNavigate();

    const fetchComments = async (page: number) => {
        try {
            const commentsResponse = await api.getComments(id!, {
                limit: 10,
                page: page,
                sort: 'created_at_asc',
            });
            setComments(commentsResponse.data);
            setTotalPages(Math.ceil(commentsResponse.pagination.total / 10));
        } catch (error) {
            console.error('Failed to fetch comments:', error);
            showAlert('Failed to load comments', 'error');
        }
    };

    useEffect(() => {
        const fetchPostAndComments = async () => {
            try {
                const fetchedPost = await api.getPost(id!);
                setPost(fetchedPost);
                await fetchComments(currentPage);
            } catch (error) {
                console.error('Failed to fetch post or comments:', error);
                showAlert('Failed to load post or comments', 'error');
            } finally {
                setLoading(false);
            }
        };

        fetchPostAndComments();
    }, [id, currentPage, showAlert]);

    const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
        setCurrentPage(value);
    };

    const handleDeleteClick = () => {
        setOpenDeleteDialog(true);
    };

    const handleDeleteConfirm = async () => {
        if (!post) return;

        try {
            await api.deletePost(post.id);
            showAlert('Post deleted successfully', 'success');
            navigate('/posts');
        } catch (error) {
            console.error('Failed to delete post:', error);
            showAlert('Failed to delete post', 'error');
        } finally {
            setOpenDeleteDialog(false);
        }
    };

    const handleDeleteCancel = () => {
        setOpenDeleteDialog(false);
    };


    const handleAddComment = async (content: string) => {
        if (!user || !post) return;

        try {
            await api.createComment({ postId: post.id, content, authorId: user.id });
            await fetchComments(currentPage);
            showAlert('Comment added successfully', 'success');
        } catch (error) {
            console.error('Failed to add comment:', error);
            showAlert('Failed to add comment', 'error');
        }
    };

    const canEditOrDelete = user && post && (user.id === post.authorId || user.role === 'admin');

    if (loading) {
        return <CircularProgress />;
    }

    if (!post) {
        return <Typography>Post not found</Typography>;
    }

    return (
        <Box>
            <Typography variant="h4" component="h1" gutterBottom>
                {post.title}
            </Typography>
            <Typography variant="subtitle1" gutterBottom>
                By {post.authorId} on {new Date(post.createdAt).toLocaleDateString()}
            </Typography>
            <Typography variant="body1" paragraph>
                {post.content}
            </Typography>

            {canEditOrDelete && (
                <Box mt={2}>
                    <Button
                        component={Link}
                        to={`/edit-post/${post.id}`}
                        variant="contained"
                        color="primary"
                        sx={{ mr: 1 }}
                    >
                        Edit
                    </Button>
                    <Button
                        onClick={handleDeleteClick}
                        variant="contained"
                        color="error"
                    >
                        Delete
                    </Button>
                </Box>
            )}

            <Button
                component={Link}
                to="/posts"
                variant="outlined"
                sx={{ mt: 2 }}
            >
                Back to Posts
            </Button>

            <Divider sx={{ my: 4 }} />

            <Typography variant="h5" component="h2" gutterBottom>
                Comments
            </Typography>

            <CommentList comments={comments}/>

            {totalPages > 1 && (
                <Pagination
                    count={totalPages}
                    page={currentPage}
                    onChange={handlePageChange}
                    sx={{ mt: 2, mb: 2 }}
                />
            )}

            <CommentForm postId={post.id} onSubmit={handleAddComment} />

            <Dialog
                open={openDeleteDialog}
                onClose={handleDeleteCancel}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    {"Confirm Post Deletion"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        Are you sure you want to delete this post? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleDeleteCancel}>Cancel</Button>
                    <Button onClick={handleDeleteConfirm} color="error" autoFocus>
                        Delete
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
