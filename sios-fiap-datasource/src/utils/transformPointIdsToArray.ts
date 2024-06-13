export const transformPointIdsToArray = (point_ids: Array<{ point_id: string }>): string[] => {
  // point_idsをstringの配列に変換します
  const transformedPointIds = point_ids.map(field => field.point_id);
  return transformedPointIds;
};
