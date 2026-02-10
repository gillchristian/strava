export interface StravaActivity {
  id: number;
  name: string;
  type: string;
  sport_type: string;
  distance: number;
  moving_time: number;
  elapsed_time: number;
  average_speed: number;
  average_heartrate?: number;
  max_heartrate?: number;
  start_date_local: string;
  start_date: string;
}

export interface CachedActivities {
  activities: StravaActivity[];
  fetchedAt: number;
}
