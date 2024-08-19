import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
    List,
    ListItem,
    ListItemText,
    Button,
    Pagination,
    Box,
    Typography,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { usePosts } from '../../hooks/usePosts';
import { useUser } from '../../contexts/UserContext';
import { useRedirect } from '../../hooks/useRedirect';
import {useAlert} from "../../contexts/AlertContext";

export const PostList: React.FC = () => {
    const { posts, loading, error, page, setPage, totalPages, handleDeletePost } = usePosts();
    const { user } = useUser();
    const [openDialog, setOpenDialog] = useState(false);
    const [postToDelete, setPostToDelete] = useState<string | null>(null);
    const { redirectToLogin } = useRedirect();
    const { showAlert } = useAlert();
    const navigate = useNavigate();

    const handleOpenDialog = (postId: string) => {
        setPostToDelete(postId);
        setOpenDialog(true);
    };

    const handleCloseDialog = () => {
        setOpenDialog(false);
        setPostToDelete(null);
    };

    const confirmDelete = () => {
        if (postToDelete) {
            handleDeletePost(postToDelete);
            handleCloseDialog();
        }
    };

    const canEditOrDelete = (postAuthorId: string) => {
        return user && (user.id === postAuthorId || user.role === 'admin');
    };

    const handleEditClick = (postId: string, postAuthorId: string) => {
        if (user) {
            if (canEditOrDelete(postAuthorId)) {
                navigate(`/edit-post/${postId}`);
            } else {
                showAlert('You do not have permission to edit this post', 'error');
            }
        } else {
            redirectToLogin(`/edit-post/${postId}`);
        }
    };

    const handleCreatePost = () => {
        if (user) {
            navigate('/create-post');
        } else {
            redirectToLogin('/create-post');
        }
    };

    if (loading) return <Typography>Loading...</Typography>;
    if (error) return <Typography color="error">{error}</Typography>;

    return (
        <Box>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Button
                    variant="contained"
                    color="primary"
                    startIcon={<AddIcon />}
                    onClick={handleCreatePost}
                >
                    Create New Post
                </Button>
            </Box>

            {(!posts || posts.length === 0) ? (
                <Typography>No posts available.</Typography>
            ) : (
                <>
                    <List>
                        {posts.map((post) => (
                            <ListItem key={post.id} divider>
                                <ListItemText
                                    primary={<Link to={`/posts/${post.id}`}>{post.title}</Link>}
                                    secondary={`By ${post.authorId} on ${new Date(post.createdAt).toLocaleDateString()}`}
                                />
                                {canEditOrDelete(post.authorId) && (
                                    <>
                                        <Button onClick={() => handleEditClick(post.id, post.authorId)}>
                                            Edit
                                        </Button>
                                        <Button
                                            onClick={() => handleOpenDialog(post.id)}
                                            color="error"
                                        >
                                            Delete
                                        </Button>
                                    </>
                                )}
                            </ListItem>
                        ))}
                    </List>
                    {totalPages > 1 && (
                        <Pagination
                            count={totalPages}
                            page={page}
                            onChange={(event, value) => setPage(value)}
                            color="primary"
                            sx={{ mt: 2, display: 'flex', justifyContent: 'center' }}
                        />
                    )}
                </>
            )}

            <Dialog
                open={openDialog}
                onClose={handleCloseDialog}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    {"Confirm Deletion"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        Are you sure you want to delete this post? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleCloseDialog}>Cancel</Button>
                    <Button onClick={confirmDelete} color="error" autoFocus>
                        Delete
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};
