import { Layout } from './components/Layout';
import { LoginButton } from './components/LoginButton';
import { ActivitiesTable } from './components/ActivitiesTable';
import { RefreshButton } from './components/RefreshButton';
import { useAuth } from './hooks/useAuth';
import { useActivities } from './hooks/useActivities';

function App() {
  const { authenticated, loading: authLoading, login, logout } = useAuth();
  const { activities, loading, lastFetched, refresh } = useActivities(authenticated);

  if (authLoading) {
    return (
      <Layout>
        <p className="py-12 text-center text-gray-400">Loading...</p>
      </Layout>
    );
  }

  if (!authenticated) {
    return (
      <Layout>
        <h1 className="mb-2 text-2xl font-bold text-gray-900">Strava Trends</h1>
        <p className="mb-8 text-gray-500">View your running trends from the last 30 days.</p>
        <LoginButton onClick={login} />
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Strava Trends</h1>
          {lastFetched && (
            <p className="text-xs text-gray-400">
              Last updated: {new Date(lastFetched).toLocaleString()}
            </p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <RefreshButton loading={loading} onClick={refresh} />
          <button
            onClick={logout}
            className="rounded-md px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-200 transition-colors"
          >
            Logout
          </button>
        </div>
      </div>
      <ActivitiesTable activities={activities} />
    </Layout>
  );
}

export default App;
