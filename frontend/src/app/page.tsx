'use client';

import { FormEvent, useEffect, useMemo, useState } from 'react';

import { MetricCard } from '@/components/MetricCard';
import { requestJson } from '@/lib/api';

type Task = {
  id: number;
  title: string;
  done: boolean;
  createdAt: string;
};

type Summary = {
  app: string;
  health: string;
  generatedAt: string;
  tasksTotal: number;
  tasksCompleted: number;
  tasksOpen: number;
  status: string;
};

export default function Home() {
  const [summary, setSummary] = useState<Summary | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [health, setHealth] = useState('checking...');
  const [title, setTitle] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        setError(null);
        const [healthResult, summaryResult, tasksResult] = await Promise.all([
          requestJson<{ status: string }>('/healthz'),
          requestJson<Summary>('/api/summary'),
          requestJson<{ tasks: Task[] }>('/api/tasks'),
        ]);
        setHealth(healthResult.status);
        setSummary(summaryResult);
        setTasks(tasksResult.tasks);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unable to reach the backend');
        setHealth('offline');
      }
    };

    void load();
  }, []);

  const completionRate = useMemo(() => {
    if (!summary || summary.tasksTotal === 0) {
      return 0;
    }
    return Math.round((summary.tasksCompleted / summary.tasksTotal) * 100);
  }, [summary]);

  async function refresh() {
    const [summaryResult, tasksResult] = await Promise.all([
      requestJson<Summary>('/api/summary'),
      requestJson<{ tasks: Task[] }>('/api/tasks'),
    ]);
    setSummary(summaryResult);
    setTasks(tasksResult.tasks);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!title.trim()) {
      return;
    }

    setBusy(true);
    setError(null);

    try {
      await requestJson<Task>('/api/tasks', {
        method: 'POST',
        body: JSON.stringify({ title }),
      });
      setTitle('');
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to create task');
    } finally {
      setBusy(false);
    }
  }

  async function toggleTask(task: Task) {
    setBusy(true);
    setError(null);

    try {
      await requestJson<Task>(`/api/tasks/${task.id}`, {
        method: 'PUT',
        body: JSON.stringify({ done: !task.done }),
      });
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to update task');
    } finally {
      setBusy(false);
    }
  }

  async function deleteTask(id: number) {
    setBusy(true);
    setError(null);

    try {
      await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/api/tasks/${id}`, {
        method: 'DELETE',
      });
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to delete task');
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="shell">
      <section className="hero">
        <div>
          <p className="eyebrow">CloudPulse demo</p>
          <h1>Minimal cloud dashboard, wired to a Go API.</h1>
          <p className="lede">
            This starter shows the shape of the roadmap app without the production-only pieces.
            It gives you live health checks, analytics, and a task list you can edit.
          </p>
        </div>
        <div className="status-card">
          <span className={`pill ${health === 'ok' ? 'ok' : 'bad'}`}>{health}</span>
          <strong>{summary?.app ?? 'CloudPulse'}</strong>
          <span>Backend: {process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}</span>
          <span>Mode: demo-ready</span>
        </div>
      </section>

      <section className="stats">
        <MetricCard label="Tasks total" value={summary?.tasksTotal ?? tasks.length} />
        <MetricCard label="Completed" value={summary?.tasksCompleted ?? tasks.filter((task) => task.done).length} />
        <MetricCard label="Open" value={summary?.tasksOpen ?? tasks.filter((task) => !task.done).length} />
        <MetricCard label="Completion" value={`${completionRate}%`} />
      </section>

      <section className="panel">
        <div className="panel-header">
          <h2>Tasks</h2>
          <p>{summary?.generatedAt ? `Refreshed at ${new Date(summary.generatedAt).toLocaleTimeString()}` : 'Loading summary...'}</p>
        </div>

        <form className="task-form" onSubmit={handleSubmit}>
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Add a new task"
            aria-label="New task title"
          />
          <button type="submit" disabled={busy}>
            Add task
          </button>
        </form>

        {error ? <p className="error">{error}</p> : null}

        <div className="task-list">
          {tasks.map((task) => (
            <article className="task" key={task.id}>
              <label>
                <input type="checkbox" checked={task.done} onChange={() => toggleTask(task)} />
                <span className={task.done ? 'done' : ''}>{task.title}</span>
              </label>
              <button type="button" onClick={() => deleteTask(task.id)}>
                Remove
              </button>
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
