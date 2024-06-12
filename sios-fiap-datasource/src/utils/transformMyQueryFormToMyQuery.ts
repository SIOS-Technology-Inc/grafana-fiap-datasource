import { MyQueryForm, MyQuery } from '../types';

// transformMyQueryFormToMyQuery関数にrefIdを引数として追加します
export const transformMyQueryFormToMyQuery = (form: MyQueryForm, refId: string): MyQuery => {
  const { point_ids, data_range, start_time, end_time } = form;

  // point_idsをstringの配列に変換します
  const transformedPointIds = point_ids.map(field => field.point_id);

  const query: MyQuery = {
    refId, // 引数から取得したrefIdを設定します
    point_ids: transformedPointIds,
    data_range,
    start_time: {
      ...start_time,
    },
    end_time: {
      ...end_time,
    },
  };

  return query;
};
