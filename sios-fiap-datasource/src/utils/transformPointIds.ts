export const transformPointIdsFieldArrayToArray = (point_ids: Array<{ point_id: string }>): string[] => {
  // point_idsをstringの配列に変換します
  const transformedPointIds = point_ids.map(field => field.point_id);
  return transformedPointIds;
};

export const transformPointIdsArrayToFieldArray = (point_ids: string[]): Array<{ point_id: string }> => {
  // point_idsをfieldArrayに変換します
  const transformedPointIds = point_ids.map(point_id => ({ point_id }));
  return transformedPointIds;
};