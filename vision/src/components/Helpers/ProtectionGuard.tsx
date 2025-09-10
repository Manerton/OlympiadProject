import { Navigate } from "react-router-dom";
import { useAuth } from "./AuthContext";

interface ProtectedRouteProps {
  children: JSX.Element;
  allowedRoles: number[];
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, allowedRoles }) => {
  const { user } = useAuth(); // допустим user = { id, role }

  if (!user) {
    console.log("User not authenticated, redirecting to /auth");
    return <Navigate to="/auth" replace />;
  }

  if (!allowedRoles.includes(user.role)) {
    return <Navigate to="/forbidden" replace />;
  }

  return children;
};

export default ProtectedRoute;
