"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Activity } from "lucide-react";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await fetch("http://localhost:8000/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        throw new Error("Failed to register. User might already exist.");
      }

      // Auto login after register
      const loginRes = await fetch("http://localhost:8000/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await loginRes.json();
      localStorage.setItem("token", data.token);
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-container" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh' }}>
      <div className="glass-panel animate-fade-in" style={{ width: '100%', maxWidth: '400px', padding: '2.5rem' }}>
        <div style={{ display: 'flex', justifyContent: 'center', marginBottom: '2rem' }}>
          <div style={{ background: 'var(--accent-color)', padding: '1rem', borderRadius: '50%' }}>
            <Activity size={32} color="white" />
          </div>
        </div>
        
        <h1 style={{ textAlign: 'center', marginBottom: '0.5rem', fontSize: '1.5rem', fontWeight: 600 }}>Create an Account</h1>
        <p style={{ textAlign: 'center', color: 'var(--text-secondary)', marginBottom: '2rem' }}>Join CloudPulse today</p>
        
        {error && (
          <div style={{ background: 'rgba(239, 68, 68, 0.1)', border: '1px solid var(--danger-color)', color: '#fca5a5', padding: '0.75rem', borderRadius: '8px', marginBottom: '1.5rem', fontSize: '0.875rem' }}>
            {error}
          </div>
        )}

        <form onSubmit={handleRegister} style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>Email</label>
            <input 
              type="email" 
              required
              className="input-field" 
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div>
            <label style={{ display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>Password</label>
            <input 
              type="password" 
              required
              className="input-field" 
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <button type="submit" className="btn" style={{ marginTop: '0.5rem' }} disabled={loading}>
            {loading ? "Creating account..." : "Register"}
          </button>
        </form>

        <p style={{ textAlign: 'center', marginTop: '2rem', fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
          Already have an account? <Link href="/login" className="nav-link" style={{ color: 'var(--accent-hover)' }}>Sign In</Link>
        </p>
      </div>
    </div>
  );
}
