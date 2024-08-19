import axios, {AxiosInstance} from 'axios';

const API_URL = 'http://localhost:8080';

export interface LoginCredentials {
    username: string;
    password: string;
}

export interface RegisterCredentials {
    username: string;
    email: string;
    password: string;
}

export interface Post {
    id: string;
    title: string;
    content: string;
    authorId: string;
    createdAt: string;
    updatedAt: string;
}

export interface NewPost {
    title: string;
    content: string;
    authorId: string;
}

export interface UpdatePost {
    title?: string;
    content?: string;
}

export interface Comment {
    id: string;
    postId: string;
    content: string;
    authorId: string;
    createdAt: string;
    updatedAt: string;
}

export interface NewComment {
    postId: string;
    content: string;
    authorId: string;
}

export interface User {
    id: string;
    username: string;
    email: string;
    role: 'user' | 'admin';
    createdAt: string;
    updatedAt: string;
}

export interface PaginationParams {
    page?: number;
    limit?: number;
    offset?: number;
    sort?: string
}

export interface Pagination {
    total: number;
    page: number;
    limit: number;
    offset: number;
}

export interface PaginatedResponse<T> {
    data: T[];
    pagination: Pagination;
}

interface ErrorResponse {
    error: string;
}

class Api {
    private instance: AxiosInstance;

    constructor() {
        this.instance = axios.create({
            baseURL: API_URL,
            headers: {
                'Content-Type': 'application/json',
            },
        });
    }

    setToken(token: string) {
        this.instance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }

    async login(credentials: LoginCredentials): Promise<{ token: string }> {
        try {
            const response = await this.instance.post('/auth/login', credentials);
            return response.data;
        } catch (error) {
            if (axios.isAxiosError(error) && error.response) {
                const errorData = error.response.data as ErrorResponse;
                throw new Error(errorData.error || 'An unknown error occurred');
            }
            throw error;
        }
    }

    async getCurrentUser(): Promise<User> {
        try {
            const response = await this.instance.get('/auth/me');
            return response.data;
        } catch (error) {
            if (axios.isAxiosError(error) && error.response) {
                const errorData = error.response.data as ErrorResponse;
                throw new Error(errorData.error || 'Failed to fetch current user data');
            }
            throw error;
        }
    }

    async logout(): Promise<void> {
        try {
            await this.instance.post('/auth/logout');
        } catch (error) {
            if (axios.isAxiosError(error) && error.response) {
                const errorData = error.response.data as ErrorResponse;
                throw new Error(errorData.error || 'Failed to logout');
            }
            throw error;
        }
    }

    async register(credentials: RegisterCredentials): Promise<User> {
        const response = await this.instance.post('/auth/register', credentials);
        return response.data;
    }

    async getPosts(params: PaginationParams = {}): Promise<PaginatedResponse<Post>> {
        const response = await this.instance.get('/api/v1/posts', { params });
        return response.data;
    }

    async createPost(post: NewPost): Promise<Post> {
        const response = await this.instance.post('/api/v1/posts', post);
        return response.data;
    }

    async getPost(postId: string): Promise<Post> {
        const response = await this.instance.get(`/api/v1/posts/${postId}`);
        return response.data;
    }

    async updatePost(postId: string, post: UpdatePost): Promise<Post> {
        const response = await this.instance.put(`/api/v1/posts/${postId}`, post);
        return response.data;
    }

    async deletePost(postId: string): Promise<void> {
        await this.instance.delete(`/api/v1/posts/${postId}`);
    }

    async getComments(postId: string, params: PaginationParams = {}): Promise<PaginatedResponse<Comment>> {
        const response = await this.instance.get(`/api/v1/posts/${postId}/comments`, { params });
        return response.data;
    }

    async createComment(comment: NewComment): Promise<Comment> {
        const response = await this.instance.post(`/api/v1/posts/${comment.postId}/comments`, comment);
        return response.data;
    }
}

export const api = new Api();
