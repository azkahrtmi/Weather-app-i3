import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";
import toast from "react-hot-toast";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error] = useState("");
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();
      if (!response.ok) {
        toast.error(data.error || "Login failed");
        return;
      }
      toast.success("Login successful!");
      navigate("/home");
    } catch (err) {
      toast.error("Network error");
    }
  };

  return (
    <div className="flex flex-col justify-center items-center h-screen bg-weather">
      <div className="flex flex-col items-center bg-white/55 w-[400px] h-[450px] rounded-[25px] shadow-xl p-6">
        <h2 className="font-extrabold text-2xl pt-5 pb-6 text-gray-800">
          Weather Application
        </h2>

        {error && (
          <div className="mb-4 text-red-600 font-semibold">{error}</div>
        )}

        <div className="w-full">
          <label
            className="block text-gray-700 font-semibold mb-2"
            htmlFor="username"
          >
            Username
          </label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white/80 text-gray-800 placeholder-gray-500 italic"
            placeholder="Enter your username"
          />
        </div>

        <div className="w-full pt-5 relative">
          <label
            className="block text-gray-700 font-semibold mb-2"
            htmlFor="password"
          >
            Password
          </label>
          <input
            type={showPassword ? "text" : "password"}
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-2 pr-12 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white/80 text-gray-800 placeholder-gray-500 italic"
            placeholder="Enter your password"
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 bottom-2 flex items-center text-gray-600 hover:text-amber-500"
          >
            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
          </button>
        </div>

        <button
          onClick={handleLogin}
          className="flex justify-center items-center mt-8 text-white bg-gray-700 hover:bg-black w-[200px] h-[50px] text-2xl font-extrabold text-center rounded-[5px] cursor-pointer"
        >
          Log in
        </button>
      </div>
    </div>
  );
}
