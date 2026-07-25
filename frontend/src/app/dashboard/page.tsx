"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Activity, LogOut, Plus, CheckCircle2, Circle, Trash2, BarChart3, ListTodo, Zap } from "lucide-react";
import { requestJson } from "@/lib/api";

interface Task {
  id: number;
  title: string;
  done: boolean;
  attachmentUrl?: string;
  attachmentName?: string;
  attachmentSize?: number;
  mimeType?: string;
  thumbnailUrl?: string;
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
  const [attachment, setAttachment] = useState<File | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const [tasksData, analyticsData] = await Promise.all([
        requestJson<{ tasks: Task[] }>("/tasks"),
        requestJson<Analytics>("/analytics")
      ]);

      setTasks(tasksData.tasks || []);
      setAnalytics(analyticsData);
    } catch (error: any) {
      console.error("Failed to fetch data:", error);
      if (error.message.includes("401") || error.message.includes("unauthorized")) {
        localStorage.removeItem("token");
        router.push("/login");
      }
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

    try {
      const task = await requestJson<Task>("/tasks", {
        method: "POST",
        body: JSON.stringify({ title: newTaskTitle }),
      });

      if (attachment) {
        const formData = new FormData();
        formData.append("attachment", attachment);
        await requestJson(`/tasks/${task.id}/attachment`, {
          method: "POST",
          body: formData,
        });
      }

      setNewTaskTitle("");
      setAttachment(null);
      fetchData();
    } catch (error) {
      console.error("Failed to add task:", error);
    }
  };

  const toggleTask = async (task: Task) => {
    try {
      await requestJson(`/tasks/${task.id}`, {
        method: "PUT",
        body: JSON.stringify({ done: !task.done }),
      });
      fetchData();
    } catch (error) {
      console.error("Failed to toggle task:", error);
    }
  };

  const deleteTask = async (id: number) => {
    try {
      await requestJson(`/tasks/${id}`, {
        method: "DELETE",
      });
      fetchData();
    } catch (error) {
      console.error("Failed to delete task:", error);
    }
  };

  const deleteAttachment = async (id: number) => {
    try {
      await requestJson(`/tasks/${id}/attachment`, {
        method: "DELETE",
      });
      fetchData();
    } catch (error) {
      console.error("Failed to delete attachment:", error);
    }
  };

  const viewAttachment = async (id: number) => {
    try {
      const data = await requestJson<{url: string}>(`/tasks/${id}/attachment/url`);
      if (data.url) {
        window.open(data.url, "_blank");
      }
    } catch (error) {
      console.error("Failed to fetch attachment url:", error);
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
            
            <form onSubmit={addTask} style={{ display: 'flex', flexDirection: 'column', gap: '1rem', marginBottom: '2rem' }}>
              <div style={{ display: 'flex', gap: '1rem' }}>
                <input
                  type="text"
                  className="input-field"
                  placeholder="What needs to be done?"
                  value={newTaskTitle}
                  onChange={(e) => setNewTaskTitle(e.target.value)}
                  style={{ flex: 1 }}
                />
                <button type="submit" className="btn" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', whiteSpace: 'nowrap' }}>
                  <Plus size={16} /> Add Task
                </button>
              </div>
              <div>
                <input
                  type="file"
                  onChange={(e) => setAttachment(e.target.files ? e.target.files[0] : null)}
                  style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}
                />
                {attachment && <span style={{ marginLeft: '0.5rem', fontSize: '0.75rem', color: 'var(--accent-color)' }}>Ready to upload</span>}
              </div>
            </form>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
              {tasks.length === 0 ? (
                <div style={{ textAlign: 'center', color: 'var(--text-secondary)', padding: '2rem 0' }}>
                  No tasks yet. Create one above!
                </div>
              ) : (
                tasks.map(task => (
                  <div key={task.id} style={{
                    display: 'flex', flexDirection: 'column',
                    padding: '1rem', background: 'rgba(15, 23, 42, 0.4)', borderRadius: '8px',
                    border: '1px solid var(--surface-border)', transition: 'all 0.2s ease'
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
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
                    {(task.attachmentUrl || task.attachmentName) && (
                      <div style={{ marginTop: '1rem', marginLeft: '3rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
                        {task.thumbnailUrl ? (
                          <button onClick={() => viewAttachment(task.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0 }} title="View Full Attachment">
                            <img src={task.thumbnailUrl} alt="Thumbnail" style={{ width: '80px', height: '80px', objectFit: 'cover', borderRadius: '4px' }} />
                          </button>
                        ) : task.mimeType?.startsWith('image/') ? (
                          <button onClick={() => viewAttachment(task.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0 }} title="View Full Attachment">
                            <div style={{ width: '80px', height: '80px', background: 'rgba(255,255,255,0.05)', display: 'flex', alignItems: 'center', justifyContent: 'center', borderRadius: '4px', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
                              Click to view
                            </div>
                          </button>
                        ) : (
                          <button onClick={() => viewAttachment(task.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--accent-color)', fontSize: '0.875rem' }}>
                            View Attachment
                          </button>
                        )}
                        <button 
                          onClick={() => deleteAttachment(task.id)}
                          style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--danger-color)', fontSize: '0.75rem' }}
                        >
                          Remove File
                        </button>
                      </div>
                    )}
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
