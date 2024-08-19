import React from 'react';
import { List, ListItem, ListItemText, Typography } from '@mui/material';
import { Comment } from '../../services/api';

interface CommentListProps {
    comments: Comment[];
}

export const CommentList: React.FC<CommentListProps> = ({ comments }) => {
    if (comments?.length === 0) {
        return <Typography>No comments available.</Typography>;
    }

    if (!comments) {
        return <Typography>Loading comments...</Typography>;
    }

    return (
        <List>
            {comments.map((comment) => (
                <ListItem key={comment.id} divider>
                    <ListItemText
                        primary={comment.content}
                        secondary={`By ${comment.authorId} on ${new Date(comment.createdAt).toLocaleDateString()}`}
                    />
                </ListItem>
            ))}
        </List>
    );
};