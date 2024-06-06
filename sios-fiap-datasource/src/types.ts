import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export type MyQueryForm = {
  point_ids: Array<{ point_id: string }>;
  data_range: string;
  start_time: {
    time: string;
    link_dashboards: boolean;
  };
  end_time: {
    time: string;
    link_dashboards: boolean;
  };
}

export interface MyQuery extends DataQuery {
  point_ids: string[];
  data_range: string;
  start_time: {
    time: string;
    link_dashboards: boolean;
  };
  end_time: {
    time: string;
    link_dashboards: boolean;
  };
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  url: string;
}