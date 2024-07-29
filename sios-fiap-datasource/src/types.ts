import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  point_ids: Array<{ point_id: string }>;
  data_range: string;
  start_time: {
    time: string;
    link_dashboard: boolean;
  };
  end_time: {
    time: string;
    link_dashboard: boolean;
  };
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  point_ids: [{point_id: ''}],
  data_range: 'period',
  start_time: {time: '', link_dashboard: true},
  end_time: {time: '', link_dashboard: true},
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  url: string;
  server_timezone: string;
}
