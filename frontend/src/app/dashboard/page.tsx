"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Activity, LogOut, Plus, CheckCircle2, Circle, Trash2, BarChart3, ListTodo, Zap } from "lucide-react";

interface Task {
  id: number;
  title: string;
  done: boolean;
  createdAt: string;
}

interface Analytics {
  tasksTotal: number;
  tasksCompleted: number;
  tasksOpen: number;
  status: string;
}

export default function Dashboard() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [analytics, setAnalytics] = useState<Analytics | null>(null);
  const [newTaskTitle, setNewTaskTitle] = useState("");
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const [tasksRes, analyticsRes] = await Promise.all([
        fetch("http://localhost:8000/api/tasks", {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch("http://localhost:8000/api/analytics", {
          headers: { Authorization: `Bearer ${token}` },
        })
      ]);

      if (tasksRes.status === 401) {
        localStorage.removeItem("token");
        router.push("/login");
        return;
      }

      const tasksData = await tasksRes.json();
      const analyticsData = await analyticsRes.json();

      setTasks(tasksData.tasks || []);
      setAnalytics(analyticsData);
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    router.push("/login");
  };

  const addTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTaskTitle.trim()) return;

    const token = localStorage.getItem("token");
    try {
      const res = await fetch("http://localhost:8000/api/tasks", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title: newTaskTitle }),
      });

      if (res.ok) {
        setNewTaskTitle("");
        fetchData();
      }
    } catch (error) {
      console.error("Failed to add task:", error);
    }
  };

  const toggleTask = async (task: Task) => {
    const token = localStorage.getItem("token");
    try {
      await fetch(`http://localhost:8000/api/tasks/${task.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ done: !task.done }),
      });
      fetchData();
    } catch (error) {
      console.error("Failed to toggle task:", error);
    }
  };

  const deleteTask = async (id: number) => {
    const token = localStorage.getItem("token");
    try {
      await fetch(`http://localhost:8000/api/tasks/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchData();
    } catch (error) {
      console.error("Failed to delete task:", error);
    }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Activity className="animate-fade-in" size={48} color="var(--accent-color)" />
      </div>
    );
  }

  return (
    <div className="page-container animate-fade-in">
      {/* Header */}
      <header className="glass-panel" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '1rem 2rem', marginBottom: '2rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <Activity color="var(--accent-color)" />
          <h1 style={{ fontSize: '1.25rem', fontWeight: 600 }}>CloudPulse</h1>
        </div>
        <button onClick={handleLogout} className="btn" style={{ background: 'transparent', border: '1px solid var(--surface-border)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <LogOut size={16} /> Logout
        </button>
      </header>

      <div className="main-content" style={{ display: 'grid', gridTemplateColumns: '1fr 300px', gap: '2rem' }}>
        {/* Main Column */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
          {/* Analytics Summary */}
          {analytics && (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem' }}>
              <div className="glass-panel" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-secondary)' }}>
                  <ListTodo size={18} /> <span>Total Tasks</span>
                </div>
                <div style={{ fontSize: '2rem', fontWeight: 700 }}>{analytics.tasksTotal}</div>
              </div>
              <div className="glass-panel" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--success-color)' }}>
                  <CheckCircle2 size={18} /> <span>Completed</span>
                </div>
                <div style={{ fontSize: '2rem', fontWeight: 700 }}>{analytics.tasksCompleted}</div>
              </div>
              <div className="glass-panel" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--accent-color)' }}>
                  <Zap size={18} /> <span>Open</span>
                </div>
                <div style={{ fontSize: '2rem', fontWeight: 700 }}>{analytics.tasksOpen}</div>
              </div>
            </div>
          )}

          {/* Task List */}
          <div className="glass-panel" style={{ padding: '2rem' }}>
            <h2 style={{ fontSize: '1.25rem', marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <BarChart3 size={20} color="var(--accent-color)" /> My Tasks
            </h2>
            
            <form onSubmit={addTask} style={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
              <input
                type="text"
                className="input-field"
                placeholder="What needs to be done?"
                value={newTaskTitle}
                onChange={(e) => setNewTaskTitle(e.target.value)}
              />
              <button type="submit" className="btn" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', whiteSpace: 'nowrap' }}>
                <Plus size={16} /> Add Task
              </button>
            </form>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
              {tasks.length === 0 ? (
                <div style={{ textAlign: 'center', color: 'var(--text-secondary)', padding: '2rem 0' }}>
                  No tasks yet. Create one above!
                </div>
              ) : (
                tasks.map(task => (
                  <div key={task.id} style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'space-between', 
                    padding: '1rem', background: 'rgba(15, 23, 42, 0.4)', borderRadius: '8px',
                    border: '1px solid var(--surface-border)', transition: 'all 0.2s ease'
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                      <button 
                        onClick={() => toggleTask(task)} 
                        style={{ background: 'none', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', color: task.done ? 'var(--success-color)' : 'var(--text-secondary)' }}
                      >
                        {task.done ? <CheckCircle2 size={24} /> : <Circle size={24} />}
                      </button>
                      <span style={{ fontSize: '1rem', color: task.done ? 'var(--text-secondary)' : 'var(--text-primary)', textDecoration: task.done ? 'line-through' : 'none' }}>
                        {task.title}
                      </span>
                    </div>
                    <button 
                      onClick={() => deleteTask(task.id)}
                      style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-secondary)', padding: '0.5rem' }}
                      onMouseEnter={(e) => e.currentTarget.style.color = 'var(--danger-color)'}
                      onMouseLeave={(e) => e.currentTarget.style.color = 'var(--text-secondary)'}
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
          <div className="glass-panel" style={{ padding: '2rem' }}>
            <h3 style={{ fontSize: '1rem', marginBottom: '1rem', color: 'var(--text-secondary)' }}>System Status</h3>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '1rem' }}>
              <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: 'var(--success-color)', boxShadow: '0 0 8px var(--success-color)' }}></div>
              <span style={{ fontSize: '0.875rem' }}>API Online</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: 'var(--success-color)', boxShadow: '0 0 8px var(--success-color)' }}></div>
              <span style={{ fontSize: '0.875rem' }}>Database Connected</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
