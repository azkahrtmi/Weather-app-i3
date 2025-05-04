import { JSX, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

export default function ProtectedRoute({
  children,
}: {
  children: JSX.Element;
}) {
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const validateUser = async () => {
      try {
        const res = await fetch("http://localhost:8080/validate", {
          credentials: "include", //
        });
        if (!res.ok) {
          navigate("/");
        }
      } catch (err) {
        navigate("/");
      } finally {
        setLoading(false);
      }
    };

    validateUser();
  }, [navigate]);

  if (loading) return <div>Loading...</div>;

  return children;
}
