import { useNavigate, useLocation } from 'react-router-dom';

export const useRedirect = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const redirectToLogin = (intendedPath: string) => {
        navigate('/login', { state: { from: intendedPath || location.pathname } });
    };

    const redirectAfterLogin = () => {
        const state = location.state as { from: string } | undefined;
        if (state && state.from) {
            navigate(state.from);
        } else {
            navigate('/');
        }
    };

    return { redirectToLogin, redirectAfterLogin };
};
